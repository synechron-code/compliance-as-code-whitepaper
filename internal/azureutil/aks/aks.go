package aks

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"citihub.com/compliance-as-code/internal/azureutil"
	"citihub.com/compliance-as-code/internal/azureutil/group"
	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-08-01/containerservice"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-02-01/resources"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

// ListAllAKS return all AKS clusters within the Subscription defined by the AZURE_SUBSCRIPTION_ID environment variable.
func ListAllAKS(ctx context.Context) (containerservice.ManagedClusterListResultIterator, error) {
	c := client()
	log.Printf("[DEBUG] subscriptionID: %v", c.SubscriptionID)
	r, err := c.ListComplete(ctx)
	if err != nil {
		log.Printf("Unable to list Managed Clusters: %v", err)
	}
	return r, err
}

// RBACEnabled checks whether or not RBAC is enabled for the Managed Cluster specified by environment variables AKS_NAME and AKS_RG.
func RBACEnabled(ctx context.Context) (*bool, error) {

	rg, bRg := os.LookupEnv("AKS_RG")
	name, bName := os.LookupEnv("AKS_NAME")

	if !bRg || !bName {
		log.Printf("Either of AKS_RG or AKS_NAME are not specified, but are required for this test.")
	}

	r, err := client().Get(ctx, rg, name)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve properties of AKS cluster - %v", err)
	}

	return r.ManagedClusterProperties.EnableRBAC, nil
}

// CreateCluster creates a cluster named 'bddaks' in the resource group 'bdd-aks-rbac-prevent-rg'
func CreateCluster(ctx context.Context) (e error) {

	e = nil
	clusterName := "bddaks"
	targetRg := "bdd-aks-rbac-prevent-rg"

	rg, e := group.Create(ctx, targetRg)
	if e != nil {
		log.Printf("failed to create resource group '%s', '%v'", targetRg, e)
		return
	}

	var agentPoolProfiles = []containerservice.ManagedClusterAgentPoolProfile{{
		Name:   &clusterName,
		Count:  to.Int32Ptr(1),
		VMSize: "Standard_B2s",
	}}

	ch := make(chan string)

	var f containerservice.ManagedClustersCreateOrUpdateFuture

	log.Println(fmt.Sprintf("creating cluster '%s' in resource group '%v'", clusterName, *rg.Name))

	c := client()

	f, e = c.CreateOrUpdate(ctx, *rg.Name, clusterName, containerservice.ManagedCluster{
		Location: to.StringPtr(os.Getenv("AZURE_LOCATION")),
		ManagedClusterProperties: &containerservice.ManagedClusterProperties{
			KubernetesVersion: to.StringPtr("1.15.5"),
			DNSPrefix:         &clusterName,
			AgentPoolProfiles: &agentPoolProfiles,
			ServicePrincipalProfile: &containerservice.ManagedClusterServicePrincipalProfile{
				ClientID: to.StringPtr(os.Getenv("AZURE_CLIENT_ID")),
				Secret:   to.StringPtr(os.Getenv("AZURE_CLIENT_SECRET")),
			},
			EnableRBAC: to.BoolPtr(true),
		},
	})

	go func() {
		if e != nil {
			log.Printf("Failed to create cluster, %v", e)
			ch <- "Failed to create cluster [1]"
			ch <- "Nothing to clean up [2]"
			return
		}
	}()

	go cleanup(ctx, &c, f, rg, ch)

	for i := 0; i < 2; i++ {
		log.Print(<-ch)
	}

	return
}

func cleanup(ctx context.Context, c *containerservice.ManagedClustersClient, f containerservice.ManagedClustersCreateOrUpdateFuture, rg resources.Group, ch chan string) error {

	cluster := containerservice.ManagedCluster{}
	var err error

	// either provisioning succeeds...
	for i := 0; i <= 15; i++ {
		cluster, err = f.Result(*c)

		if err == nil && !strings.EqualFold(*cluster.ProvisioningState, "Succeeded") {
			log.Printf("Cluster not provisioned after %vs, checking again in 60s", i*60)
			time.Sleep(60 * time.Second)
			continue
		}

		ch <- "Cluster provisioned [1]" //need to clean up and send another message
		break
	}

	// ...or it doesn't
	if err != nil || cluster.Name == nil {
		log.Printf("Failed to provision Cluster in '%s' after 15 mins, %v", *rg.Name, err)
		ch <- "Provisioning timeout, not waiting any longer [1]"
		ch <- "May need manual cleanup [2]"
		return err
	}

	// if we've got this far, we can delete the cluster
	if cluster.Name != nil && strings.EqualFold(*cluster.ProvisioningState, "Succeeded") {
		_, err = c.Delete(ctx, *rg.Name, *cluster.Name)
		if err != nil {
			log.Printf("Failed to request Deletion of '%s' in '%s', %v", *cluster.Name, *rg.Name, err)
			ch <- "Deletion request failed [2]"
			return err
		}
		ch <- "Deletion requested [2]"
	}

	return nil
}

func client() containerservice.ManagedClustersClient {
	c := containerservice.NewManagedClustersClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorize Managed Clusters client: %v", err)
	}
	return c
}

package network

import (
	"context"
	"fmt"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-08-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

func getIPClient() (ipClient network.PublicIPAddressesClient, err error) {
	ipClient = network.NewPublicIPAddressesClient(azureutil.GetAzureSubscriptionID())
	// create an authorizer from env vars or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		ipClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return
}

// CreatePublicIP creates a new public IP
func CreatePublicIP(ctx context.Context, ipName string, tags map[string]*string) (ip network.PublicIPAddress, err error) {
	ipClient, err := getIPClient()
	if err != nil {
		log.Fatalf("Unabled to get IPClient: %v", err)
		return
	}
	future, err := ipClient.CreateOrUpdate(
		ctx,
		azureutil.GetAzureResourceGP(),
		ipName,
		network.PublicIPAddress{
			Name:     to.StringPtr(ipName),
			Location: to.StringPtr(azureutil.GetAzureLocation()),
			PublicIPAddressPropertiesFormat: &network.PublicIPAddressPropertiesFormat{
				PublicIPAddressVersion:   network.IPv4,
				PublicIPAllocationMethod: network.Static,
			},
			Tags: tags,
		},
	)

	if err != nil {
		return ip, fmt.Errorf("cannot create public ip address: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, ipClient.Client)
	if err != nil {
		return ip, fmt.Errorf("cannot get public ip address create or update future response: %v", err)
	}

	return future.Result(ipClient)
}

// DeletePublicIP deletes an existing public IP
func DeletePublicIP(ctx context.Context, ipName string) error {
	ipClient, err := getIPClient()
	if err != nil {
		log.Fatalf("Unabled to get IPClient: %v", err)
		return err
	}
	future, err := ipClient.Delete(ctx, azureutil.GetAzureResourceGP(), ipName)

	if err != nil {
		return fmt.Errorf("cannot delete public ip [ %v ] address: %v", ipName, err)
	}

	err = future.WaitForCompletionRef(ctx, ipClient.Client)
	if err != nil {
		return fmt.Errorf("cannot get delete ip address future response: %v", err)
	}
	log.Printf("%v publicIP should be deleted", ipName)
	return nil
}

// GetPublicIP returns an existing public IP
func GetPublicIP(ctx context.Context, ipName string) (publicIP network.PublicIPAddress, err error) {
	ipClient, err := getIPClient()
	if err != nil {
		log.Fatalf("Unabled to get IPClient: %v", err)
	}
	publicIP, err = ipClient.Get(ctx, azureutil.GetAzureResourceGP(), ipName, "")
	return
}

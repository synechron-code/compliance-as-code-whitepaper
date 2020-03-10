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

// CreatePublicIP creates a new public IP
func CreatePublicIP(ctx context.Context, ipName string, tags map[string]*string) (ip network.PublicIPAddress, err error) {
	c, err := ipClient()
	if err != nil {
		log.Fatalf("Unabled to get IPClient: %v", err)
		return
	}
	future, err := c.CreateOrUpdate(
		ctx,
		azureutil.ResourceGroup(),
		ipName,
		network.PublicIPAddress{
			Name:     to.StringPtr(ipName),
			Location: to.StringPtr(azureutil.Location()),
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

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return ip, fmt.Errorf("cannot get public ip address create or update future response: %v", err)
	}

	return future.Result(c)
}

// DeletePublicIP deletes an existing public IP
func DeletePublicIP(ctx context.Context, ipName string) error {
	c, err := ipClient()
	if err != nil {
		log.Fatalf("Unabled to get IPClient: %v", err)
		return err
	}
	future, err := c.Delete(ctx, azureutil.ResourceGroup(), ipName)

	if err != nil {
		return fmt.Errorf("cannot delete public ip [ %v ] address: %v", ipName, err)
	}

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return fmt.Errorf("cannot get delete ip address future response: %v", err)
	}
	log.Printf("[DEBUG] %v publicIP should be deleted", ipName)
	return nil
}

// PublicIP returns an existing Public IP by name from the Resource Group created during testing (a test Resource Group name in the form 'test[a-z]{6}resourceGP').
func PublicIP(ctx context.Context, name string) (network.PublicIPAddress, error) {
	c, err := ipClient()
	if err != nil {
		log.Fatalf("Unabled to get IPClient: %v", err)
	}
	return c.Get(ctx, azureutil.ResourceGroup(), name, "")
}

func ipClient() (c network.PublicIPAddressesClient, err error) {
	c = network.NewPublicIPAddressesClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorise Public IP Addresses client: %v", err)
	}
	return
}

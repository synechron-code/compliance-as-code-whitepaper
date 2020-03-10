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

// CreateVirtualNetwork creates a Virtual Network with CIDR 10.0.0.0/8  in the Subscription configured by environment variable AZURE_SUBSCRIPTION_ID.
func CreateVirtualNetwork(ctx context.Context, name string) (vnet network.VirtualNetwork, err error) {
	c := vnetClient()
	future, err := c.CreateOrUpdate(
		ctx,
		azureutil.ResourceGroup(),
		name,
		network.VirtualNetwork{
			Location: to.StringPtr(azureutil.Location()),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{"10.0.0.0/8"},
				},
			},
		})

	if err != nil {
		return vnet, fmt.Errorf("cannot create virtual network: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return vnet, fmt.Errorf("cannot get the vnet create or update future response: %v", err)
	}

	return future.Result(c)
}

// CreateVirtualNetworkAndSubnets creates a Virtual Network with CIDR 10.0.0.0/8 and Subnets 10.0.0.0/16 and 10.1.0.0/16 in the Subscription configured by environment variable AZURE_SUBSCRIPTION_ID.
func CreateVirtualNetworkAndSubnets(ctx context.Context, name, subnet1Name, subnet2Name string, tags map[string]*string) (vnet network.VirtualNetwork, err error) {
	c := vnetClient()
	future, err := c.CreateOrUpdate(
		ctx,
		azureutil.ResourceGroup(),
		name,
		network.VirtualNetwork{
			Location: to.StringPtr(azureutil.Location()),
			VirtualNetworkPropertiesFormat: &network.VirtualNetworkPropertiesFormat{
				AddressSpace: &network.AddressSpace{
					AddressPrefixes: &[]string{"10.0.0.0/8"},
				},
				Subnets: &[]network.Subnet{
					{
						Name: to.StringPtr(subnet1Name),
						SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
							AddressPrefix: to.StringPtr("10.0.0.0/16"),
						},
					},
					{
						Name: to.StringPtr(subnet2Name),
						SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
							AddressPrefix: to.StringPtr("10.1.0.0/16"),
						},
					},
				},
			},
			Tags: tags,
		})

	if err != nil {
		return vnet, fmt.Errorf("cannot create virtual network: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return vnet, fmt.Errorf("cannot get the vnet create or update future response: %v", err)
	}

	return future.Result(c)
}

// DeleteVirtualNetwork deletes a Virtual Network by name in the Subscription configured by environment variable AZURE_SUBSCRIPTION_ID.
func DeleteVirtualNetwork(ctx context.Context, name string) (network.VirtualNetworksDeleteFuture, error) {
	vnetClient := vnetClient()
	return vnetClient.Delete(ctx, azureutil.ResourceGroup(), name)
}

// ListAllVNetByResourceGroup returns the VNets in the given Resource Group in the Subscription configured by environment variable AZURE_SUBSCRIPTION_ID.
func ListAllVNetByResourceGroup(ctx context.Context, rgName string) (result network.VirtualNetworkListResultIterator, err error) {
	return vnetClient().ListComplete(ctx, rgName)
}

func vnetClient() network.VirtualNetworksClient {
	c := network.NewVirtualNetworksClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorise Virtual Network client: %v", err)
	}
	return c
}

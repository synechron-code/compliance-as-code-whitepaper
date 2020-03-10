package network

import (
	"context"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-08-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

// CreateNIC creates a new network interface. The Network Security Group is not a required parameter.
func CreateNIC(ctx context.Context, vnetName, subnetName, nsgName, ipName, nicName string, tags map[string]*string) (nic network.Interface, err error) {
	subnet, err := GetVirtualNetworkSubnet(ctx, vnetName, subnetName)
	if err != nil {
		log.Fatalf("failed to get subnet: %v", err)
	}

	ip, err := PublicIP(ctx, ipName)
	if err != nil {
		log.Fatalf("failed to get ip address: %v", err)
	}

	nicParams := network.Interface{
		Name:     to.StringPtr(nicName),
		Location: to.StringPtr(azureutil.Location()),
		InterfacePropertiesFormat: &network.InterfacePropertiesFormat{
			IPConfigurations: &[]network.InterfaceIPConfiguration{
				{
					Name: to.StringPtr("ipConfig1"),
					InterfaceIPConfigurationPropertiesFormat: &network.InterfaceIPConfigurationPropertiesFormat{
						Subnet:                    &subnet,
						PrivateIPAllocationMethod: network.Dynamic,
						PublicIPAddress:           &ip,
					},
				},
			},
		},
		Tags: tags,
	}

	if nsgName != "" {
		nsg, err := SecurityGroup(ctx, nsgName)
		if err != nil {
			log.Fatalf("failed to get nsg: %v", err)
		}
		nicParams.NetworkSecurityGroup = &nsg
	}

	c := nicClient()
	future, err := c.CreateOrUpdate(ctx, azureutil.ResourceGroup(), nicName, nicParams)
	if err != nil {
		return nic, err
	}

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return nic, err
	}

	return future.Result(c)
}

// NIC returns an existing network interface by name
func NIC(ctx context.Context, name string) (network.Interface, error) {
	return nicClient().Get(ctx, azureutil.ResourceGroup(), name, "")
}

// DeleteNIC deletes an existing network interface by name.
func DeleteNIC(ctx context.Context, name string) (network.InterfacesDeleteFuture, error) {
	return nicClient().Delete(ctx, azureutil.ResourceGroup(), name)
}

func nicClient() network.InterfacesClient {
	c := network.NewInterfacesClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorise Network Interfaces client: %v", err)
	}
	return c
}

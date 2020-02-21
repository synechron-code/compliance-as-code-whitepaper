package network

import (
	"context"
	"fmt"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-08-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
	"net/http"
)

// VNet Subnets

func getSubnetsClient() network.SubnetsClient {
	subnetsClient := network.NewSubnetsClient(azureutil.GetAzureSubscriptionID())
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		subnetsClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return subnetsClient
}

// CreateVirtualNetworkSubnet creates a subnet in an existing vnet
func CreateVirtualNetworkSubnet(ctx context.Context, vnetName, subnetName string) (subnet network.Subnet, err error) {
	subnetsClient := getSubnetsClient()

	future, err := subnetsClient.CreateOrUpdate(
		ctx,
		azureutil.GetAzureResourceGP(),
		vnetName,
		subnetName,
		network.Subnet{
			SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
				AddressPrefix: to.StringPtr("10.0.0.0/16"),
			},
		})
	if err != nil {
		return subnet, fmt.Errorf("cannot create subnet: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, subnetsClient.Client)
	if err != nil {
		return subnet, fmt.Errorf("cannot get the subnet create or update future response: %v", err)
	}

	return future.Result(subnetsClient)
}

// CreateSubnetWithNetworkSecurityGroup create a subnet referencing a network security group
func CreateSubnetWithNetworkSecurityGroup(ctx context.Context, vnetName, subnetName, addressPrefix, nsgName string) (subnet network.Subnet, err error) {
	nsg, err := GetNetworkSecurityGroup(ctx, nsgName)
	if err != nil {
		return subnet, fmt.Errorf("cannot get nsg: %v", err)
	}

	subnetsClient := getSubnetsClient()
	future, err := subnetsClient.CreateOrUpdate(
		ctx,
		azureutil.GetAzureResourceGP(),
		vnetName,
		subnetName,
		network.Subnet{
			SubnetPropertiesFormat: &network.SubnetPropertiesFormat{
				AddressPrefix:        to.StringPtr(addressPrefix),
				NetworkSecurityGroup: &nsg,
			},
		})
	if err != nil {
		return subnet, fmt.Errorf("cannot create subnet: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, subnetsClient.Client)
	if err != nil {
		return subnet, fmt.Errorf("cannot get the subnet create or update future response: %v", err)
	}

	return future.Result(subnetsClient)
}

// DeleteVirtualNetworkSubnet deletes a subnet
func DeleteVirtualNetworkSubnet() {}

// GetVirtualNetworkSubnet returns an existing subnet from a virtual network
func GetVirtualNetworkSubnet(ctx context.Context, vnetName string, subnetName string) (network.Subnet, error) {
	subnetsClient := getSubnetsClient()
	return subnetsClient.Get(ctx, azureutil.GetAzureResourceGP(), vnetName, subnetName, "")
}

// GetVirtualNetworkSubnetByResourceGroup returns an existing subnet from a virtual network
func GetVirtualNetworkSubnetByResourceGroup(ctx context.Context, resourceGroup, vnetName, subnetName string) (network.Subnet, error) {
	subnetsClient := getSubnetsClient()
	return subnetsClient.Get(ctx, resourceGroup, vnetName, subnetName, "")
}

// GetSubnetPreparerWithID prepares the Get request.
func GetSubnetPreparerWithID(ctx context.Context, resourceID string, expand string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceId": resourceID,
	}

	subnetsClient := getSubnetsClient()

	const APIVersion = "2019-08-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(expand) > 0 {
		queryParameters["$expand"] = autorest.Encode("query", expand)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(subnetsClient.BaseURI),
		autorest.WithPathParameters("/{resourceId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))

	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSubnetByID gets the specified subnet by resourceID.
// Parameters:
// resourceID - resource ID of the subnet
// expand - expands referenced resources.
func GetSubnetByID(ctx context.Context, resourceID string, expand string) (result network.Subnet, err error) {
	subnetsClient := getSubnetsClient()
	req, err := GetSubnetPreparerWithID(ctx, resourceID, expand)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.SubnetsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := subnetsClient.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.SubnetsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = subnetsClient.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.SubnetsClient", "Get", resp, "Failure responding to request")
	}

	return
}

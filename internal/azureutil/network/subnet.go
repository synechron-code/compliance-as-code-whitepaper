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

// CreateVirtualNetworkSubnet creates a subnet in an existing vnet
func CreateVirtualNetworkSubnet(ctx context.Context, vnetName, subnetName string) (subnet network.Subnet, err error) {
	c := subnetsClient()

	future, err := c.CreateOrUpdate(
		ctx,
		azureutil.ResourceGroup(),
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

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return subnet, fmt.Errorf("cannot get the subnet create or update future response: %v", err)
	}

	return future.Result(c)
}

// CreateSubnetWithNetworkSecurityGroup create a subnet referencing a network security group
func CreateSubnetWithNetworkSecurityGroup(ctx context.Context, vnetName, subnetName, addressPrefix, nsgName string) (subnet network.Subnet, err error) {
	nsg, err := SecurityGroup(ctx, nsgName)
	if err != nil {
		return subnet, fmt.Errorf("cannot get nsg: %v", err)
	}

	c := subnetsClient()
	future, err := c.CreateOrUpdate(
		ctx,
		azureutil.ResourceGroup(),
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

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return subnet, fmt.Errorf("cannot get the subnet create or update future response: %v", err)
	}

	return future.Result(c)
}

// GetVirtualNetworkSubnet returns an existing subnet from a virtual network
func GetVirtualNetworkSubnet(ctx context.Context, vnetName string, subnetName string) (network.Subnet, error) {
	return subnetsClient().Get(ctx, azureutil.ResourceGroup(), vnetName, subnetName, "")
}

// GetVirtualNetworkSubnetByResourceGroup returns an existing subnet from a virtual network
func GetVirtualNetworkSubnetByResourceGroup(ctx context.Context, resourceGroup, vnetName, subnetName string) (network.Subnet, error) {
	return subnetsClient().Get(ctx, resourceGroup, vnetName, subnetName, "")
}

// GetSubnetPreparerWithID prepares the Get request.
func GetSubnetPreparerWithID(ctx context.Context, resourceID string, expand string) (*http.Request, error) {

	c := subnetsClient()

	queryParameters := map[string]interface{}{
		"api-version": "2019-08-01",
	}

	if len(expand) > 0 {
		queryParameters["$expand"] = autorest.Encode("query", expand)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(c.BaseURI),
		autorest.WithPathParameters("/{resourceId}", map[string]interface{}{
			"resourceId": resourceID,
		}),
		autorest.WithQueryParameters(queryParameters))

	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSubnetByID gets the specified subnet by resourceID.
// Parameters:
// resourceID - resource ID of the subnet
// expand - expands referenced resources.
func GetSubnetByID(ctx context.Context, resourceID string, expand string) (result network.Subnet, err error) {

	c := subnetsClient()

	req, err := GetSubnetPreparerWithID(ctx, resourceID, expand)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.SubnetsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := c.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.SubnetsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = c.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.SubnetsClient", "Get", resp, "Failure responding to request")
	}

	return
}

func subnetsClient() network.SubnetsClient {
	c := network.NewSubnetsClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return c
}

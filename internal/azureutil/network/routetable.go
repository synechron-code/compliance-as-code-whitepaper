package network

import (
	"context"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-08-01/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"net/http"
)

// VNet RouteTables

func getRouteTablesClient() network.RouteTablesClient {
	RouteTablesClient := network.NewRouteTablesClient(azureutil.GetAzureSubscriptionID())
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		RouteTablesClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return RouteTablesClient
}

// GetRouteTablePreparerWithID prepares the Get request.
func GetRouteTablePreparerWithID(ctx context.Context, resourceID string, expand string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceId": resourceID,
	}

	RouteTablesClient := getRouteTablesClient()

	const APIVersion = "2019-08-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if len(expand) > 0 {
		queryParameters["$expand"] = autorest.Encode("query", expand)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(RouteTablesClient.BaseURI),
		autorest.WithPathParameters("/{resourceId}", pathParameters),
		autorest.WithQueryParameters(queryParameters))

	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetRouteTableByID gets the specified RouteTable by resourceID.
// Parameters:
// resourceID - resource ID of the RouteTable
// expand - expands referenced resources.
func GetRouteTableByID(ctx context.Context, resourceID string, expand string) (result network.RouteTable, err error) {
	RouteTablesClient := getRouteTablesClient()
	req, err := GetRouteTablePreparerWithID(ctx, resourceID, expand)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.RouteTablesClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := RouteTablesClient.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.RouteTablesClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = RouteTablesClient.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.RouteTablesClient", "Get", resp, "Failure responding to request")
	}

	return
}

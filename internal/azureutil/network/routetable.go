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

// RouteTableByID gets the specified RouteTable by resourceID.
// Parameters:
// resourceID - resource ID of the RouteTable
// expand - expands referenced resources.
func RouteTableByID(ctx context.Context, resourceID string, expand string) (result network.RouteTable, err error) {

	c := routeTableClient()

	req, err := routeTablePreparerWithID(ctx, resourceID, expand)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.c", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := c.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "network.c", "Get", resp, "Failure sending request")
		return
	}

	result, err = c.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "network.c", "Get", resp, "Failure responding to request")
	}

	return
}

// GetRouteTablePreparerWithID prepares the Get request.
func routeTablePreparerWithID(ctx context.Context, resourceID string, expand string) (*http.Request, error) {

	c := routeTableClient()

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

func routeTableClient() network.RouteTablesClient {
	c := network.NewRouteTablesClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorise Route Table client: %v", err)
	}
	return c
}

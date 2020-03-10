package network

import (
	"citihub.com/compliance-as-code/internal/azureutil"
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-08-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// AzureFirewalls returns all Azure Firewall instances within the Subscription (configured by the AZURE_SUBSCRIPTION_ID environment variable).
func AzureFirewalls(ctx context.Context) (network.AzureFirewallListResultIterator, error) {
	c := fwClient()
	log.Printf("[DEBUG] subscriptionID: %v", c.SubscriptionID)
	r, err := c.ListAllComplete(ctx)
	if err == nil {
		log.Println("[DEBUG] Successfully listed all FW in subscription")
	}
	return r, err
}

func fwClient() network.AzureFirewallsClient {
	c := network.NewAzureFirewallsClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorise Azure Firewall client: %v", err)
	}
	return c
}

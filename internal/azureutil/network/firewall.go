package network

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-08-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func getFwClient(subscriptionID string) (fwClient network.AzureFirewallsClient, err error) {
	fwClient = network.NewAzureFirewallsClient(subscriptionID)
	// create an authorizer from env vars or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		fwClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return
}

// ListAllAzureFirewall return all AzureFirewall within the subscription
func ListAllAzureFirewall(ctx context.Context, fwSubscription string) (result network.AzureFirewallListResultIterator, err error) {
	fwClient, err := getFwClient(fwSubscription)
	log.Printf("subscriptionID: %v", fwClient.SubscriptionID)
	result, err = fwClient.ListAllComplete(ctx)
	if err == nil {
		log.Println("Successfully listed all FW in subscription")
	}
	return
}

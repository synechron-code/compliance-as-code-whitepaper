package containerservice

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-08-01/containerservice"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// Container Service

func getManagedClustersClient(aksSubscription string) containerservice.ManagedClustersClient {
	mcClient := containerservice.NewManagedClustersClient(aksSubscription)
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		mcClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}

	return mcClient
}

// ListAllAKS return all AKS clusters within the subscription
func ListAllAKS(ctx context.Context, aksSubscription string) (result containerservice.ManagedClusterListResultIterator, err error) {
	mcClient := getManagedClustersClient(aksSubscription)
	log.Printf("subscriptionID: %v", mcClient.SubscriptionID)
	result, err = mcClient.ListComplete(ctx)
	if err == nil {
		log.Println("Successfully listed all AKS in subscription")
	}
	return
}

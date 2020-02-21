package policy

import (
	"context"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-01-01/policy"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func getAssignmentsClient() policy.AssignmentsClient {
	assignmentsClient := policy.NewAssignmentsClient(azureutil.GetAzureSubscriptionID())
	// create an authorizer from env vars or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		assignmentsClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return assignmentsClient
}

// GetAssignmentBySubscription - Get a Policy definition according to name
func GetAssignmentBySubscription(ctx context.Context, subscriptionID, policyAssignmentName string) (policy.Assignment, error) {
	assignmentsClient := getAssignmentsClient()
	scope := "/subscriptions/" + subscriptionID
	log.Printf("Getting Policy Assignment with subscriptionID: %v", scope)
	return assignmentsClient.Get(ctx, scope, policyAssignmentName)
}

// GetAssignmentByManagementGroup - Get a Policy definition according to name in the management group
func GetAssignmentByManagementGroup(ctx context.Context, managementGroup, policyAssignmentName string) (policy.Assignment, error) {
	assignmentsClient := getAssignmentsClient()
	scope := "/providers/Microsoft.Management/managementGroups/" + managementGroup
	log.Printf("Getting Policy Assignment with scope: %v", scope)
	return assignmentsClient.Get(ctx, scope, policyAssignmentName)
}

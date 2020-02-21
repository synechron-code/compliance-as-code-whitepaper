package policy

import (
	"context"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-01-01/policy"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

func getPolicyClient() policy.DefinitionsClient {
	definitionClient := policy.NewDefinitionsClient(azureutil.GetAzureSubscriptionID())
	// create an authorizer from env vars or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		definitionClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return definitionClient
}

// GetPolicyDefinition - Get a Policy definition according to name
func GetPolicyDefinition(ctx context.Context, policyDefinitionName string) (policy.Definition, error) {
	definitionClient := getPolicyClient()

	return definitionClient.Get(ctx, policyDefinitionName)
}

package network

import (
	"context"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-08-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// Network Security Groups

func getNsrClient() network.SecurityRulesClient {
	nsrClient := network.NewSecurityRulesClient(azureutil.GetAzureSubscriptionID())
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		nsrClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}

	return nsrClient
}

// CreateNetworkSecurityRule creates a new network security rule
func CreateNetworkSecurityRule(ctx context.Context, networkSecurityGroupName string, securityRuleName string, securityRuleParameters network.SecurityRule) (nsr network.SecurityRule, err error) {
	nsrClient := getNsrClient()
	future, err := nsrClient.CreateOrUpdate(
		ctx,
		azureutil.GetAzureResourceGP(),
		networkSecurityGroupName,
		securityRuleName,
		securityRuleParameters)

	if err != nil {
		return nsr, err
	}

	err = future.WaitForCompletionRef(ctx, nsrClient.Client)
	if err != nil {
		return nsr, err
	}

	return future.Result(nsrClient)
}

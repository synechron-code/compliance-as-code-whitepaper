package network

import (
	"context"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-08-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// CreateSecurityRule creates a new network security rule
func CreateSecurityRule(ctx context.Context, nsgName string, nsrName string, parameters network.SecurityRule) (nsr network.SecurityRule, err error) {
	c := nsrClient()
	future, err := c.CreateOrUpdate(
		ctx,
		azureutil.ResourceGroup(),
		nsgName,
		nsrName,
		parameters)

	if err != nil {
		return nsr, err
	}

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return nsr, err
	}

	return future.Result(c)
}

func nsrClient() network.SecurityRulesClient {
	c := network.NewSecurityRulesClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorise Security Rules client: %v", err)
	}
	return c
}

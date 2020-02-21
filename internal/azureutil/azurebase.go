package azureutil

import (
	"log"
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
	EnvPolicyAssignmentManagementGroup string = "AZURE_POLICY_ASSIGNMENT_MANAGEMENT_GROUP"
)

var testprefix string
var azureResourceGp string

// GetTestPrefix return a random test prefix with test + 6 random characters
func GetTestPrefix() string {
	if testprefix == "" {
		testprefix = "test" + RandStringBytesMaskImprSrcUnsafe(6) + ""
	}
	return testprefix
}

func getFromEnvVar(varName string) string {
	result := os.Getenv(varName)
	if result == "" {
		log.Fatalf("Environment variable \"%v\" is not defined", varName)
	}
	return result
}

//GetAzureResourceGP - Default resource GP
func GetAzureResourceGP() string {
	if azureResourceGp == "" {
		azureResourceGp = GetTestPrefix() + "resourceGP"
	}
	return azureResourceGp
}

//GetAzureLocation - Default location
func GetAzureLocation() string {
	return getFromEnvVar("AZURE_LOCATION")
}

//GetAzureSubscriptionID - Return subscriptionID
func GetAzureSubscriptionID() string {
	return getFromEnvVar("AZURE_SUBSCRIPTION_ID")
}

//GetAzureAuthorizer - return an Azure Authorizer
func GetAzureAuthorizer() autorest.Authorizer {
	// create an authorizer from env vars or Azure Managed Service Identity

	Authorizer, err := auth.NewAuthorizerFromEnvironment()

	if err != nil {
		log.Panicf("Unable to load Azure credential due to %v", err)
	}
	return Authorizer
}

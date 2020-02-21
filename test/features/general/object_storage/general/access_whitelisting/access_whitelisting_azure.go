package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"citihub.com/compliance-as-code/internal/azureutil"
	"citihub.com/compliance-as-code/internal/azureutil/policy"
	"citihub.com/compliance-as-code/internal/azureutil/resource"
	"citihub.com/compliance-as-code/internal/azureutil/storage"
	azurePolicy "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-01-01/policy"
	azureStorage "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
)

const (
	policyAssignmentName = "deny_storage_wo_net_acl"
)

const storageRgEnvVar = "STORAGE_ACCOUNT_RESOURCE_GROUP"

// AccessWhitelistingAzure azure implementation of the encryption in flight for Object Storage feature
type AccessWhitelistingAzure struct {
	ctx                       context.Context
	policyAssignmentMgmtGroup string
	tags                      map[string]*string
	bucketName                string
	storageAccount            azureStorage.Account
	runningErr                error
}

func (state *AccessWhitelistingAzure) setup() {
	log.Println("Setting up \"AccessWhitelistingAzure\"")
	state.ctx = context.Background()
	state.policyAssignmentMgmtGroup = os.Getenv(azureutil.EnvPolicyAssignmentManagementGroup)
	if state.policyAssignmentMgmtGroup == "" {
		log.Printf("'%v' environment variable is not defined. Policy assignment check against subscription", azureutil.EnvPolicyAssignmentManagementGroup)
	}

	state.tags = map[string]*string{
		"project": to.StringPtr("gitlab-CICD"),
		"env":     to.StringPtr("test"),
		"tier":    to.StringPtr("internal"),
	}

	_, err := resource.CreateGroupWithTags(state.ctx, azureutil.GetAzureResourceGP(), state.tags)

	if err != nil {
		log.Fatalf("failed to create group: %v\n", err.Error())
	}
	log.Printf("Created Resource Group: %v", azureutil.GetAzureResourceGP())
}

func (state *AccessWhitelistingAzure) teardown() {
	err := resource.Cleanup(state.ctx)
	if err != nil {
		log.Fatalf("Failed to teardown: %v\n", err.Error())
	}
	log.Println("Teardown completed")
}

func (state *AccessWhitelistingAzure) checkPolicyAssigned() error {
	var policyAssignment azurePolicy.Assignment
	var aerr error
	// Search assignment from Management Group instead of subscription
	if state.policyAssignmentMgmtGroup != "" {
		policyAssignment, aerr = policy.GetAssignmentByManagementGroup(state.ctx, state.policyAssignmentMgmtGroup, policyAssignmentName)
	} else {
		policyAssignment, aerr = policy.GetAssignmentBySubscription(state.ctx, azureutil.GetAzureSubscriptionID(), policyAssignmentName)
	}

	if aerr != nil {
		log.Printf("Get policy assignment error: %v", aerr)
		return aerr
	}

	log.Printf("Policy assignment check: %v [Step PASSED]", *policyAssignment.Name)
	return nil
}

func (state *AccessWhitelistingAzure) prepareToCreateStorageContainer() error {
	state.bucketName = azureutil.RandStringBytesMaskImprSrcUnsafe(10)
	return nil
}

func (state *AccessWhitelistingAzure) createWithWhiteList(ipRange string) error {
	var networkRuleSet azureStorage.NetworkRuleSet
	if ipRange == "nil" {
		networkRuleSet = azureStorage.NetworkRuleSet{
			DefaultAction: azureStorage.DefaultActionAllow,
		}
	} else {
		// ipRule := &azureStorage.IPRule{
		// 	Action: azureStorage.Allow,
		// 	IPAddressOrRange: to.StringPtr(ipRange),
		// }
		var ipRules *[]azureStorage.IPRule

		networkRuleSet = azureStorage.NetworkRuleSet{
			IPRules:       ipRules,
			DefaultAction: azureStorage.DefaultActionDeny,
		}
	}

	state.storageAccount, state.runningErr = storage.CreateStorageAccountWithNetworkRuleSet(state.ctx, state.bucketName, azureutil.GetAzureResourceGP(), state.tags, true, &networkRuleSet)

	return nil
}

func (state *AccessWhitelistingAzure) creationWill(result string) error {
	if result == "Fail" {
		if state.runningErr == nil {
			return fmt.Errorf("incorrectly created storage account: %v", *state.storageAccount.ID)
		}
		return nil
	}

	if state.runningErr == nil {
		return nil
	}
	return state.runningErr
}

func (state *AccessWhitelistingAzure) isCspCapable() error {
	return nil
}

func (state *AccessWhitelistingAzure) examineStorageContainer(containerNameEnvVar string) error {
	accountName := os.Getenv(containerNameEnvVar)
	if accountName == "" {
		return fmt.Errorf("environment variable \"%s\" is not defined test can't run", containerNameEnvVar)
	}

	resourceGroup := os.Getenv(storageRgEnvVar)
	if resourceGroup == "" {
		return fmt.Errorf("environment variable \"%s\" is not defined test can't run", storageRgEnvVar)
	}

	state.storageAccount, state.runningErr = storage.GetStorageAccountProperties(state.ctx, resourceGroup, accountName)

	if state.runningErr != nil {
		return state.runningErr
	}

	networkRuleSet := state.storageAccount.AccountProperties.NetworkRuleSet
	result := false
	// Default action is deny
	if networkRuleSet.DefaultAction == azureStorage.DefaultActionAllow {
		return fmt.Errorf("%s has not configured with firewall network rule default action is not deny", accountName)
	}

	// Check if it has IP white listing
	for _, ipRule := range *networkRuleSet.IPRules {
		result = true
		log.Printf("IP WhiteListing: %v, %v", *ipRule.IPAddressOrRange, ipRule.Action)
	}

	// Check if it has private Endpoint white listing
	for _, vnetRule := range *networkRuleSet.VirtualNetworkRules {
		result = true
		log.Printf("VNet whitelisting: %v, %v", *vnetRule.VirtualNetworkResourceID, vnetRule.Action)
	}

	// TODO: Private Endpoint implementation when it's GA

	if result {
		log.Printf("Whitelisting rule exist. [Step PASSED]")
		return nil
	}
	return fmt.Errorf("no whitelisting has been defined for %v", accountName)
}

func (state *AccessWhitelistingAzure) whitelistingIsConfigured() error {
	// Checked in previous step
	return nil
}

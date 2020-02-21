package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"citihub.com/compliance-as-code/internal/azureutil"
	"citihub.com/compliance-as-code/internal/azureutil/policy"
	"citihub.com/compliance-as-code/internal/azureutil/resource"
	"citihub.com/compliance-as-code/internal/azureutil/storage"
	azurePolicy "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-01-01/policy"
	azureStorage "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
)

const (
	policyName = "deny_http_storage"
)

// EncryptionInFlightAzure azure implementation of the encryption in flight for Object Storage feature
type EncryptionInFlightAzure struct {
	ctx                       context.Context
	tags                      map[string]*string
	httpOption                bool
	httpsOption               bool
	policyAssignmentMgmtGroup string
}

func (state *EncryptionInFlightAzure) setup() {
	log.Println("Setting up \"EncryptionInFlightAzure\"")
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

func (state *EncryptionInFlightAzure) teardown() {
	resource.Cleanup(state.ctx)
	log.Println("Teardown completed")
}

func (state *EncryptionInFlightAzure) securityControlsThatRestrictDataFromBeingUnencryptedInFlight() error {
	var policyAssignment azurePolicy.Assignment
	var aerr error
	// Search assignment from Management Group instead of subscription
	if state.policyAssignmentMgmtGroup != "" {
		policyAssignment, aerr = policy.GetAssignmentByManagementGroup(state.ctx, state.policyAssignmentMgmtGroup, policyName)
	} else {
		policyAssignment, aerr = policy.GetAssignmentBySubscription(state.ctx, azureutil.GetAzureSubscriptionID(), policyName)
	}

	if aerr != nil {
		log.Printf("Get policy assignment error: %v", aerr)
		return aerr
	}

	log.Printf("Policy assignment check: %v [Step PASSED]", *policyAssignment.Name)
	return nil
}

func (state *EncryptionInFlightAzure) weProvisionAnObjectStorageBucket() error {
	// Nothing to do here
	return nil
}

func (state *EncryptionInFlightAzure) httpAccessIs(arg1 string) error {
	if arg1 == "enabled" {
		state.httpOption = true
	} else {
		state.httpOption = false
	}
	return nil
}

func (state *EncryptionInFlightAzure) httpsAccessIs(arg1 string) error {
	if arg1 == "enabled" {
		state.httpsOption = true
	} else {
		state.httpsOption = false
	}
	return nil
}

func (state *EncryptionInFlightAzure) creationWillWithAnErrorMatching(expectation, errDescription string) error {
	accountName := azureutil.RandStringBytesMaskImprSrcUnsafe(5) + "storageac"

	var err error

	networkRuleSet := azureStorage.NetworkRuleSet{
		DefaultAction: azureStorage.DefaultActionDeny,
	}

	// Both true take it as http option is try
	if state.httpsOption && state.httpOption {
		_, err = storage.CreateStorageAccountWithNetworkRuleSet(state.ctx, accountName,
			azureutil.GetAzureResourceGP(), state.tags, false, &networkRuleSet)
	} else if state.httpsOption {
		_, err = storage.CreateStorageAccountWithNetworkRuleSet(state.ctx, accountName,
			azureutil.GetAzureResourceGP(), state.tags, state.httpsOption, &networkRuleSet)
	} else if state.httpOption {
		_, err = storage.CreateStorageAccountWithNetworkRuleSet(state.ctx, accountName,
			azureutil.GetAzureResourceGP(), state.tags, state.httpsOption, &networkRuleSet)
	}

	if expectation == "Fail" {

		if err == nil {
			return fmt.Errorf("Storage Account was created, but should not have been: policy is not working or incorrectly configured")
		}

		detailedError := err.(autorest.DetailedError)
		originalErr := detailedError.Original
		detailed := originalErr.(*azure.ServiceError)

		log.Printf("Detailed Error: %v", detailed)
		
		if strings.EqualFold(detailed.Code, "RequestDisallowedByPolicy") {
			// Now check if it is the right policy
			if strings.Contains(detailed.Message, policyName) {
				log.Printf("Request was Disallowed By Policy: %v [Step PASSED]", policyName)
				return nil
			}
			return fmt.Errorf("Storage Account was not created but blocked not by the right policy: %v", detailed.Message)
		} 
		
		return fmt.Errorf("Storage Account was not created")
	} else if expectation == "Succeed" {
		if err != nil {
			log.Printf("Unexpected failure in create storage ac [Step FAILED]")
			return err
		}
		return nil
	}

	return fmt.Errorf("unsupported `result` option '%s' in the Gherkin feature - use either 'Fail' or 'Succeed'", expectation)
}

func (state *EncryptionInFlightAzure) cSPProvideDetectiveMeasureForNonComplianceSecureTransferOnObjectStorage() error {
	return nil
}

func (state *EncryptionInFlightAzure) weExamineTheDetectiveMeasure() error {
	return nil
}

func (state *EncryptionInFlightAzure) theDetectiveMeasureIsEnabled() error {
	var policyAssignment azurePolicy.Assignment
	var aerr error
	// Search assignment from Management Group instead of subscription
	if state.policyAssignmentMgmtGroup != "" {
		policyAssignment, aerr = policy.GetAssignmentByManagementGroup(state.ctx, state.policyAssignmentMgmtGroup, policyName)
	} else {
		policyAssignment, aerr = policy.GetAssignmentBySubscription(state.ctx, azureutil.GetAzureSubscriptionID(), policyName)
	}

	if aerr != nil {
		log.Printf("Get policy assignment error: %v", aerr)
		return aerr
	}

	log.Printf("Policy assignment check: %v [Step PASSED]", *policyAssignment.Name)
	return nil
}

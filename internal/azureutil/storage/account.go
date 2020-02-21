package storage

import (
	"context"
	"fmt"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

func getStorageAccountsClient() storage.AccountsClient {
	storageAccountsClient := storage.NewAccountsClient(azureutil.GetAzureSubscriptionID())
	// create an authorizer from env vars or Azure Managed Service Identity
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		storageAccountsClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return storageAccountsClient
}

// CreateStorageAccount starts creation of a new storage account and waits for
// the account to be created.
func CreateStorageAccount(ctx context.Context, accountName, accountGroupName string, tags map[string]*string) (storage.Account, error) {
	return CreateStorageAccountWithHTTPSOption(ctx, accountName, accountGroupName, tags, true)
}

// CreateStorageAccountWithHTTPSOption starts creation of a new storage account and waits for
// the account to be created.
func CreateStorageAccountWithHTTPSOption(ctx context.Context, accountName, accountGroupName string, tags map[string]*string, httpsOnly bool) (storage.Account, error) {
	return CreateStorageAccountWithNetworkRuleSet(ctx, accountName, accountGroupName, tags, httpsOnly, &storage.NetworkRuleSet{})
}

// CreateStorageAccountWithNetworkRuleSet starts creation of a new storage account and waits for
// the account to be created.
func CreateStorageAccountWithNetworkRuleSet(ctx context.Context, accountName, accountGroupName string, tags map[string]*string, httpsOnly bool, networkRuleSet *storage.NetworkRuleSet) (storage.Account, error) {
	var s storage.Account
	storageAccountsClient := getStorageAccountsClient()

	result, err := storageAccountsClient.CheckNameAvailability(
		ctx,
		storage.AccountCheckNameAvailabilityParameters{
			Name: to.StringPtr(accountName),
			Type: to.StringPtr("Microsoft.Storage/storageAccounts"),
		})
	if err != nil {
		return s, err
	}

	if *result.NameAvailable != true {
		return s, fmt.Errorf(
			"storage account name [%s] not available: %v\nserver message: %v",
			accountName, err, *result.Message)
	}

	networkRuleSetParam := &storage.AccountPropertiesCreateParameters{
		EnableHTTPSTrafficOnly: to.BoolPtr(httpsOnly),
		NetworkRuleSet:         networkRuleSet,
	}

	future, err := storageAccountsClient.Create(
		ctx,
		accountGroupName,
		accountName,
		storage.AccountCreateParameters{
			Sku: &storage.Sku{
				Name: storage.StandardLRS},
			Kind:                              storage.Storage,
			Location:                          to.StringPtr(azureutil.GetAzureLocation()),
			AccountPropertiesCreateParameters: networkRuleSetParam,
			Tags:                              tags,
		})

	if err != nil {
		return s, err
	}

	err = future.WaitForCompletionRef(ctx, storageAccountsClient.Client)
	if err != nil {
		return s, err
	}

	return future.Result(storageAccountsClient)
}

// DeleteStorageAccount deletes an existing storage account
func DeleteStorageAccount(ctx context.Context, accountName, accountGroupName string) (autorest.Response, error) {
	storageAccountsClient := getStorageAccountsClient()
	return storageAccountsClient.Delete(ctx, accountGroupName, accountName)
}

// GetAccountKeys gets the storage account keys
func GetAccountKeys(ctx context.Context, accountName, accountGroupName string) (storage.AccountListKeysResult, error) {
	accountsClient := getStorageAccountsClient()
	return accountsClient.ListKeys(ctx, accountGroupName, accountName, "")
}

// GetAccountPrimaryKey return the primary key
func GetAccountPrimaryKey(ctx context.Context, accountName, accountGroupName string) string {
	response, err := GetAccountKeys(ctx, accountName, accountGroupName)
	if err != nil {
		log.Fatalf("failed to list keys: %v", err)
	}
	return *(((*response.Keys)[0]).Value)
}

// GetStorageAccountProperties - return the properties of the storageAccounts
func GetStorageAccountProperties(ctx context.Context, resourceGroupName, accountName string) (storage.Account, error) {
	accountsClient := getStorageAccountsClient()
	return accountsClient.GetProperties(ctx, resourceGroupName, accountName, "")
}

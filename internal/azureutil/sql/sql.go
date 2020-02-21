package sql

import (
	"context"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/2015-05-01-preview/sql"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

// Servers

func getServersClient() sql.ServersClient {
	serversClient := sql.NewServersClient(azureutil.GetAzureSubscriptionID())
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		serversClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return serversClient
}

// CreateServer creates a new SQL Server
func CreateServer(ctx context.Context, resourceGPName, serverName, dbLogin, dbPassword string, tags map[string]*string) (server sql.Server, err error) {
	serversClient := getServersClient()
	future, err := serversClient.CreateOrUpdate(
		ctx,
		resourceGPName,
		serverName,
		sql.Server{
			Location: to.StringPtr(azureutil.GetAzureLocation()),
			ServerProperties: &sql.ServerProperties{
				AdministratorLogin:         to.StringPtr(dbLogin),
				AdministratorLoginPassword: to.StringPtr(dbPassword),
			},
			Tags: tags,
		})

	if err != nil {
		return server, err
	}

	err = future.WaitForCompletionRef(ctx, serversClient.Client)
	if err != nil {
		return server, err
	}

	return future.Result(serversClient)
}

// Databases

func getDbClient() sql.DatabasesClient {
	dbClient := sql.NewDatabasesClient(azureutil.GetAzureSubscriptionID())
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		dbClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return dbClient
}

// CreateDB creates a new SQL Database on a given server
func CreateDB(ctx context.Context, resourceGPName, serverName, dbName string) (db sql.Database, err error) {
	dbClient := getDbClient()
	future, err := dbClient.CreateOrUpdate(
		ctx,
		resourceGPName,
		serverName,
		dbName,
		sql.Database{
			Location: to.StringPtr(azureutil.GetAzureLocation()),
		})
	if err != nil {
		return db, err
	}

	err = future.WaitForCompletionRef(ctx, dbClient.Client)
	if err != nil {
		return db, err
	}

	return future.Result(dbClient)
}

// DeleteDB deletes an existing database from a server
func DeleteDB(ctx context.Context, resourceGPName, serverName, dbName string) (autorest.Response, error) {
	dbClient := getDbClient()
	return dbClient.Delete(
		ctx,
		resourceGPName,
		serverName,
		dbName,
	)
}

// Firewall rules

func getFwRulesClient() sql.FirewallRulesClient {
	fwrClient := sql.NewFirewallRulesClient(azureutil.GetAzureSubscriptionID())
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		fwrClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return fwrClient
}

// CreateFirewallRules creates new firewall rules for a given server
func CreateFirewallRules(ctx context.Context, resourceGPName, serverName string) error {
	fwrClient := getFwRulesClient()

	_, err := fwrClient.CreateOrUpdate(
		ctx,
		resourceGPName,
		serverName,
		"unsafe open to world",
		sql.FirewallRule{
			FirewallRuleProperties: &sql.FirewallRuleProperties{
				StartIPAddress: to.StringPtr("0.0.0.0"),
				EndIPAddress:   to.StringPtr("255.255.255.255"),
			},
		},
	)
	if err != nil {
		return err
	}

	_, err = fwrClient.CreateOrUpdate(
		ctx,
		resourceGPName,
		serverName,
		"open to Azure internal",
		sql.FirewallRule{
			FirewallRuleProperties: &sql.FirewallRuleProperties{
				StartIPAddress: to.StringPtr("0.0.0.0"),
				EndIPAddress:   to.StringPtr("0.0.0.0"),
			},
		},
	)

	return err
}

// PrintInfo logs information on SQL user agent and ARM client
func PrintInfo() {
	log.Printf("user agent string: %s\n", sql.UserAgent())
	log.Printf("SQL ARM Client version: %s\n", sql.Version())
}

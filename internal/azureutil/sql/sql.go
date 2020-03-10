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

//Servers

// CreateServer creates or updates a SQL Server instance and waits for request completion.
func CreateServer(ctx context.Context, rgName, serverName, dbLogin, dbPassword string, tags map[string]*string) (server sql.Server, err error) {
	c := serverClient()
	future, err := c.CreateOrUpdate(
		ctx,
		rgName,
		serverName,
		sql.Server{
			Location: to.StringPtr(azureutil.Location()),
			ServerProperties: &sql.ServerProperties{
				AdministratorLogin:         to.StringPtr(dbLogin),
				AdministratorLoginPassword: to.StringPtr(dbPassword),
			},
			Tags: tags,
		})

	if err != nil {
		return server, err
	}

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return server, err
	}

	return future.Result(c)
}

func serverClient() sql.ServersClient {
	c := sql.NewServersClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorise SQL Server client: %v", err)
	}
	return c
}

// Databases

// CreateServer creates or updates a SQL Database instance on the given server and waits for request completion.
func CreateDB(ctx context.Context, rgName, serverName, dbName string) (db sql.Database, err error) {
	c := dbClient()
	future, err := c.CreateOrUpdate(
		ctx,
		rgName,
		serverName,
		dbName,
		sql.Database{
			Location: to.StringPtr(azureutil.Location()),
		})
	if err != nil {
		return db, err
	}

	err = future.WaitForCompletionRef(ctx, c.Client)
	if err != nil {
		return db, err
	}

	return future.Result(c)
}

// DeleteDB deletes an existing database from a server.
func DeleteDB(ctx context.Context, rgName, serverName, dbName string) (autorest.Response, error) {
	return dbClient().Delete(ctx, rgName, serverName, dbName)
}

func dbClient() sql.DatabasesClient {
	c := sql.NewDatabasesClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorise SQL Database client: %v", err)
	}
	return c
}

// Firewall rules

// CreateFirewallRules creates or updates two SQL Firewall Rules (open to world and open to the Azure network)
func CreateFirewallRules(ctx context.Context, rgName, serverName string) error {
	c := fwRulesClient()

	_, err := c.CreateOrUpdate(
		ctx,
		rgName,
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

	_, err = c.CreateOrUpdate(
		ctx,
		rgName,
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

func fwRulesClient() sql.FirewallRulesClient {
	c := sql.NewFirewallRulesClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorise SQL Firewall client: %v", err)
	}
	return c
}

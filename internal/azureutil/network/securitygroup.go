package network

import (
	"context"
	"fmt"
	"log"

	"citihub.com/compliance-as-code/internal/azureutil"
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-08-01/network"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/Azure/go-autorest/autorest/to"
)

// Network Security Groups

func getNsgClient() network.SecurityGroupsClient {
	nsgClient := network.NewSecurityGroupsClient(azureutil.GetAzureSubscriptionID())
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		nsgClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}

	return nsgClient
}

// CreateNetworkSecurityGroup creates a new network security group with rules set for allowing SSH and HTTPS use
func CreateNetworkSecurityGroup(ctx context.Context, nsgName string, tags map[string]*string) (nsg network.SecurityGroup, err error) {
	nsgClient := getNsgClient()
	future, err := nsgClient.CreateOrUpdate(
		ctx,
		azureutil.GetAzureResourceGP(),
		nsgName,
		network.SecurityGroup{
			Location: to.StringPtr(azureutil.GetAzureLocation()),
			SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
				SecurityRules: &[]network.SecurityRule{
					{
						Name: to.StringPtr("allow_ssh"),
						SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
							Protocol:                 network.SecurityRuleProtocolTCP,
							SourceAddressPrefix:      to.StringPtr("0.0.0.0/0"),
							SourcePortRange:          to.StringPtr("1-65535"),
							DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
							DestinationPortRange:     to.StringPtr("22"),
							Access:                   network.SecurityRuleAccessAllow,
							Direction:                network.SecurityRuleDirectionInbound,
							Priority:                 to.Int32Ptr(100),
						},
					},
					{
						Name: to.StringPtr("allow_https"),
						SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
							Protocol:                 network.SecurityRuleProtocolTCP,
							SourceAddressPrefix:      to.StringPtr("0.0.0.0/0"),
							SourcePortRange:          to.StringPtr("1-65535"),
							DestinationAddressPrefix: to.StringPtr("0.0.0.0/0"),
							DestinationPortRange:     to.StringPtr("443"),
							Access:                   network.SecurityRuleAccessAllow,
							Direction:                network.SecurityRuleDirectionInbound,
							Priority:                 to.Int32Ptr(200),
						},
					},
				},
			},
			Tags: tags,
		},
	)

	if err != nil {
		return nsg, fmt.Errorf("cannot create nsg: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, nsgClient.Client)
	if err != nil {
		return nsg, fmt.Errorf("cannot get nsg create or update future response: %v", err)
	}

	return future.Result(nsgClient)
}

// CreateCustomNetworkSecurityGroup creates a new network security group with rules specified in 3rd argument
func CreateCustomNetworkSecurityGroup(ctx context.Context, nsgName string, securityRules []network.SecurityRule) (nsg network.SecurityGroup, err error) {
	return CreateCustomNetworkSecurityGroupWithTags(ctx, nsgName, securityRules, nil)
}

// CreateCustomNetworkSecurityGroupWithTags creates a new network security group with rules specified in 3rd argument
func CreateCustomNetworkSecurityGroupWithTags(ctx context.Context, nsgName string, securityRules []network.SecurityRule, tags map[string]*string) (nsg network.SecurityGroup, err error) {
	nsgClient := getNsgClient()
	future, err := nsgClient.CreateOrUpdate(
		ctx,
		azureutil.GetAzureResourceGP(),
		nsgName,
		network.SecurityGroup{
			Location: to.StringPtr(azureutil.GetAzureLocation()),
			SecurityGroupPropertiesFormat: &network.SecurityGroupPropertiesFormat{
				SecurityRules: &securityRules,
			},
			Tags: tags,
		},
	)

	if err != nil {
		return nsg, err
	}

	err = future.WaitForCompletionRef(ctx, nsgClient.Client)
	if err != nil {
		return nsg, err
	}

	return future.Result(nsgClient)
}

// CreateSimpleNetworkSecurityGroup creates a new network security group, without rules (rules can be set later)
func CreateSimpleNetworkSecurityGroup(ctx context.Context, nsgName string) (nsg network.SecurityGroup, err error) {
	nsgClient := getNsgClient()
	future, err := nsgClient.CreateOrUpdate(
		ctx,
		azureutil.GetAzureResourceGP(),
		nsgName,
		network.SecurityGroup{
			Location: to.StringPtr(azureutil.GetAzureLocation()),
		},
	)

	if err != nil {
		return nsg, fmt.Errorf("cannot create nsg: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, nsgClient.Client)
	if err != nil {
		return nsg, fmt.Errorf("cannot get nsg create or update future response: %v", err)
	}

	return future.Result(nsgClient)
}

// DeleteNetworkSecurityGroup deletes an existing network security group
func DeleteNetworkSecurityGroup(ctx context.Context, nsgName string) (result network.SecurityGroupsDeleteFuture, err error) {
	nsgClient := getNsgClient()
	return nsgClient.Delete(ctx, azureutil.GetAzureResourceGP(), nsgName)
}

// GetNetworkSecurityGroup returns an existing network security group
func GetNetworkSecurityGroup(ctx context.Context, nsgName string) (network.SecurityGroup, error) {
	nsgClient := getNsgClient()
	return nsgClient.Get(ctx, azureutil.GetAzureResourceGP(), nsgName, "")
}

// Network security group rules

func getSecurityRulesClient() network.SecurityRulesClient {
	rulesClient := network.NewSecurityRulesClient(azureutil.GetAzureSubscriptionID())
	authorizer, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		rulesClient.Authorizer = authorizer
	} else {
		log.Fatalf("Unable to get Authorization: %v", err)
	}
	return rulesClient
}

// CreateSSHRule creates an inbound network security rule that allows using port 22
func CreateSSHRule(ctx context.Context, nsgName string) (rule network.SecurityRule, err error) {
	rulesClient := getSecurityRulesClient()
	future, err := rulesClient.CreateOrUpdate(ctx,
		azureutil.GetAzureResourceGP(),
		nsgName,
		"ALLOW-SSH",
		network.SecurityRule{
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Access:                   network.SecurityRuleAccessAllow,
				DestinationAddressPrefix: to.StringPtr("*"),
				DestinationPortRange:     to.StringPtr("22"),
				Direction:                network.SecurityRuleDirectionInbound,
				Description:              to.StringPtr("Allow SSH"),
				Priority:                 to.Int32Ptr(103),
				Protocol:                 network.SecurityRuleProtocolTCP,
				SourceAddressPrefix:      to.StringPtr("*"),
				SourcePortRange:          to.StringPtr("*"),
			},
		})
	if err != nil {
		return rule, fmt.Errorf("cannot create SSH security rule: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, rulesClient.Client)
	if err != nil {
		return rule, fmt.Errorf("cannot get security rule create or update future response: %v", err)
	}

	return future.Result(rulesClient)
}

// CreateHTTPRule creates an inbound network security rule that allows using port 80
func CreateHTTPRule(ctx context.Context, nsgName string) (rule network.SecurityRule, err error) {
	rulesClient := getSecurityRulesClient()
	future, err := rulesClient.CreateOrUpdate(ctx,
		azureutil.GetAzureResourceGP(),
		nsgName,
		"ALLOW-HTTP",
		network.SecurityRule{
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Access:                   network.SecurityRuleAccessAllow,
				DestinationAddressPrefix: to.StringPtr("*"),
				DestinationPortRange:     to.StringPtr("80"),
				Direction:                network.SecurityRuleDirectionInbound,
				Description:              to.StringPtr("Allow HTTP"),
				Priority:                 to.Int32Ptr(101),
				Protocol:                 network.SecurityRuleProtocolTCP,
				SourceAddressPrefix:      to.StringPtr("*"),
				SourcePortRange:          to.StringPtr("*"),
			},
		})
	if err != nil {
		return rule, fmt.Errorf("cannot create HTTP security rule: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, rulesClient.Client)
	if err != nil {
		return rule, fmt.Errorf("cannot get security rule create or update future response: %v", err)
	}

	return future.Result(rulesClient)
}

// CreateSQLRule creates an inbound network security rule that allows using port 1433
func CreateSQLRule(ctx context.Context, nsgName, frontEndAddressPrefix string) (rule network.SecurityRule, err error) {
	rulesClient := getSecurityRulesClient()
	future, err := rulesClient.CreateOrUpdate(ctx,
		azureutil.GetAzureResourceGP(),
		nsgName,
		"ALLOW-SQL",
		network.SecurityRule{
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Access:                   network.SecurityRuleAccessAllow,
				DestinationAddressPrefix: to.StringPtr("*"),
				DestinationPortRange:     to.StringPtr("1433"),
				Direction:                network.SecurityRuleDirectionInbound,
				Description:              to.StringPtr("Allow SQL"),
				Priority:                 to.Int32Ptr(102),
				Protocol:                 network.SecurityRuleProtocolTCP,
				SourceAddressPrefix:      &frontEndAddressPrefix,
				SourcePortRange:          to.StringPtr("*"),
			},
		})
	if err != nil {
		return rule, fmt.Errorf("cannot create SQL security rule: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, rulesClient.Client)
	if err != nil {
		return rule, fmt.Errorf("cannot get security rule create or update future response: %v", err)
	}

	return future.Result(rulesClient)
}

// CreateDenyOutRule creates an network security rule that denies outbound traffic
func CreateDenyOutRule(ctx context.Context, nsgName string) (rule network.SecurityRule, err error) {
	rulesClient := getSecurityRulesClient()
	future, err := rulesClient.CreateOrUpdate(ctx,
		azureutil.GetAzureResourceGP(),
		nsgName,
		"DENY-OUT",
		network.SecurityRule{
			SecurityRulePropertiesFormat: &network.SecurityRulePropertiesFormat{
				Access:                   network.SecurityRuleAccessDeny,
				DestinationAddressPrefix: to.StringPtr("*"),
				DestinationPortRange:     to.StringPtr("*"),
				Direction:                network.SecurityRuleDirectionOutbound,
				Description:              to.StringPtr("Deny outbound traffic"),
				Priority:                 to.Int32Ptr(100),
				Protocol:                 network.SecurityRuleProtocolAsterisk,
				SourceAddressPrefix:      to.StringPtr("*"),
				SourcePortRange:          to.StringPtr("*"),
			},
		})
	if err != nil {
		return rule, fmt.Errorf("cannot create deny out security rule: %v", err)
	}

	err = future.WaitForCompletionRef(ctx, rulesClient.Client)
	if err != nil {
		return rule, fmt.Errorf("cannot get security rule create or update future response: %v", err)
	}

	return future.Result(rulesClient)
}

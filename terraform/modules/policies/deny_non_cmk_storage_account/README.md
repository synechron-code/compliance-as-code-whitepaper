# Deny creating storage account using Microsoft managed Key

Deny creating storage account using microsoft managed key.

## Cloud Controls Objectives

This policy help to satisfy the following Common Control Objectives:

| Controls ID  | Objectives |
|---|---|
|AGP135|Ensure encryption keys are owned and managed by the FI following industry best practice for key management|

## Intended Use

Prevent creating storage without customer managed key. This policy is defaulted to audit as the portal does not yet support creating storage account with CMK in one step. This results in a chicken and egg scenario.

Moreover, Terraform seems to have its own problems when it comes to creating a storage account that is using a customer managed key; see <https://github.com/terraform-providers/terraform-provider-azurerm/issues/658> for details.

### Variables

management_group_id : the management group Id that the policy definition is created against.

## Apply with Terraform

This should be applied to Azure as a policy and then assigned with appropriate parameters. This would be applied with the main azure-policy module.
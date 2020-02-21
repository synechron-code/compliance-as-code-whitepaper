// Policy Definition
resource "azurerm_policy_definition" "deny_non_cmk_storage_ac" {
  name                = "deny_non_cmk_storage_ac"
  policy_type         = "Custom"
  mode                = "Indexed"
  display_name        = "Deny storage account using MS Managed Key [BDD]"
  description         = "Deny storage account using Microsoft Managed Key. Storage account should be using Customer Managed Key from keyvault"
  management_group_id = var.definition_management_group_id
  metadata            = <<METADATA
  {
    "category": "Storage"
  }
  METADATA

  lifecycle {
    ignore_changes = [
      metadata
    ]
  }

  parameters = <<PARAMETERS
  {
    "effect": {
        "type": "String",
        "metadata": {
          "displayName": "Effect",
          "description": "Enable or disable the execution of the policy"
        },
        "allowedValues": [
          "Deny",
          "Audit",
          "Disabled"
        ],
        "defaultValue": "Deny"
      }
  }

  PARAMETERS

  policy_rule = file("${path.module}/deny_non_cmk_storage_account_rule.json")
}

// Policy Assignment
resource "azurerm_policy_assignment" "audit_non_cmk_storage_ac" {
  name                 = "audit_non_cmk_storage_ac"
  scope                = var.assignment_scope
  policy_definition_id = azurerm_policy_definition.deny_non_cmk_storage_ac.id
  display_name         = "Audit storage account using MS Managed Key [BDD]"
  description          = "Audit storage account using MS Managed Key [BDD]"
  location             = var.location
  identity {
    type = "SystemAssigned"
  }

  parameters = <<PARAMETERS
  {
    "effect": {
      "value":"Audit"
    }
  }
  PARAMETERS

  not_scopes = var.audit_exclusion_list
}
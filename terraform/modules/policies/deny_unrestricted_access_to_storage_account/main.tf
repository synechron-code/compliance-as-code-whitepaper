// Policy Definition
resource "azurerm_policy_definition" "deny_unrestricted_access_to_storage_account" {
  name                = "deny_unrestricted_network_access_to_storage_account"
  policy_type         = "Custom"
  mode                = "All"
  display_name        = "Deny unrestricted network access to storage account"
  description         = "Deny unrestricted network access in your storage account firewall settings. Instead, configure network rules so only applications from allowed networks can access the storage account. To allow connections from specific internet or on-premise clients, access can be granted to traffic from specific Azure virtual networks or to public internet IP address ranges"
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

  parameters = file("${path.module}/policy_parameters.json")

  policy_rule = file("${path.module}/deny_unrestricted_access_to_storage_account.json")
}

// Policy Assignment
resource "azurerm_policy_assignment" "audit_storage_wo_net_acl" {
  name                 = "audit_storage_wo_net_acl"
  scope                = var.assignment_scope
  policy_definition_id = azurerm_policy_definition.deny_unrestricted_access_to_storage_account.id
  display_name         = "Audit unrestricted network access to storage account"
  description          = "Audit unrestricted network access to storage account"
  location             = var.location
  identity {
    type = "SystemAssigned"
  }

  parameters = jsonencode(
    {
      "effect" : {
        value : "Audit"
      },
      allowedAddressRanges : {
        "value" : var.whitelist_ips
      }
  })

  not_scopes = var.audit_exclusion_list
}

resource "azurerm_policy_assignment" "deny_storage_wo_net_acl" {
  name                 = "deny_storage_wo_net_acl"
  scope                = var.assignment_scope
  policy_definition_id = azurerm_policy_definition.deny_unrestricted_access_to_storage_account.id
  display_name         = "Deny unrestricted network access to storage account [BDD]"
  description          = "Deny unrestricted network access to storage account"
  location             = var.location
  identity {
    type = "SystemAssigned"
  }

  parameters = jsonencode(
    {
      "effect" : {
        "value" : "Deny"
      },
      allowedAddressRanges : {
        "value" : var.whitelist_ips
      }
  })

  not_scopes = var.deny_exclusion_list
}
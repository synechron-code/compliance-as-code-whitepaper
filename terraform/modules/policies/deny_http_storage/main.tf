locals {
  builtin_policy_id = "/providers/Microsoft.Authorization/policyDefinitions/404c3081-a854-4457-ae30-26a93ef643f9"
}

resource "azurerm_policy_assignment" "deny_http_storage" {
  name                 = "deny_http_storage"
  scope                = var.assignment_scope
  policy_definition_id = local.builtin_policy_id
  display_name         = "Secure transfer to storage accounts should be enabled"
  description          = "Secure transfer to storage accounts should be enabled"
  location             = var.location
  identity {
    type = "SystemAssigned"
  }

  parameters = <<PARAMETERS
  {
    "effect": {
      "value":"Deny"
    }
  }
  PARAMETERS

  not_scopes = var.deny_exclusion_list
}

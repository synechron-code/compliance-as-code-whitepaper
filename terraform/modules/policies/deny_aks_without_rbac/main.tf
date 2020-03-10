locals {
  name         = "deny_aks_without_rbac"
  display_name = "[AKS] Azure Kubernetes clusters must not be created without RBAC enabled [BDD]"
}

resource "azurerm_policy_definition" "deny_aks_without_rbac" {
  name                = local.name
  display_name        = local.display_name
  policy_type         = "Custom"
  mode                = "Indexed"
  policy_rule         = file("${path.module}/../../../resources/azure_policy/aks_rbac_deny.json")
  management_group_id = var.definition_management_group_id
  metadata            = <<METADATA
  {
    "category": "AKS"
  }
  METADATA

  lifecycle {
    ignore_changes = [
      metadata
    ]
  }
}

resource "azurerm_policy_assignment" "deny_aks_without_rbac" {
  name                 = local.name
  scope                = var.assignment_scope
  policy_definition_id = azurerm_policy_definition.deny_aks_without_rbac.id
  display_name         = local.display_name
  description          = "Azure Kubernetes clusters must not be created without RBAC enabled"
  location             = var.location
  identity {
    type = "SystemAssigned"
  }

  not_scopes = var.deny_exclusion_list
}


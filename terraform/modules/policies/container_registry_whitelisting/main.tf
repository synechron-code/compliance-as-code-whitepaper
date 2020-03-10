locals {
  builtin_policy_id = "/providers/Microsoft.Authorization/policyDefinitions/5f86cb6e-c4da-441b-807c-44bd0cc14e66"
}

resource "azurerm_policy_assignment" "container_registry_whitelisting" {
  name                 = "cr_whitelisting"
  scope                = var.assignment_scope
  policy_definition_id = local.builtin_policy_id
  display_name         = "[AKS] Ensure only allowed container images in AKS [BDD]"
  description          = "This policy ensures only allowed container images are running in an Azure Kubernetes Service cluster."
  location             = var.location
  identity {
    type = "SystemAssigned"
  }

  parameters = <<PARAMETERS
  {
    "effect": {
      "value": "EnforceRegoPolicy"
    },

    "allowedContainerImagesRegex": {
      "value": "${var.allowed_container_images_regex}"
    }
  }
  PARAMETERS

  not_scopes = var.deny_exclusion_list
}


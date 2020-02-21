variable "root_management_group_id" {
  type        = map(string)
  description = "ID of route level management group for different environment"
  default = {
    "dev" : "boxbank-root",
    "demo" : "ccasc-demo",
  }
}

output "root_management_group_id" {
  value = lookup(var.root_management_group_id, var.env)
}

variable "root_management_group_scope" {
  type        = map(string)
  description = "Scope of route level management group for different environment"
  default = {
    "dev" : "/providers/Microsoft.Management/managementgroups/boxbank-root",
    "demo" : "/providers/Microsoft.Management/managementgroups/ccasc-demo"
  }
}

output "root_management_group_scope" {
  value = lookup(var.root_management_group_scope, var.env)
}

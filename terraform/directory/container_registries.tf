// container_registry_whitelisting
variable "container_registry_whitelisting_allowed_regex" {
  type        = map(string)
  description = "regex for container_registry_whitelisting"
  default = {
    "dev"  = "^bdddevacr.azurecr.io/.+|^mcr.microsoft.com/.+",
    "demo" = "^bddacr.azurecr.io/.+|^mcr.microsoft.com/.+",
  }
}

output "container_registry_whitelisting_allowed_regex" {
  value = lookup(var.container_registry_whitelisting_allowed_regex, var.env)
}

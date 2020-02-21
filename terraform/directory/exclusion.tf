// deny_http_storage
variable "deny_http_storage_exclusion" {
  type        = map(list(string))
  description = "exclusion for deny_http_storage"
  default = {
    "dev"  = [],
    "demo" = [],
  }
}

output "deny_http_storage_exclusion" {
  value = var.deny_http_storage_exclusion[var.env]
}

variable "audit_http_storage_exclusion" {
  type        = map(list(string))
  description = "exclusion for audit_subnet_without_routetable"
  default = {
    "dev"  = [],
    "demo" = [],
  }
}

output "audit_http_storage_exclusion" {
  value = var.audit_http_storage_exclusion[var.env]
}

// deny_non_cmk_storage_account
variable "deny_non_cmk_storage_account_exclusion" {
  type        = map(list(string))
  description = "exclusion for deny_non_cmk_storage_account"
  default = {
    "dev"  = [],
    "demo" = [],
  }
}

output "deny_non_cmk_storage_account_exclusion" {
  value = var.deny_non_cmk_storage_account_exclusion[var.env]
}

variable "audit_non_cmk_storage_account_exclusion" {
  type        = map(list(string))
  description = "exclusion for audit_non_cmk_storage_account"
  default = {
    "dev"  = [],
    "demo" = [],
  }
}

output "audit_non_cmk_storage_account_exclusion" {
  value = var.audit_non_cmk_storage_account_exclusion[var.env]
}

// deny_unrestricted_access_to_storage_account
variable "deny_unrestricted_access_to_storage_account_exclusion" {
  type        = map(list(string))
  description = "exclusion for deny_unrestricted_access_to_storage_account"
  default = {
    "dev"  = [],
    "demo" = [],
  }
}

output "deny_unrestricted_access_to_storage_account_exclusion" {
  value = var.deny_unrestricted_access_to_storage_account_exclusion[var.env]
}

variable "audit_unrestricted_access_to_storage_account_exclusion" {
  type        = map(list(string))
  description = "exclusion for audit_unrestricted_access_to_storage_account"
  default = {
    "dev"  = [],
    "demo" = [],
  }
}

output "audit_unrestricted_access_to_storage_account_exclusion" {
  value = var.audit_unrestricted_access_to_storage_account_exclusion[var.env]
}

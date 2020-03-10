variable "definition_management_group_id" {
  description = "Policy Definition management group id."
  type        = string
}

variable "assignment_scope" {
  description = "Scope for assigning this policy"
  type        = string
}

variable "deny_exclusion_list" {
  description = "list of management group or subscription to be excluded for the deny assignment"
  type        = list(string)
}

variable "audit_exclusion_list" {
  description = "list of management group or subscription to be excluded for the audit assignment"
  type        = list(string)
}

variable "location" {
  description = "Azure location"
  type        = string
}

variable "whitelist_ips" {
  description = "Whitelisted IPs"
  type        = list(string)
}
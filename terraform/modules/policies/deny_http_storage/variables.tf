variable "assignment_scope" {
  description = "Scope for assigning this policy"
  type        = string
}

variable "deny_exclusion_list" {
  description = "list of management group or Subscription to be excluded for the deny assignment"
  type        = list(string)
}

variable "audit_exclusion_list" {
  description = "list of management group or Subscription to be excluded for the audit assignment"
  type        = list(string)
}

variable "location" {
  description = "Azure location"
  type        = string
}
// Use by deny_unrestricted_access_to_storage_account
variable "storage_account_whitelisted_ip" {
  type        = map(list(string))
  description = "IP ranges of whitelisted IPs for storage account"
  default = {
    "dev" : [
      "219.79.19.0/24",
      "170.74.231.168"
    ],
    "demo" : [
      "219.79.19.0/24",
      "170.74.231.168"
    ],
  }
}

output "storage_account_whitelisted_ip" {
  value = var.storage_account_whitelisted_ip[var.env]
}
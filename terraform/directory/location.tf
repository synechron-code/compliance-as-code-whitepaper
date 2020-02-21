variable "location" {
  type        = map(string)
  description = "azure location"
  default = {
    "dev" : "eastasia",
    "demo" : "uksouth",
  }
}

output "location" {
  value = lookup(var.location, var.env)
}

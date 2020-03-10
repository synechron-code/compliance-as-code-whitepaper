variable "location" {
  type = string
}

variable "name_prefix" {
  type = string
}

variable "vnet_rg_name" {
  type = string
}

variable "vnet_name" {
  type = string
}

variable "cluster_subnet_id" {
  type = string
}

variable "cluster_spn_id" {
  type = string
}

variable "cluster_spn_secret" {
  type = string
}

variable "node_public_ssh_key" {
  type    = string
  default = "generate"
}

variable "node_count" {
  type = string
}

variable "node_sku" {
  type = string
}

variable "tags" {
  type = map(string)
}

variable "rg_name" {
  type    = string
  default = "create_one_for_me"
}

variable "routes" {
  type = list(string)
  default = [
    "10.2.0.0/16",
  "10.3.0.0/16"]
}

variable "router_ip_address" {
  type    = string
  default = "10.1.0.36"
}

variable "aks_admin_username" {
  type    = string
  default = "citihub_admin"
}

variable "env" {
  default = "dev"
}


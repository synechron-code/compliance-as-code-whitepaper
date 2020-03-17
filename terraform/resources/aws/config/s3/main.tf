provider "aws" {}

locals {
  name_prefix = "bdd-demo"
}

module "encryption_as_rest_config_remediation" {
  source = "../../../../modules/aws/config/s3/encryption-at-rest-remediate"

  config_rule_name = var.encryption_at_rest_rule_name
  name_prefix = local.name_prefix
  remediation_action_name = var.encryption_at_rest_remediation_name
}

module "encryption_in_flight" {
  source = "../../../../modules/aws/config/s3/encryption-in-flight-remediate"

  config_rule_name = var.encryption_in_flight_rule_name
  name_prefix = local.name_prefix
  remediation_action_name = var.encryption_in_flight_remediation_name
}

module "ip_whitelisting" {
  source = "../../../../modules/aws/config/s3/ip-whitelist"

  config_rule_name = var.ip_whitelist_rule_name
  name_prefix = local.name_prefix
  ip_addresses = [ "219.79.19.94/24" ]
}
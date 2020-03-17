resource "aws_config_config_rule" "ip-whitelist-aws-config-rule" {
  name = lower(replace(var.config_rule_name, "_", "-"))

  source {
    owner             = "AWS"
    source_identifier = upper(replace(var.config_rule_name, "-", "_"))
  }

  scope {
    compliance_resource_types = [ "AWS::S3::Bucket" ]
  }

  input_parameters = <<PARAMS
[
  {
    "ParameterKey": "ipAddresses",
    "ParameterValue": "${join(", ", var.ip_addresses)}"
  }
]
PARAMS
}
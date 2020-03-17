
resource "aws_config_config_rule" "encryption-at-rest-aws-config-rule" {
  name = lower(replace(var.config_rule_name, "_", "-"))

  source {
    owner             = "AWS"
    source_identifier = upper(replace(var.config_rule_name, "-", "_"))
  }

  scope {
    compliance_resource_types = [ "AWS::S3::Bucket" ]
  }

}


resource "aws_cloudformation_stack" "encryption-at-rest-aws-config-remediation" {
  name = "${var.name_prefix}-${lower(replace(var.config_rule_name, "_", "-"))}"

  parameters = {
    ConfigRuleName = lower(replace(var.config_rule_name, "_", "-"))
    RemediationActionName = var.remediation_action_name
  }

  template_body = file("${path.module}/cloudformation.yaml")
  depends_on = [aws_config_config_rule.encryption-at-rest-aws-config-rule]
}
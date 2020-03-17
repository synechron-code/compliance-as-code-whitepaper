resource "aws_config_config_rule" "encryption-in-flight-aws-config-rule" {
  name = lower(replace(var.config_rule_name, "_", "-"))

  source {
    owner             = "AWS"
    source_identifier = upper(replace(var.config_rule_name, "-", "_"))
  }

  scope {
    compliance_resource_types = [ "AWS::S3::Bucket" ]
  }
}

resource "aws_cloudformation_stack" "encryption-in-flight-aws-config-remediation" {
  name = "${var.name_prefix}-${lower(replace(var.config_rule_name, "_", "-"))}"

  parameters = {
    ConfigRuleName = lower(replace(var.config_rule_name, "_", "-"))
    RemediationActionName = aws_ssm_document.encryption-in-flight-remediation-automation.name
  }

  template_body = file("${path.module}/cloudformation.yaml")

  depends_on = [aws_config_config_rule.encryption-in-flight-aws-config-rule]
}

resource "aws_ssm_document" "encryption-in-flight-remediation-automation" {
  name          = var.remediation_action_name
  document_type = "Automation"
  document_format = "YAML"

  content = file("${path.module}/ssm_document.yaml")
}
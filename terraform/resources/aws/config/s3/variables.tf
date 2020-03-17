variable "encryption_at_rest_rule_name" {
  type = string
  default = "S3_BUCKET_SERVER_SIDE_ENCRYPTION_ENABLED"
}

variable "encryption_at_rest_remediation_name" {
  type = string
  default = "AWS-EnableS3BucketEncryption"
}

variable "encryption_in_flight_rule_name" {
  type = string
  default = "S3_BUCKET_SSL_REQUESTS_ONLY"
}

variable "encryption_in_flight_remediation_name" {
  type = string
  default = "Citihub-set-s3-ssl-request-only"
}

variable "ip_whitelist_rule_name" {
  type = string
  default = "S3_BUCKET_POLICY_GRANTEE_CHECK"
}

provider "aws" {}

data "aws_region" "current" {}

resource "aws_s3_bucket" "bddipaddronlytest" {
  bucket = "bddipaddronlytest${data.aws_region.current.name}"
  acl    = "private"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
}

resource "aws_s3_bucket_policy" "bddipaddronlytest" {

  bucket = aws_s3_bucket.bddipaddronlytest.id
  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Id": "VPCe and SourceIP",
  "Statement": [
      {
          "Sid": "VPCe and SourceIP",
          "Effect": "Deny",
          "Principal": "*",
          "Action": "s3:*Object*",
          "Resource": [
              "arn:aws:s3:::bddipaddronlytest${data.aws_region.current.name}",
              "arn:aws:s3:::bddipaddronlytest${data.aws_region.current.name}/*"
          ],
          "Condition": {
              "StringNotLike": {
                  "aws:sourceVpce": [
                      "vpce-1111111",
                      "vpce-2222222"
                  ]
              },
              "NotIpAddress": {
                  "aws:SourceIp": [
                      "11.11.11.11/32",
                      "22.22.22.22/32"
                  ]
              }
          }
      },
      {
          "Effect": "Deny",
          "Principal": "*",
          "Action": "*",
          "Resource": "arn:aws:s3:::bddipaddronlytest${data.aws_region.current.name}/*",
          "Condition": {
              "Bool": {
                  "aws:SecureTransport": "false"
              }
          }
      }
  ]
}
POLICY
}
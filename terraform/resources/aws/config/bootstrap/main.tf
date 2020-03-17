provider "aws" {}

data "aws_region" "current" {}


resource "aws_iam_role" "config_iam_role" {
  name = "bdd-awsconfig-role"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "config.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
POLICY
}


resource "aws_iam_role_policy" "config_iam_policy" {
  name = "bdd-awsconfig-policy"
  role = aws_iam_role.config_iam_role.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
        "Action": "config:Put*",
        "Effect": "Allow",
        "Resource": "*"

    }
  ]
}
POLICY
}


resource "aws_config_configuration_recorder" "config_recorder" {
  name     = "BDDAWSConfigRecorder"
  role_arn = aws_iam_role.config_iam_role.arn
}

resource "aws_s3_bucket" "aws_config_recording_bucket" {
  bucket = "awsconfig-bucket-${data.aws_region.current.name}"
  acl = "private"
}

resource "aws_config_delivery_channel" "config_channel" {
  name           = "aws_config_delivery_channel"
  s3_bucket_name = aws_s3_bucket.aws_config_recording_bucket.bucket
}

resource "aws_config_configuration_recorder_status" "config_recorder_status" {
  name       = aws_config_configuration_recorder.config_recorder.name
  is_enabled = true
  depends_on = [ aws_config_delivery_channel.config_channel ]
}

resource "aws_iam_role_policy_attachment" "a" {
  role       = aws_iam_role.config_iam_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSConfigRole"
}

resource "aws_iam_role_policy" "p" {
  name = "awsconfig-example"
  role = aws_iam_role.config_iam_role.id

  policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:*"
      ],
      "Effect": "Allow",
      "Resource": [
        "${aws_s3_bucket.aws_config_recording_bucket.arn}",
        "${aws_s3_bucket.aws_config_recording_bucket.arn}/*"
      ]
    }
  ]
}
POLICY

}

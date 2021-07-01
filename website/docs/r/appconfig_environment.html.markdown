---
subcategory: "AppConfig"
layout: "aws"
page_title: "AWS: aws_appconfig_environment"
description: |-
  Provides an AppConfig Environment resource.
---

# Resource: aws_appconfig_environment

Provides an AppConfig Environment resource.

## Example Usage

### AppConfig Environment

```hcl
resource "aws_appconfig_application" "this" {
  name = "sample"
}
resource "aws_cloudwatch_metric_alarm" "test_alarm" {
  alarm_name                = "sample"
  comparison_operator       = "GreaterThanOrEqualToThreshold"
  evaluation_periods        = "2"
  metric_name               = "CPUUtilization"
  namespace                 = "AWS/EC2"
  period                    = "120"
  statistic                 = "Average"
  threshold                 = "80"
  alarm_description         = "This metric monitors ec2 cpu utilization"
  insufficient_data_actions = []
}
resource "aws_iam_role" "test_role" {
  name = "sample"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "ssm.amazonaws.com"
        }
      },
    ]
  })
}
resource "aws_appconfig_environment" "this" {
  application_id = aws_appconfig_application.this.id
  name           = "sample"
  description    = "just a sample"
  monitor {
    alarm_arn      = aws_cloudwatch_metric_alarm.test_alarm.arn
    alarm_role_arn = aws_iam_role.test_role.arn
  }
  tags = {
    Name = "AppConfig Environment"
  }
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required) The id of the parent AppConfig application for the environment.
* `name` - (Required) The name to use for the environment. Must be between 1 and 64 characters in length.
* `description` - (Optional) The description of the environment. Can be at most 1024 characters
* `monitor` - (Optional) The Cloudwatch monitors to connect to the environment.
* `tags` - (Optional) A map of tags to assign to the resource.

The `monitor` object supports the following:

* `alarm_arn` - (Required) ARN of the Amazon CloudWatch alarm.
* `alarm_role_arn` - (Required) ARN of the IAM role for AWS AppConfig to monitor AlarmArn.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the AppConfig Environment.
* `id` - The AppConfig Environment ID

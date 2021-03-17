---
subcategory: "AppConfig"
layout: "aws"
page_title: "AWS: aws_appconfig_deployment"
description: |-
  Provides an AppConfig Deployment resource.
---

# Resource: aws_appconfig_deployment

Provides an AppConfig Deployment resource.

## Example Usage

### AppConfig Deployment

```hcl
resource "aws_appconfig_application" "app" {
  name = "sample"
}
resource "aws_appconfig_environment" "env" {
  application_id = aws_appconfig_application.app.id
  name           = "sample"
}
resource "aws_appconfig_deployment_strategy" "strategy" {
  name                           = "sample"
  deployment_duration_in_minutes = 10
  growth_type                    = "LINEAR"
  replicate_to                   = "NONE"
}
resource "aws_appconfig_configuration_profile" "config" {
  application_id = aws_appconfig_application.app.id
  location_uri   = "hosted"
  name           = "sample"
}
resource "aws_appconfig_hosted_configuration_version" "config" {
  application_id           = aws_appconfig_application.app.id
  configuration_profile_id = aws_appconfig_configuration_profile.config.id
  content                  = "Settings"
  content_type             = "text/plain"
}
resource "aws_appconfig_deployment" "this" {
  application_id           = aws_appconfig_application.this.id
  environment_id           = aws_appconfig_environment.this.id
  deployment_strategy_id   = aws_appconfig_deployment_strategy.this.id
  configuration_profile_id = aws_appconfig_configuration_profile.this.id
  configuration_version    = aws_appconfig_hosted_configuration_version.config.version_number
  description              = "just a sample"
  tags = {
    Name = "AppConfig Deployment"
  }
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required) The id of deploy’s target application.
* `application_id` - (Required) The id of deploy’s target environment.
* `deployment_strategy_id`- (Required) The id of the deployment strategy to use.
* `configuration_profile_id`- (Required) The id of the configuration profile that contains the settings.
* `configuration_version` - (Required) The version of the settings to deploy.
* `description` - (Optional) The description of the hosted configuration version. Can be at most 1024 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the AppConfig Application.
* `id` - The AppConfig Application ID
* `deployment_number` - The deployment number

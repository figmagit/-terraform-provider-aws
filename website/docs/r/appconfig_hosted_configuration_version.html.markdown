---
subcategory: "AppConfig"
layout: "aws"
page_title: "AWS: aws_appconfig_hosted_configuration_version"
description: |-
  Provides an AppConfig Application resource.
---

# Resource: aws_appconfig_hosted_configuration_version

Provides an AppConfig Hosted Configuration Version resource.

## Example Usage

### AppConfig Hosted Configuration Version

```hcl
resource "aws_appconfig_application" "app" {
  name = "sample"
}
resource "aws_appconfig_configuration_profile" "config" {
  application_id = aws_appconfig_application.app.id
  location_uri   = "hosted"
  name           = "sample"
}
resource "aws_appconfig_hosted_configuration_version" "hosted" {
  application_id           = aws_appconfig_application.app.id
  configuration_profile_id = aws_appconfig_configuration_profile.config.id
  content                  = "Settings"
  content_type             = "text/plain"
  description              = "test"
}

```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required) The id of the parent AppConfig application for the hosted configuration version.
* `configuration_profile_id`-  (Required) The id of the parent AppConfig configuration profile for the hosted configuration version.
* `content` -  (Required) The content of the configuration or the configuration data.
* `content_type` -  (Required) A standard MIME type describing the format of the configuration content. AppConfig supports JSON ("application/json"), YAML ("application/x-yaml"), and plain text ("text/plain").
* `description` -  (Optional) The description of the hosted configuration version. Can be at most 1024 characters.
* `tags` - (Optional) A map of tags to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) specifying the role.
* `id` - Deployment ID.
* `version_number` - The version number of the hosted settings

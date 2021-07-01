---
subcategory: "AppConfig"
layout: "aws"
page_title: "AWS: aws_appconfig_configuration_profile"
description: |-
  Provides an AppConfig Application resource.
---

# Resource: aws_appconfig_configuration_profile

Provides an AppConfig Configuration Profile resource.

## Example Usage

### AppConfig Configuration Profile

```hcl
resource "aws_appconfig_application" "this" {
  name = "sample"
}

resource "aws_appconfig_configuration_profile" "this" {
  application_id     = aws_appconfig_application.this.id
  name               = "sample"
  description        = "just a sample"
  location_uri       = "hosted"
  retrieval_role_arn = aws_iam_role.retrieve_ssm.arn
  validator {
    type    = "JSON_SCHEMA"
    content = <<EOF
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "title": "$id$",
  "description": "BasicFeatureToggle-1",
  "type": "object",
  "additionalProperties": false,
  "patternProperties": {
    "[^\\s]+$": {
      "type": "boolean"
    }
  },
  "minProperties": 1
}
	EOF
  }
  tags = {
    Name = "AppConfig Configuration Profile"
  }
}
```

## Argument Reference

The following arguments are supported:

* `application_id` - (Required) The id of the parent AppConfig application for the configuration profile.
* `name` - (Required) The name to use for the configuration profile. Must be between 1 and 64 characters in length.
* `description` - (Optional) The description of the configuration profile. Can be at most 1024 characters.
* `location_uri` - (Required) A URI to locate the configuration. You can specify a Systems Manager (SSM) document or an SSM Parameter Store parameter. You also can use the AppConfig storage by entering “hosted”.
* `retrieval_role_arn` - (Optional) The description of the application. Can be at most 1024 characters.
* `validator` - (Optional) The validators to use when validating new configurations.
* `tags` - (Optional) A map of tags to assign to the resource.

The `validator` object supports the following:

* `type` - (Required) AppConfig supports validators of type JSON_SCHEMA and LAMBDA.
* `content` - (Required) Either the JSON Schema content or the Amazon Resource Name (ARN) of an AWS Lambda function that validates new configurations

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - Amazon Resource Name (ARN) of the Configuration Profile.
* `id` - Configuration Profile ID

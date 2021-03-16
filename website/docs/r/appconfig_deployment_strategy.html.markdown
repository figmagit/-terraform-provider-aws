---
subcategory: "AppConfig"
layout: "aws"
page_title: "AWS: aws_appconfig_deployment_strategy"
description: |-
  Provides an AppConfig Deployment Strategy resource.
---

# Resource: aws_appconfig_deployment_strategy

Provides an AppConfig Deployment Strategy resource.

## Example Usage

### AppConfig Deployment Strategy

```hcl
resource "aws_appconfig_deployment_strategy" "this" {
  name                           = "sample"
  description                    = "just a sample"
  deployment_duration_in_minutes = 15
  final_bake_time_in_minutes     = 30
  growth_factor                  = 0.1
  growth_type                    = "LINEAR"
  replicate_to                   = "NONE"
  tags = {
    Name = "AppConfig Deployment Strategy"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use for the deployment strategy. Must be between 1 and 64 characters in length.
* `description` - (Optional) The description of the deployment strategy. Can be at most 1024 characters.
* `deployment_duration_in_minutes` - (Required) The length of the deployment in minutes. Can be at most 1440 minutes.
* `final_bake_time_in_minutes` - (Optional) The length of the bake time at the end of the deployment. Can be at most 1440 minutes.
* `growth_factor` - (Optional) The percentage of targets that receive the new settings every interval. Defaults to 1. Must be between 1 and 100.
* `growth_type` - (Required) - The algorithm used to define how the percentage grows over time.
* `replicate_to` - (Required) - The Systems Manager (SSM) document where the deployment strategy is saved. Can be "NONE" if the deployment strategy is not saved.
* `tags` - (Optional) A map of tags to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The Amazon Resource Name (ARN) of the AppConfig Deployment Strategy.
* `id` - The AppConfig Application ID

## Import

Deployment Strategies can be imported using their ID, e.g.

```
$ terraform import aws_appconfig_deployment_strategy.bar l4fppbi
```

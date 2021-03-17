package aws

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsAppconfigDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsAppconfigDeploymentCreate,
		Read:   resourceAwsAppconfigDeploymentRead,
		Update: resourceAwsAppconfigDeploymentUpdate,
		Delete: resourceAwsAppconfigDeploymentDelete,

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"deployment_strategy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration_profile_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 1024),
				),
			},
			"tags": tagsSchema(),
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAwsAppconfigDeploymentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn
	applicationID := aws.String(d.Get("application_id").(string))
	environmentID := aws.String(d.Get("environment_id").(string))
	deploymentStrategyID := aws.String(d.Get("deployment_strategy_id").(string))
	configProfileID := aws.String(d.Get("configuration_profile_id").(string))
	deploymentDescription := aws.String(d.Get("description").(string))
	configVersion := aws.String(d.Get("configuration_version").(string))

	input := &appconfig.StartDeploymentInput{
		ApplicationId:          applicationID,
		EnvironmentId:          environmentID,
		DeploymentStrategyId:   deploymentStrategyID,
		ConfigurationProfileId: configProfileID,
		ConfigurationVersion:   configVersion,
		Description:            deploymentDescription,
		Tags:                   keyvaluetags.New(d.Get("tags").(map[string]interface{})).IgnoreAws().AppconfigTags(),
	}

	var deploy *appconfig.StartDeploymentOutput
	var err error
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		deploy, err = conn.StartDeployment(input)
		if err != nil {
			if isAWSErr(err, appconfig.ErrCodeConflictException, "") {
				log.Printf("[DEBUG] Retrying AppConfig Deployment creation: %s", err)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if isResourceTimeoutError(err) {
		deploy, err = conn.StartDeployment(input)
	}

	if err != nil {
		return fmt.Errorf("Creating AppConfig Deployment failed: %s", err)
	}

	d.SetId(fmt.Sprintf("%s-%s-%d", aws.StringValue(applicationID), aws.StringValue(environmentID), aws.Int64Value(deploy.DeploymentNumber)))
	d.Set("deployment_number", aws.Int64Value(deploy.DeploymentNumber))

	return resourceAwsAppconfigDeploymentRead(d, meta)
}

func resourceAwsAppconfigDeploymentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn
	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig

	applicationID := aws.String(d.Get("application_id").(string))
	environmentID := aws.String(d.Get("environment_id").(string))

	input := &appconfig.GetDeploymentInput{
		ApplicationId:    applicationID,
		EnvironmentId:    environmentID,
		DeploymentNumber: aws.Int64(int64(d.Get("deployment_number").(int))),
	}

	output, err := conn.GetDeployment(input)

	if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
		log.Printf("[WARN] Appconfig Deployment (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting AppConfig Deployment (%s): %s", d.Id(), err)
	}

	if output == nil {
		return fmt.Errorf("error getting AppConfig Deployment (%s): empty response", d.Id())
	}

	currentState := aws.StringValue(output.State)
	if currentState == appconfig.DeploymentStateRolledBack || currentState == appconfig.DeploymentStateRollingBack {
		log.Printf("[WARN] Appconfig Deployment (%s) rolled back, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	appID := aws.StringValue(output.ApplicationId)
	envID := aws.StringValue(output.EnvironmentId)
	deployNum := fmt.Sprintf("%d", aws.Int64Value(output.DeploymentNumber))
	arnResource := fmt.Sprintf("application/%s/environment/%s/deployment/%s", appID, envID, deployNum)
	appARN := arn.ARN{
		AccountID: meta.(*AWSClient).accountid,
		Partition: meta.(*AWSClient).partition,
		Region:    meta.(*AWSClient).region,
		Resource:  arnResource,
		Service:   "appconfig",
	}.String()

	d.Set("arn", appARN)
	d.Set("description", output.Description)
	d.Set("deployment_number", aws.Int64Value(output.DeploymentNumber))

	tags, err := keyvaluetags.AppconfigListTags(conn, appARN)
	if err != nil {
		return fmt.Errorf("error getting tags for AppConfig Deployment (%s): %s", d.Id(), err)
	}

	if err := d.Set("tags", tags.IgnoreAws().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

func resourceAwsAppconfigDeploymentUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := resourceAwsAppconfigDeploymentDelete(d, meta); err != nil {
		return fmt.Errorf("error rolling back existing deployment during update: %s", err)
	}

	if err := resourceAwsAppconfigDeploymentCreate(d, meta); err != nil {
		return fmt.Errorf("error starting new deployment during update: %s", err)
	}

	return resourceAwsAppconfigDeploymentRead(d, meta)
}

func resourceAwsAppconfigDeploymentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn

	input := &appconfig.StopDeploymentInput{
		ApplicationId:    aws.String(d.Get("application_id").(string)),
		EnvironmentId:    aws.String(d.Get("environment_id").(string)),
		DeploymentNumber: aws.Int64(int64(d.Get("deployment_number").(int))),
	}

	_, err := conn.StopDeployment(input)

	if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
		return nil
	}

	if isAWSErr(err, appconfig.ErrCodeBadRequestException, "it has a status of ROLLED_BACK") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error stopping Appconfig Deployment (%s): %s", d.Id(), err)
	}

	return nil
}

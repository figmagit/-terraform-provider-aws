package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAwsAppconfigHostedConfigurationVersion() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsAppconfigHostedConfigurationVersionCreate,
		Read:   resourceAwsAppconfigHostedConfigurationVersionRead,
		Update: resourceAwsAppconfigHostedConfigurationVersionUpdate,
		Delete: resourceAwsAppconfigHostedConfigurationVersionDelete,

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration_profile_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"application/json",
					"application/x-yaml",
					"text/plain",
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 2048),
				),
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAwsAppconfigHostedConfigurationVersionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn
	applicationID := aws.String(d.Get("application_id").(string))
	configProfileID := aws.String(d.Get("configuration_profile_id").(string))
	content := []byte(d.Get("content").(string))
	contentType := aws.String(d.Get("content_type").(string))
	configVersionDescription := aws.String(d.Get("description").(string))

	input := &appconfig.CreateHostedConfigurationVersionInput{
		ApplicationId:          aws.String(d.Get("application_id").(string)),
		ConfigurationProfileId: configProfileID,
		Content:                content,
		ContentType:            contentType,
		Description:            configVersionDescription,
	}

	hostedConfig, err := conn.CreateHostedConfigurationVersion(input)
	if err != nil {
		return fmt.Errorf("Creating AppConfig Hosted Configuration Version failed: %s", err)
	}
	log.Printf("[DEBUG] AppConfig Hosted Configuration Version created: %s", hostedConfig)

	d.SetId(fmt.Sprintf("%s-%s-%d", aws.StringValue(applicationID), aws.StringValue(configProfileID), aws.Int64Value(hostedConfig.VersionNumber)))
	d.Set("version_number", aws.Int64Value(hostedConfig.VersionNumber))

	return resourceAwsAppconfigHostedConfigurationVersionRead(d, meta)
}

func resourceAwsAppconfigHostedConfigurationVersionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn

	input := &appconfig.GetHostedConfigurationVersionInput{
		ApplicationId:          aws.String(d.Get("application_id").(string)),
		ConfigurationProfileId: aws.String(d.Get("configuration_profile_id").(string)),
		VersionNumber:          aws.Int64(int64(d.Get("version_number").(int))),
	}

	output, err := conn.GetHostedConfigurationVersion(input)

	if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
		log.Printf("[WARN] Appconfig Hosted Configuration Version (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting AppConfig Hosted Configuration Version (%s): %s", d.Id(), err)
	}

	if output == nil {
		return fmt.Errorf("error getting AppConfig Hosted Configuration Version (%s): empty response", d.Id())
	}

	appID := aws.StringValue(output.ApplicationId)
	profileID := aws.StringValue(output.ConfigurationProfileId)
	versionNum := fmt.Sprintf("%d", aws.Int64Value(output.VersionNumber))
	arnResource := fmt.Sprintf("application/%s/configurationprofile/%s/hostedconfigurationversion/%s", appID, profileID, versionNum)
	appARN := arn.ARN{
		AccountID: meta.(*AWSClient).accountid,
		Partition: meta.(*AWSClient).partition,
		Region:    meta.(*AWSClient).region,
		Resource:  arnResource,
		Service:   "appconfig",
	}.String()

	d.Set("arn", appARN)
	d.Set("description", output.Description)
	d.Set("version_number", aws.Int64Value(output.VersionNumber))
	d.Set("content", string(output.Content))
	d.Set("content_type", output.ContentType)

	return nil
}

func resourceAwsAppconfigHostedConfigurationVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := resourceAwsAppconfigHostedConfigurationVersionDelete(d, meta); err != nil {
		return fmt.Errorf("error deleting hosted configuration version during update: %s", err)
	}

	if err := resourceAwsAppconfigHostedConfigurationVersionCreate(d, meta); err != nil {
		return fmt.Errorf("error creating hosted configuration version during update: %s", err)
	}

	return resourceAwsAppconfigHostedConfigurationVersionRead(d, meta)
}

func resourceAwsAppconfigHostedConfigurationVersionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn

	input := &appconfig.DeleteHostedConfigurationVersionInput{
		ApplicationId:          aws.String(d.Get("application_id").(string)),
		ConfigurationProfileId: aws.String(d.Get("configuration_profile_id").(string)),
		VersionNumber:          aws.Int64(int64(d.Get("version_number").(int))),
	}

	_, err := conn.DeleteHostedConfigurationVersion(input)

	if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Appconfig Hosted Configuration Version (%s): %s", d.Id(), err)
	}

	return nil
}

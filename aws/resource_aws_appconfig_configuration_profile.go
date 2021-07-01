package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
	iamwaiter "github.com/terraform-providers/terraform-provider-aws/aws/internal/service/iam/waiter"
)

func resourceAwsAppconfigConfigurationProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsAppconfigConfigurationProfileCreate,
		Read:   resourceAwsAppconfigConfigurationProfileRead,
		Update: resourceAwsAppconfigConfigurationProfileUpdate,
		Delete: resourceAwsAppconfigConfigurationProfileDelete,

		Schema: map[string]*schema.Schema{
			"application_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
				),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 1024),
				),
			},
			"location_uri": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 2048),
				),
			},
			"retrieval_role_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateArn,
			},
			"validator": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{appconfig.ValidatorTypeJsonSchema, appconfig.ValidatorTypeLambda}, false),
						},
						"content": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"tags": tagsSchema(),
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAwsAppconfigConfigurationProfileCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn
	applicationID := aws.String(d.Get("application_id").(string))
	configProfileName := aws.String(d.Get("name").(string))
	configProfileDescription := aws.String(d.Get("description").(string))
	locationURI := aws.String(d.Get("location_uri").(string))
	retrievalRoleArn := aws.String(d.Get("retrieval_role_arn").(string))
	if *retrievalRoleArn == "" {
		retrievalRoleArn = nil
	}

	var validatorList []*appconfig.Validator
	if validators := d.Get("validator").(*schema.Set).List(); len(validators) > 0 {
		validatorList = convertMapToValidators(validators)
	}

	input := &appconfig.CreateConfigurationProfileInput{
		ApplicationId:    applicationID,
		Name:             configProfileName,
		Description:      configProfileDescription,
		LocationUri:      locationURI,
		RetrievalRoleArn: retrievalRoleArn,
		Validators:       validatorList,
		Tags:             keyvaluetags.New(d.Get("tags").(map[string]interface{})).IgnoreAws().AppconfigTags(),
	}

	expectedErrMsg := fmt.Sprintf("Error trying to assume role %s", d.Get("retrieval_role_arn").(string))
	expectedErrMsg2 := "is not authorized to perform the operation"
	var config *appconfig.CreateConfigurationProfileOutput
	var err error
	err = resource.Retry(iamwaiter.PropagationTimeout, func() *resource.RetryError {
		config, err = conn.CreateConfigurationProfile(input)
		if err != nil {
			if isAWSErr(err, "BadRequestException", expectedErrMsg) || isAWSErr(err, "BadRequestException", expectedErrMsg2) {
				log.Printf("[DEBUG] Retrying AppConfig Configuration Profile creation: %s", err)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if isResourceTimeoutError(err) {
		config, err = conn.CreateConfigurationProfile(input)
	}
	if err != nil {
		return fmt.Errorf("Creating AppConfig Configuration Profile failed: %s", err)
	}
	log.Printf("[DEBUG] AppConfig Configuration Profile created: %s", config)

	d.SetId(aws.StringValue(config.Id))

	return resourceAwsAppconfigConfigurationProfileRead(d, meta)
}

func resourceAwsAppconfigConfigurationProfileRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn
	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig

	input := &appconfig.GetConfigurationProfileInput{
		ApplicationId:          aws.String(d.Get("application_id").(string)),
		ConfigurationProfileId: aws.String(d.Id()),
	}

	output, err := conn.GetConfigurationProfile(input)

	if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
		log.Printf("[WARN] Appconfig Configuration Profile (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting AppConfig Configuration Profile (%s): %s", d.Id(), err)
	}

	if output == nil {
		return fmt.Errorf("error getting AppConfig Configuration Profile (%s): empty response", d.Id())
	}

	appARN := arn.ARN{
		AccountID: meta.(*AWSClient).accountid,
		Partition: meta.(*AWSClient).partition,
		Region:    meta.(*AWSClient).region,
		Resource:  fmt.Sprintf("application/%s/configurationprofile/%s", aws.StringValue(output.ApplicationId), aws.StringValue(output.Id)),
		Service:   "appconfig",
	}.String()

	d.Set("arn", appARN)
	d.Set("name", output.Name)
	d.Set("description", output.Description)
	d.Set("location_uri", output.LocationUri)
	d.Set("retrieval_role_arn", output.RetrievalRoleArn)
	d.Set("validator", convertValidatorsToMap(output.Validators))

	tags, err := keyvaluetags.AppconfigListTags(conn, appARN)
	if err != nil {
		return fmt.Errorf("error getting tags for AppConfig Configuration Profile (%s): %s", d.Id(), err)
	}

	if err := d.Set("tags", tags.IgnoreAws().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

func resourceAwsAppconfigConfigurationProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn

	if d.HasChange("tags") {
		o, n := d.GetChange("tags")
		if err := keyvaluetags.AppconfigUpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating AppConfig (%s) tags: %s", d.Id(), err)
		}
	}

	applicationID := aws.String(d.Get("application_id").(string))
	configProfileID := aws.String(d.Id())
	configProfileName := aws.String(d.Get("name").(string))
	configProfileDescription := aws.String(d.Get("description").(string))
	retrievalRoleArn := aws.String(d.Get("retrieval_role_arn").(string))
	if *retrievalRoleArn == "" {
		retrievalRoleArn = nil
	}
	var validatorList []*appconfig.Validator
	if validators := d.Get("validator").(*schema.Set).List(); len(validators) > 0 {
		validatorList = convertMapToValidators(validators)
	}

	updateInput := &appconfig.UpdateConfigurationProfileInput{
		ApplicationId:          applicationID,
		ConfigurationProfileId: configProfileID,
		Description:            configProfileDescription,
		Name:                   configProfileName,
		RetrievalRoleArn:       retrievalRoleArn,
		Validators:             validatorList,
	}

	if d.HasChange("description") {
		_, n := d.GetChange("description")
		updateInput.Description = aws.String(n.(string))
	}

	if d.HasChange("name") {
		_, n := d.GetChange("name")
		updateInput.Name = aws.String(n.(string))
	}

	expectedErrMsg := fmt.Sprintf("Error trying to assume role %s", d.Get("retrieval_role_arn").(string))
	expectedErrMsg2 := "is not authorized to perform the operation"
	var err error
	err = resource.Retry(iamwaiter.PropagationTimeout, func() *resource.RetryError {
		_, err = conn.UpdateConfigurationProfile(updateInput)
		if err != nil {
			if isAWSErr(err, "BadRequestException", expectedErrMsg) || isAWSErr(err, "BadRequestException", expectedErrMsg2) {
				log.Printf("[DEBUG] Retrying AppConfig Configuration Profile update: %s", err)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if isResourceTimeoutError(err) {
		_, err = conn.UpdateConfigurationProfile(updateInput)
	}
	if err != nil {
		return fmt.Errorf("Updating AppConfig Configuration Profile failed: %s", err)
	}

	return resourceAwsAppconfigConfigurationProfileRead(d, meta)
}

func resourceAwsAppconfigConfigurationProfileDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn

	input := &appconfig.DeleteConfigurationProfileInput{
		ApplicationId:          aws.String(d.Get("application_id").(string)),
		ConfigurationProfileId: aws.String(d.Id()),
	}

	_, err := conn.DeleteConfigurationProfile(input)

	if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Appconfig Configuration Profile (%s): %s", d.Id(), err)
	}

	return nil
}

func convertMapToValidators(validatorMap []interface{}) []*appconfig.Validator {
	validatorList := make([]*appconfig.Validator, len(validatorMap))
	for i, v := range validatorMap {
		vMap := v.(map[string]interface{})
		validator := appconfig.Validator{
			Type:    aws.String(vMap["type"].(string)),
			Content: aws.String(vMap["content"].(string)),
		}
		validatorList[i] = &validator
	}
	return validatorList
}

func convertValidatorsToMap(validatorList []*appconfig.Validator) []map[string]string {
	validatorMap := make([]map[string]string, len(validatorList))
	for i, v := range validatorList {
		validator := map[string]string{
			"type":    aws.StringValue(v.Type),
			"content": aws.StringValue(v.Content),
		}
		validatorMap[i] = validator
	}
	return validatorMap
}

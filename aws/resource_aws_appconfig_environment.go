package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func resourceAwsAppconfigEnvironment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsAppconfigEnvironmentCreate,
		Read:   resourceAwsAppconfigEnvironmentRead,
		Update: resourceAwsAppconfigEnvironmentUpdate,
		Delete: resourceAwsAppconfigEnvironmentDelete,

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
			"monitor": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_arn": {
							Type:     schema.TypeString,
							Required: true,
						},
						"alarm_role_arn": {
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

func resourceAwsAppconfigEnvironmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn
	applicationID := aws.String(d.Get("application_id").(string))
	environmentName := aws.String(d.Get("name").(string))
	environmentDescription := aws.String(d.Get("description").(string))

	var monitorList []*appconfig.Monitor
	if monitors := d.Get("monitor").(*schema.Set).List(); len(monitors) > 0 {
		monitorList = convertMapToMonitors(monitors)
	}

	input := &appconfig.CreateEnvironmentInput{
		ApplicationId: applicationID,
		Name:          environmentName,
		Description:   environmentDescription,
		Monitors:      monitorList,
		Tags:          keyvaluetags.New(d.Get("tags").(map[string]interface{})).IgnoreAws().AppconfigTags(),
	}

	env, err := conn.CreateEnvironment(input)
	if err != nil {
		return fmt.Errorf("Creating AppConfig Environment failed: %s", err)
	}
	log.Printf("[DEBUG] AppConfig Environment created: %s", env)

	d.SetId(aws.StringValue(env.Id))

	return resourceAwsAppconfigEnvironmentRead(d, meta)
}

func resourceAwsAppconfigEnvironmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn
	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig

	input := &appconfig.GetEnvironmentInput{
		ApplicationId: aws.String(d.Get("application_id").(string)),
		EnvironmentId: aws.String(d.Id()),
	}

	output, err := conn.GetEnvironment(input)

	if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
		log.Printf("[WARN] Appconfig Environment (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting AppConfig Environment (%s): %s", d.Id(), err)
	}

	if output == nil {
		return fmt.Errorf("error getting AppConfig Environment (%s): empty response", d.Id())
	}

	appARN := arn.ARN{
		AccountID: meta.(*AWSClient).accountid,
		Partition: meta.(*AWSClient).partition,
		Region:    meta.(*AWSClient).region,
		Resource:  fmt.Sprintf("application/%s/environment/%s", aws.StringValue(output.ApplicationId), aws.StringValue(output.Id)),
		Service:   "appconfig",
	}.String()

	d.Set("arn", appARN)
	d.Set("name", output.Name)
	d.Set("description", output.Description)
	d.Set("monitor", convertMonitorsToMap(output.Monitors))

	tags, err := keyvaluetags.AppconfigListTags(conn, appARN)
	if err != nil {
		return fmt.Errorf("error getting tags for AppConfig Environment (%s): %s", d.Id(), err)
	}

	if err := d.Set("tags", tags.IgnoreAws().IgnoreConfig(ignoreTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}

func resourceAwsAppconfigEnvironmentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn

	if d.HasChange("tags") {
		o, n := d.GetChange("tags")
		if err := keyvaluetags.AppconfigUpdateTags(conn, d.Get("arn").(string), o, n); err != nil {
			return fmt.Errorf("error updating AppConfig (%s) tags: %s", d.Id(), err)
		}
	}

	if d.HasChanges("name", "description", "monitors") {
		applicationID := aws.String(d.Get("application_id").(string))
		environmentID := aws.String(d.Id())
		environmentName := aws.String(d.Get("name").(string))
		environmentDescription := aws.String(d.Get("description").(string))

		var monitorList []*appconfig.Monitor
		if monitors := d.Get("monitor").(*schema.Set).List(); len(monitors) > 0 {
			monitorList = convertMapToMonitors(monitors)
		}

		updateInput := &appconfig.UpdateEnvironmentInput{
			ApplicationId: applicationID,
			EnvironmentId: environmentID,
			Description:   environmentDescription,
			Name:          environmentName,
			Monitors:      monitorList,
		}

		_, err := conn.UpdateEnvironment(updateInput)
		if err != nil {
			return fmt.Errorf("Updating AppConfig Environment failed: %s", err)
		}
	}

	return resourceAwsAppconfigEnvironmentRead(d, meta)
}

func resourceAwsAppconfigEnvironmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).appconfigconn

	input := &appconfig.DeleteEnvironmentInput{
		ApplicationId: aws.String(d.Get("application_id").(string)),
		EnvironmentId: aws.String(d.Id()),
	}

	_, err := conn.DeleteEnvironment(input)

	if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting Appconfig Environment (%s): %s", d.Id(), err)
	}

	return nil
}

func convertMapToMonitors(monitorMap []interface{}) []*appconfig.Monitor {
	monitorList := make([]*appconfig.Monitor, len(monitorMap))
	for i, v := range monitorMap {
		vMap := v.(map[string]interface{})
		monitor := appconfig.Monitor{
			AlarmArn:     aws.String(vMap["alarm_arn"].(string)),
			AlarmRoleArn: aws.String(vMap["alarm_role_arn"].(string)),
		}
		monitorList[i] = &monitor
	}
	return monitorList
}

func convertMonitorsToMap(monitorList []*appconfig.Monitor) []map[string]string {
	monitorMap := make([]map[string]string, len(monitorList))
	for i, v := range monitorList {
		monitor := map[string]string{
			"alarm_arn":      aws.StringValue(v.AlarmArn),
			"alarm_role_arn": aws.StringValue(v.AlarmRoleArn),
		}
		monitorMap[i] = monitor
	}
	return monitorMap
}

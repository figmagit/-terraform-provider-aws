package aws

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAWSAppConfigDeployment_basic(t *testing.T) {
	var deployment appconfig.GetDeploymentOutput

	baseName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("desc")
	resourceName := "aws_appconfig_deployment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigDeploymentName(baseName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentExists(resourceName, &deployment),
					testAccCheckAWSAppConfigDeploymentARN(resourceName, &deployment),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "configuration_version", "1"),
				),
			},
		},
	})
}

func TestAccAWSAppConfigDeployment_disappears(t *testing.T) {
	var deployment appconfig.GetDeploymentOutput

	baseName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_deployment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigDeploymentName(baseName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentExists(resourceName, &deployment),
					testAccCheckAWSAppConfigDeploymentDisappears(&deployment),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSAppConfigDeployment_Tags(t *testing.T) {
	var deployment appconfig.GetDeploymentOutput

	baseName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_deployment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigDeploymentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigDeploymentTags1(baseName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentExists(resourceName, &deployment),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccAWSAppConfigDeploymentTags2(baseName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentExists(resourceName, &deployment),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckAppConfigDeploymentDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).appconfigconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_appconfig_deployment" {
			continue
		}

		deploymentNumber, err := strconv.ParseInt(rs.Primary.Attributes["deployment_number"], 10, 64)
		if err != nil {
			return err
		}

		input := &appconfig.GetDeploymentInput{
			ApplicationId:    aws.String(rs.Primary.Attributes["application_id"]),
			EnvironmentId:    aws.String(rs.Primary.Attributes["environment_id"]),
			DeploymentNumber: aws.Int64(deploymentNumber),
		}

		output, err := conn.GetDeployment(input)

		if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
			continue
		}

		if err != nil {
			return err
		}

		currentState := aws.StringValue(output.State)
		if currentState == appconfig.DeploymentStateRolledBack ||
			currentState == appconfig.DeploymentStateRollingBack ||
			currentState == appconfig.DeploymentStateComplete {
			return nil
		}

		if output != nil {
			return fmt.Errorf("AppConfig Deployment (%s) still exists", rs.Primary.ID)
		}
	}

	return nil

}

func testAccCheckAWSAppConfigDeploymentDisappears(deployment *appconfig.GetDeploymentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).appconfigconn

		input := &appconfig.StopDeploymentInput{
			ApplicationId:    deployment.ApplicationId,
			EnvironmentId:    deployment.EnvironmentId,
			DeploymentNumber: deployment.DeploymentNumber,
		}

		_, err := conn.StopDeployment(input)

		return err
	}
}

func testAccCheckAWSAppConfigDeploymentExists(resourceName string, deployment *appconfig.GetDeploymentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		conn := testAccProvider.Meta().(*AWSClient).appconfigconn

		deploymentNumber, err := strconv.ParseInt(rs.Primary.Attributes["deployment_number"], 10, 64)
		if err != nil {
			return err
		}

		input := &appconfig.GetDeploymentInput{
			ApplicationId:    aws.String(rs.Primary.Attributes["application_id"]),
			EnvironmentId:    aws.String(rs.Primary.Attributes["environment_id"]),
			DeploymentNumber: aws.Int64(deploymentNumber),
		}

		output, err := conn.GetDeployment(input)
		if err != nil {
			return err
		}

		*deployment = *output

		return nil
	}
}

func testAccCheckAWSAppConfigDeploymentARN(resourceName string, deployment *appconfig.GetDeploymentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		appID := aws.StringValue(deployment.ApplicationId)
		envID := aws.StringValue(deployment.EnvironmentId)
		deployNum := fmt.Sprintf("%d", aws.Int64Value(deployment.DeploymentNumber))
		arnResource := fmt.Sprintf("application/%s/environment/%s/deployment/%s", appID, envID, deployNum)
		return testAccCheckResourceAttrRegionalARN(resourceName, "arn", "appconfig", arnResource)(s)
	}
}

func testAccAWSAppConfigDeploymentSetup(baseName string) string {
	appName := fmt.Sprintf("%s-app", baseName)
	envName := fmt.Sprintf("%s-env", baseName)
	strategyName := fmt.Sprintf("%s-strategy", baseName)
	configName := fmt.Sprintf("%s-config", baseName)
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_appconfig_environment" "env" {
  application_id = aws_appconfig_application.app.id
  name           = %[2]q
}
resource "aws_appconfig_deployment_strategy" "strategy" {
  name                           = %[3]q
  deployment_duration_in_minutes = 10
  growth_type                    = "LINEAR"
  replicate_to                   = "NONE"
}
resource "aws_appconfig_configuration_profile" "config" {
  application_id = aws_appconfig_application.app.id
  location_uri   = "hosted"
  name           = %[4]q
}
resource "aws_appconfig_hosted_configuration_version" "hosted" {
  application_id           = aws_appconfig_application.app.id
  configuration_profile_id = aws_appconfig_configuration_profile.config.id
  content                  = "Settings"
  content_type             = "text/plain"
}
`, appName, envName, strategyName, configName)
}

func testAccAWSAppConfigDeploymentName(baseName, rDesc string) string {
	return fmt.Sprintf(`
%[1]s
resource "aws_appconfig_deployment" "test" {
  application_id           = aws_appconfig_application.app.id
  environment_id           = aws_appconfig_environment.env.id
  deployment_strategy_id   = aws_appconfig_deployment_strategy.strategy.id
  configuration_profile_id = aws_appconfig_configuration_profile.config.id
  configuration_version    = aws_appconfig_hosted_configuration_version.hosted.version_number
  description              = %[2]q
}
`, testAccAWSAppConfigDeploymentSetup(baseName), rDesc)
}

func testAccAWSAppConfigDeploymentTags1(baseName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
%[1]s
resource "aws_appconfig_deployment" "test" {
  application_id           = aws_appconfig_application.app.id
  environment_id           = aws_appconfig_environment.env.id
  deployment_strategy_id   = aws_appconfig_deployment_strategy.strategy.id
  configuration_profile_id = aws_appconfig_configuration_profile.config.id
  configuration_version    = aws_appconfig_hosted_configuration_version.hosted.version_number
  tags = {
    %[2]q = %[3]q
  }
}
`, testAccAWSAppConfigDeploymentSetup(baseName), tagKey1, tagValue1)
}

func testAccAWSAppConfigDeploymentTags2(baseName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
%[1]s
resource "aws_appconfig_deployment" "test" {
  application_id           = aws_appconfig_application.app.id
  environment_id           = aws_appconfig_environment.env.id
  deployment_strategy_id   = aws_appconfig_deployment_strategy.strategy.id
  configuration_profile_id = aws_appconfig_configuration_profile.config.id
  configuration_version    = aws_appconfig_hosted_configuration_version.hosted.version_number
  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, testAccAWSAppConfigDeploymentSetup(baseName), tagKey1, tagValue1, tagKey2, tagValue2)
}

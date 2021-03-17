package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAWSAppConfigEnvironment_basic(t *testing.T) {
	var environment appconfig.GetEnvironmentOutput
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("desc")
	resourceName := "aws_appconfig_environment.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigEnvironmentName(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigEnvironmentARN(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "monitor.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
				),
			},
		},
	})
}

func TestAccAWSAppConfigEnvironment_disappears(t *testing.T) {
	var environment appconfig.GetEnvironmentOutput

	rName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_environment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigEnvironmentName(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigEnvironmentExists(resourceName, &environment),
					testAccCheckAWSAppConfigEnvironmentDisappears(&environment),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSAppConfigEnvironment_Monitors(t *testing.T) {
	var environment appconfig.GetEnvironmentOutput
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_environment.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigEnvironmentMonitors(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigEnvironmentARN(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "monitor.#", "1"),
				),
			},
		},
	})
}

func TestAccAWSAppConfigEnvironment_Tags(t *testing.T) {
	var environment appconfig.GetEnvironmentOutput

	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_environment.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigEnvironmentTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccAWSAppConfigEnvironmentTags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigEnvironmentExists(resourceName, &environment),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckAppConfigEnvironmentDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).appconfigconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_appconfig_environment" {
			continue
		}

		input := &appconfig.GetEnvironmentInput{
			ApplicationId: aws.String(rs.Primary.Attributes["application_id"]),
			EnvironmentId: aws.String(rs.Primary.ID),
		}

		output, err := conn.GetEnvironment(input)

		if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
			continue
		}

		if err != nil {
			return err
		}

		if output != nil {
			return fmt.Errorf("AppConfig Environment (%s) still exists", rs.Primary.ID)
		}
	}

	return nil

}

func testAccCheckAWSAppConfigEnvironmentDisappears(environment *appconfig.GetEnvironmentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).appconfigconn

		input := &appconfig.DeleteEnvironmentInput{
			ApplicationId: aws.String(*environment.ApplicationId),
			EnvironmentId: aws.String(*environment.Id),
		}

		_, err := conn.DeleteEnvironment(input)

		return err
	}
}

func testAccCheckAWSAppConfigEnvironmentExists(resourceName string, environment *appconfig.GetEnvironmentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		conn := testAccProvider.Meta().(*AWSClient).appconfigconn

		input := &appconfig.GetEnvironmentInput{
			ApplicationId: aws.String(rs.Primary.Attributes["application_id"]),
			EnvironmentId: aws.String(rs.Primary.ID),
		}

		output, err := conn.GetEnvironment(input)
		if err != nil {
			return err
		}

		*environment = *output

		return nil
	}
}

func testAccCheckAWSAppConfigEnvironmentARN(resourceName string, environment *appconfig.GetEnvironmentOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		arnResource := fmt.Sprintf("application/%s/environment/%s", aws.StringValue(environment.ApplicationId), aws.StringValue(environment.Id))
		return testAccCheckResourceAttrRegionalARN(resourceName, "arn", "appconfig", arnResource)(s)
	}
}

func testAccAWSAppConfigEnvironmentName(rName, rDesc string) string {
	appName := fmt.Sprintf("%s-app", rName)
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_appconfig_environment" "test" {
  name           = %[2]q
  application_id = aws_appconfig_application.app.id
  description    = %[3]q
}
`, appName, rName, rDesc)
}

func testAccAWSAppConfigEnvironmentMonitors(rName string) string {
	alarmName := acctest.RandomWithPrefix("test-alarm")
	roleName := acctest.RandomWithPrefix("test-role")
	appName := fmt.Sprintf("%s-app", rName)
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_appconfig_environment" "test" {
  name           = %[2]q
  application_id = aws_appconfig_application.app.id
  monitor {
    alarm_arn      = aws_cloudwatch_metric_alarm.test_alarm.arn
    alarm_role_arn = aws_iam_role.test_role.arn
  }
}
resource "aws_cloudwatch_metric_alarm" "test_alarm" {
  alarm_name                = %[3]q
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
  name = %[4]q

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
`, appName, rName, alarmName, roleName)
}

func testAccAWSAppConfigEnvironmentTags1(rName, tagKey1, tagValue1 string) string {
	appName := fmt.Sprintf("%s-app", rName)
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_appconfig_environment" "test" {
  name           = %[2]q
  application_id = aws_appconfig_application.app.id
  tags = {
    %[3]q = %[4]q
  }
}
`, appName, rName, tagKey1, tagValue1)
}

func testAccAWSAppConfigEnvironmentTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	appName := fmt.Sprintf("%s-app", rName)
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_appconfig_environment" "test" {
  name           = %[2]q
  application_id = aws_appconfig_application.app.id
  tags = {
    %[3]q = %[4]q
    %[5]q = %[6]q
  }
}
`, appName, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

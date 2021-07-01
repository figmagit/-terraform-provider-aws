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

func TestAccAWSAppConfigConfigurationProfile_basic(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput

	appName := acctest.RandomWithPrefix("tf-acc-test")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileName(appName, rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigConfigurationProfileARN(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "location_uri", "hosted"),
					resource.TestCheckResourceAttr(resourceName, "retrieval_role_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "validator.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "application_id"),
				),
			},
		},
	})
}

func TestAccAWSAppConfigConfigurationProfile_disappears(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput

	appName := acctest.RandomWithPrefix("tf-acc-test")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileName(appName, rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigConfigurationProfileExists(resourceName, &profile),
					testAccCheckAWSAppConfigConfigurationProfileDisappears(&profile),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
func TestAccAWSAppConfigConfigurationProfile_LocationURI_SSMParameter(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput
	appName := acctest.RandomWithPrefix("tf-acc-test")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileLocationSSMParameter(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigConfigurationProfileARN(resourceName, &profile),
					resource.TestCheckResourceAttrSet(resourceName, "location_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "retrieval_role_arn"),
				),
			},
		},
	})
}
func TestAccAWSAppConfigConfigurationProfile_LocationURI_SSMDocument(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput
	appName := acctest.RandomWithPrefix("tf-acc-test")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileLocationSSMDocument(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigConfigurationProfileARN(resourceName, &profile),
					resource.TestCheckResourceAttrSet(resourceName, "location_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "retrieval_role_arn"),
				),
			},
		},
	})
}

func TestAccAWSAppConfigConfigurationProfile_Validators(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput
	appName := acctest.RandomWithPrefix("tf-acc-test")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileValidator(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigConfigurationProfileARN(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "validator.#", "2"),
				),
			},
		},
	})
}
func TestAccAWSAppConfigConfigurationProfile_Tags(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput

	appName := acctest.RandomWithPrefix("tf-acc-test")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileTags1(appName, rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccAWSAppConfigConfigurationProfileTags2(appName, rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAWSAppConfigConfigurationProfileTags1(appName, rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckAWSAppConfigConfigurationProfileARN(resourceName string, config *appconfig.GetConfigurationProfileOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceArn := fmt.Sprintf("application/%s/configurationprofile/%s", aws.StringValue(config.ApplicationId), aws.StringValue(config.Id))
		return testAccCheckResourceAttrRegionalARN(resourceName, "arn", "appconfig", resourceArn)(s)
	}
}

func testAccCheckAppConfigConfigurationProfileDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).appconfigconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_appconfig_configuration_profile" {
			continue
		}

		input := &appconfig.GetConfigurationProfileInput{
			ApplicationId:          aws.String(rs.Primary.Attributes["application_id"]),
			ConfigurationProfileId: aws.String(rs.Primary.ID),
		}

		output, err := conn.GetConfigurationProfile(input)

		if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
			continue
		}

		if err != nil {
			return err
		}

		if output != nil {
			return fmt.Errorf("AppConfig Configuration Profile (%s) still exists", rs.Primary.ID)
		}
	}

	return nil

}

func testAccCheckAWSAppConfigConfigurationProfileDisappears(profile *appconfig.GetConfigurationProfileOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).appconfigconn

		input := &appconfig.DeleteConfigurationProfileInput{
			ApplicationId:          aws.String(*profile.ApplicationId),
			ConfigurationProfileId: aws.String(*profile.Id),
		}

		_, err := conn.DeleteConfigurationProfile(input)

		return err
	}
}

func testAccCheckAWSAppConfigConfigurationProfileExists(resourceName string, profile *appconfig.GetConfigurationProfileOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		conn := testAccProvider.Meta().(*AWSClient).appconfigconn

		input := &appconfig.GetConfigurationProfileInput{
			ApplicationId:          aws.String(rs.Primary.Attributes["application_id"]),
			ConfigurationProfileId: aws.String(rs.Primary.ID),
		}

		output, err := conn.GetConfigurationProfile(input)

		if err != nil {
			return err
		}

		*profile = *output

		return nil
	}
}

func testAccAWSAppConfigConfigurationProfileName(appName, rName, rDesc string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_appconfig_configuration_profile" "test" {
  name           = %[2]q
  description    = %[3]q
  application_id = aws_appconfig_application.app.id
  location_uri   = "hosted"
}
`, appName, rName, rDesc)
}

func testAccAWSAppConfigConfigurationProfileIAMRole() string {
	roleName := acctest.RandomWithPrefix("test-role")
	attatchmentName := acctest.RandomWithPrefix("test-attatchment")
	policyName := acctest.RandomWithPrefix("test-policy")

	return fmt.Sprintf(`
resource "aws_iam_role" "test_role" {
  name = %[1]q

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
resource "aws_iam_policy_attachment" "test_attach" {
  name       = %[2]q
  roles      = [aws_iam_role.test_role.name]
  policy_arn = aws_iam_policy.test_policy.arn
}
resource "aws_iam_policy" "test_policy" {
  name = %[3]q

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "ssm:GetParameter*",
          "ssm:DescribeParameters",
          "ssm:PutParameter",
          "ssm:GetDocument"
        ]
        Effect   = "Allow"
        Resource = "*"
      },
    ]
  })
}
`, roleName, attatchmentName, policyName)
}

func testAccAWSAppConfigConfigurationProfileLocationSSMParameter(appName, rName string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_ssm_parameter" "ssm_param" {
  name  = "foo"
  type  = "String"
  value = "bar"
}
resource "aws_appconfig_configuration_profile" "test" {
  name               = %[2]q
  application_id     = aws_appconfig_application.app.id
  location_uri       = aws_ssm_parameter.ssm_param.arn
  retrieval_role_arn = aws_iam_role.test_role.arn
}
%[3]s
`, appName, rName, testAccAWSAppConfigConfigurationProfileIAMRole())
}

func testAccAWSAppConfigConfigurationProfileLocationSSMDocument(appName, rName string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_ssm_document" "ssm_doc" {
  name          = "test_document"
  document_type = "Command"

  content = <<DOC
	{
	  "schemaVersion": "1.2",
	  "description": "Check ip configuration of a Linux instance.",
	  "parameters": {
  
	  },
	  "runtimeConfig": {
		"aws:runShellScript": {
		  "properties": [
			{
			  "id": "0.aws:runShellScript",
			  "runCommand": ["ifconfig"]
			}
		  ]
		}
	  }
	}
  DOC
}
resource "aws_appconfig_configuration_profile" "test" {
  name               = %[2]q
  application_id     = aws_appconfig_application.app.id
  location_uri       = aws_ssm_document.ssm_doc.arn
  retrieval_role_arn = aws_iam_role.test_role.arn
}
%[3]s
`, appName, rName, testAccAWSAppConfigConfigurationProfileIAMRole())
}

func testAccAWSAppConfigConfigurationProfileValidator(appName, rName string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_lambda_function" "test_lambda" {
  function_name = "lambda_function_name"
  role          = aws_iam_role.test_role_lambda.arn
  filename      = "test-fixtures/lambdatest.zip"
  handler       = "exports.test"
  runtime       = "nodejs12.x"
}
resource "aws_iam_role" "test_role_lambda" {
  name = "test_role_lambda"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Sid    = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      },
    ]
  })
}
resource "aws_appconfig_configuration_profile" "test" {
  name           = %[2]q
  application_id = aws_appconfig_application.app.id
  location_uri   = "hosted"
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
  validator {
    type    = "LAMBDA"
    content = aws_lambda_function.test_lambda.arn
  }
}
`, appName, rName)
}

func testAccAWSAppConfigConfigurationProfileTags1(appName, rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_appconfig_configuration_profile" "test" {
  name           = %[2]q
  application_id = aws_appconfig_application.app.id
  location_uri   = "hosted"
  tags = {
    %[3]q = %[4]q
  }
}
`, appName, rName, tagKey1, tagValue1)
}

func testAccAWSAppConfigConfigurationProfileTags2(appName, rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_appconfig_configuration_profile" "test" {
  name           = %[2]q
  application_id = aws_appconfig_application.app.id
  location_uri   = "hosted"
  tags = {
    %[3]q = %[4]q
    %[5]q = %[6]q
  }
}
`, appName, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

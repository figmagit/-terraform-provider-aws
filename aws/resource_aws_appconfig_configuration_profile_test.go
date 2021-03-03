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

func TestAccAwsAppConfigConfigurationProfile_basic(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput

	appName := acctest.RandomWithPrefix("tf-acc-test")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileName(appName, rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigConfigurationProfileARN(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "location_uri", "hosted"),
					resource.TestCheckResourceAttr(resourceName, "retrieval_role_arn", ""),
					resource.TestCheckResourceAttr(resourceName, "validator.%", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "application_id"),
				),
			},
		},
	})
}

func TestAccAwsAppConfigConfigurationProfile_disappears(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput

	appName := acctest.RandomWithPrefix("tf-acc-test")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileName(appName, rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					testAccCheckAwsAppConfigConfigurationProfileDisappears(&profile),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
func TestAccAWSAppConfigConfigurationProfile_LocationURIs(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput
	appName := acctest.RandomWithPrefix("tf-acc-test")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileLocationSSMParameter(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigConfigurationProfileARN(resourceName, &profile),
					resource.TestCheckResourceAttrSet(resourceName, "location_uri"),
					resource.TestCheckResourceAttrSet(resourceName, "retrieval_role_arn"),
				),
			},
			{
				Config: testAccAWSAppConfigConfigurationProfileLocationSSMDocument(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
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
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileValidator(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigConfigurationProfileARN(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccAWSAppConfigConfigurationProfile_RetrievalARN(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput
	appName := acctest.RandomWithPrefix("tf-acc-test")
	rName := acctest.RandomWithPrefix("tf-acc-test")

	resourceName := "aws_appconfig_configuration_profile.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileRetreivalARN(appName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigConfigurationProfileARN(resourceName, &profile),
					resource.TestCheckResourceAttrSet(resourceName, "retrieval_role_arn"),
				),
			},
		},
	})
}

func TestAccAwsAppConfigConfigurationProfile_Tags(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput

	appName := acctest.RandomWithPrefix("tf-acc-test")
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigConfigurationProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileTags1(appName, rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				Config: testAccAWSAppConfigConfigurationProfileTags2(appName, rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAWSAppConfigConfigurationProfileTags1(appName, rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
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

func testAccCheckAwsAppConfigConfigurationProfileDisappears(profile *appconfig.GetConfigurationProfileOutput) resource.TestCheckFunc {
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

func testAccCheckAwsAppConfigConfigurationProfileExists(resourceName string, profile *appconfig.GetConfigurationProfileOutput) resource.TestCheckFunc {
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
  name = %[2]q
  description = %[3]q
  application_id = aws_appconfig_application.app.id
  location_uri = "hosted"
}
`, appName, rName, rDesc)
}

func testAccAWSAppConfigConfigurationProfileLocationSSMParameter(appName, rName string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
	name = %[1]q
}
resource "aws_iam_role" "test_role" {
	name = "test_role"
  
	assume_role_policy = jsonencode({
	  Version = "2012-10-17"
	  Statement = [
		{
		  Action = "sts:AssumeRole"
		  Effect = "Allow"
		  Sid    = ""
		  Principal = {
			Service = "ec2.amazonaws.com"
		  }
		},
	  ]
	}) 
}
resource "aws_ssm_parameter" "ssm_param" {
	name  = "foo"
	type  = "String"
	value = "bar"
  }
resource "aws_appconfig_configuration_profile" "test" {
  name = %[2]q
  application_id = aws_appconfig_application.app.id
  location_uri = aws_ssm_parameter.ssm_param.arn
  retrieval_role_arn = aws_iam_role.test_role.arn
}
`, appName, rName)
}

func testAccAWSAppConfigConfigurationProfileLocationSSMDocument(appName, rName string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
	name = %[1]q
}
resource "aws_iam_role" "test_role" {
	name = "test_role"
  
	assume_role_policy = jsonencode({
	  Version = "2012-10-17"
	  Statement = [
		{
		  Action = "sts:AssumeRole"
		  Effect = "Allow"
		  Sid    = ""
		  Principal = {
			Service = "ec2.amazonaws.com"
		  }
		},
	  ]
	}) 
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
  name = %[2]q
  application_id = aws_appconfig_application.app.id
  location_uri = aws_ssm_document.ssm_doc.arn
  retrieval_role_arn = aws_iam_role.test_role.arn
}
`, appName, rName)
}

func testAccAWSAppConfigConfigurationProfileValidator(appName, rName string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
	name = %[1]q
}
resource "aws_appconfig_configuration_profile" "test" {
  name = %[2]q
  application_id = aws_appconfig_application.app.id
  location_uri = "hosted"
  validator {
    type = "JSON_SCHEMA"
    content = "JSON Schema content or AWS Lambda function name"
  }
  validator {
    type = "LAMBDA"
    content = "JSON Schema content or AWS Lambda function name"
  }
}
`, appName, rName)
}

func testAccAWSAppConfigConfigurationProfileRetreivalARN(appName, rName string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
	name = %[1]q
}
resource "aws_iam_role" "test_role" {
	name = "test_role"
  
	assume_role_policy = jsonencode({
	  Version = "2012-10-17"
	  Statement = [
		{
		  Action = "sts:AssumeRole"
		  Effect = "Allow"
		  Sid    = ""
		  Principal = {
			Service = "ec2.amazonaws.com"
		  }
		},
	  ]
	}) 
}
resource "aws_appconfig_configuration_profile" "test" {
  name = %[2]q
  application_id = aws_appconfig_application.app.id
  location_uri = "hosted"
  retrieval_role_arn = aws_iam_role.test_role.arn
}
`, appName, rName)
}
func testAccAWSAppConfigConfigurationProfileTags1(appName, rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
	resource "aws_appconfig_application" "app" {
		name = %[1]q
	}
	resource "aws_appconfig_configuration_profile" "test" {
	  name = %[2]q
	  application_id = aws_appconfig_application.app.id
	  location_uri = "hosted"
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
	  name = %[2]q
	  application_id = aws_appconfig_application.app.id
	  location_uri = "hosted"
	  tags = {
		%[3]q = %[4]q
		%[5]q = %[6]q
	  }
	}
`, appName, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

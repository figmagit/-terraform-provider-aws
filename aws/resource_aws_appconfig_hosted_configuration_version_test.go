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

func TestAccAWSAppConfigHostedConfigurationVersion_basic(t *testing.T) {
	var hostedConfigurationVersion appconfig.GetHostedConfigurationVersionOutput
	rDesc := acctest.RandomWithPrefix("desc")
	resourceName := "aws_appconfig_hosted_configuration_version.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigHostedConfigurationVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigHostedConfigurationVersionName(rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigHostedConfigurationVersionExists(resourceName, &hostedConfigurationVersion),
					testAccCheckAWSAppConfigHostedConfigurationVersionARN(resourceName, &hostedConfigurationVersion),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "content", "Settings"),
					resource.TestCheckResourceAttr(resourceName, "content_type", "text/plain"),
				),
			},
		},
	})
}

func TestAccAWSAppConfigHostedConfigurationVersion_disappears(t *testing.T) {
	var hostedConfigurationVersion appconfig.GetHostedConfigurationVersionOutput

	rDesc := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_hosted_configuration_version.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigHostedConfigurationVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigHostedConfigurationVersionName(rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigHostedConfigurationVersionExists(resourceName, &hostedConfigurationVersion),
					testAccCheckAWSAppConfigHostedConfigurationVersionDisappears(&hostedConfigurationVersion),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSAppConfigHostedConfigurationVersion_Plain(t *testing.T) {
	var hostedConfigurationVersion appconfig.GetHostedConfigurationVersionOutput
	resourceName := "aws_appconfig_hosted_configuration_version.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigHostedConfigurationVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigHostedConfigurationVersionPlainText(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigHostedConfigurationVersionExists(resourceName, &hostedConfigurationVersion),
					testAccCheckAWSAppConfigHostedConfigurationVersionARN(resourceName, &hostedConfigurationVersion),
					resource.TestCheckResourceAttr(resourceName, "version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "content", "This is a list of the new settings! \n1. A \n2. B \n3. C\n"),
					resource.TestCheckResourceAttr(resourceName, "content_type", "text/plain"),
				),
			},
		},
	})
}

func TestAccAWSAppConfigHostedConfigurationVersion_JSON(t *testing.T) {
	var hostedConfigurationVersion appconfig.GetHostedConfigurationVersionOutput
	resourceName := "aws_appconfig_hosted_configuration_version.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigHostedConfigurationVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigHostedConfigurationVersionJSON(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigHostedConfigurationVersionExists(resourceName, &hostedConfigurationVersion),
					testAccCheckAWSAppConfigHostedConfigurationVersionARN(resourceName, &hostedConfigurationVersion),
					resource.TestCheckResourceAttr(resourceName, "version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "content", `{"hello":"world"}`),
					resource.TestCheckResourceAttr(resourceName, "content_type", "application/json"),
				),
			},
		},
	})
}

func TestAccAWSAppConfigHostedConfigurationVersion_YAML(t *testing.T) {
	var hostedConfigurationVersion appconfig.GetHostedConfigurationVersionOutput
	resourceName := "aws_appconfig_hosted_configuration_version.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigHostedConfigurationVersionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigHostedConfigurationVersionYAML(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigHostedConfigurationVersionExists(resourceName, &hostedConfigurationVersion),
					testAccCheckAWSAppConfigHostedConfigurationVersionARN(resourceName, &hostedConfigurationVersion),
					resource.TestCheckResourceAttr(resourceName, "version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "content", "\"a\": \"b\"\n\"c\": \"d\"\n"),
					resource.TestCheckResourceAttr(resourceName, "content_type", "application/x-yaml"),
				),
			},
		},
	})
}

func testAccCheckAppConfigHostedConfigurationVersionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).appconfigconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_appconfig_hosted_configuration_version" {
			continue
		}

		versionNumber, err := strconv.ParseInt(rs.Primary.Attributes["version_number"], 10, 64)
		if err != nil {
			return err
		}

		input := &appconfig.GetHostedConfigurationVersionInput{
			ApplicationId:          aws.String(rs.Primary.Attributes["application_id"]),
			ConfigurationProfileId: aws.String(rs.Primary.Attributes["configuration_profile_id"]),
			VersionNumber:          aws.Int64(versionNumber),
		}

		output, err := conn.GetHostedConfigurationVersion(input)

		if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
			continue
		}

		if err != nil {
			return err
		}

		if output != nil {
			return fmt.Errorf("AppConfig HostedConfigurationVersion (%s) still exists", rs.Primary.ID)
		}
	}

	return nil

}

func testAccCheckAWSAppConfigHostedConfigurationVersionDisappears(hostedConfigurationVersion *appconfig.GetHostedConfigurationVersionOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).appconfigconn

		input := &appconfig.DeleteHostedConfigurationVersionInput{
			ApplicationId:          hostedConfigurationVersion.ApplicationId,
			ConfigurationProfileId: hostedConfigurationVersion.ConfigurationProfileId,
			VersionNumber:          hostedConfigurationVersion.VersionNumber,
		}

		_, err := conn.DeleteHostedConfigurationVersion(input)

		return err
	}
}

func testAccCheckAWSAppConfigHostedConfigurationVersionExists(resourceName string, hostedConfigurationVersion *appconfig.GetHostedConfigurationVersionOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		conn := testAccProvider.Meta().(*AWSClient).appconfigconn

		versionNumber, err := strconv.ParseInt(rs.Primary.Attributes["version_number"], 10, 64)
		if err != nil {
			return err
		}

		input := &appconfig.GetHostedConfigurationVersionInput{
			ApplicationId:          aws.String(rs.Primary.Attributes["application_id"]),
			ConfigurationProfileId: aws.String(rs.Primary.Attributes["configuration_profile_id"]),
			VersionNumber:          aws.Int64(versionNumber),
		}

		output, err := conn.GetHostedConfigurationVersion(input)
		if err != nil {
			return err
		}

		*hostedConfigurationVersion = *output

		return nil
	}
}

func testAccCheckAWSAppConfigHostedConfigurationVersionARN(resourceName string, hostedConfigurationVersion *appconfig.GetHostedConfigurationVersionOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		appID := aws.StringValue(hostedConfigurationVersion.ApplicationId)
		profileID := aws.StringValue(hostedConfigurationVersion.ConfigurationProfileId)
		versionNum := fmt.Sprintf("%d", aws.Int64Value(hostedConfigurationVersion.VersionNumber))
		arnResource := fmt.Sprintf("application/%s/configurationprofile/%s/hostedconfigurationversion/%s", appID, profileID, versionNum)
		return testAccCheckResourceAttrRegionalARN(resourceName, "arn", "appconfig", arnResource)(s)
	}
}

func testAccAWSAppConfigHostedConfigurationVersionSetup() string {
	baseName := acctest.RandomWithPrefix("tf-acc-test")
	appName := fmt.Sprintf("%s-app", baseName)
	configName := fmt.Sprintf("%s-config", baseName)
	return fmt.Sprintf(`
resource "aws_appconfig_application" "app" {
  name = %[1]q
}
resource "aws_appconfig_configuration_profile" "config" {
  application_id = aws_appconfig_application.app.id
  location_uri   = "hosted"
  name           = %[2]q
}
`, appName, configName)
}

func testAccAWSAppConfigHostedConfigurationVersionName(rDesc string) string {
	return fmt.Sprintf(`
%[1]s
resource "aws_appconfig_hosted_configuration_version" "test" {
  application_id           = aws_appconfig_application.app.id
  configuration_profile_id = aws_appconfig_configuration_profile.config.id
  content                  = "Settings"
  content_type             = "text/plain"
  description              = %[2]q
}
`, testAccAWSAppConfigHostedConfigurationVersionSetup(), rDesc)
}

func testAccAWSAppConfigHostedConfigurationVersionPlainText() string {
	return fmt.Sprintf(`
%[1]s
resource "aws_appconfig_hosted_configuration_version" "test" {
  application_id           = aws_appconfig_application.app.id
  configuration_profile_id = aws_appconfig_configuration_profile.config.id
  content                  = <<-EOF
	This is a list of the new settings! 
	1. A 
	2. B 
	3. C
	EOF
  content_type             = "text/plain"
}
`, testAccAWSAppConfigHostedConfigurationVersionSetup())
}

func testAccAWSAppConfigHostedConfigurationVersionJSON() string {
	return fmt.Sprintf(`
%[1]s
resource "aws_appconfig_hosted_configuration_version" "test" {
  application_id           = aws_appconfig_application.app.id
  configuration_profile_id = aws_appconfig_configuration_profile.config.id
  content                  = jsonencode({ "hello" = "world" })
  content_type             = "application/json"
}
`, testAccAWSAppConfigHostedConfigurationVersionSetup())
}

func testAccAWSAppConfigHostedConfigurationVersionYAML() string {
	return fmt.Sprintf(`
%[1]s
resource "aws_appconfig_hosted_configuration_version" "test" {
  application_id           = aws_appconfig_application.app.id
  configuration_profile_id = aws_appconfig_configuration_profile.config.id
  content                  = yamlencode({ "a" : "b", "c" : "d" })
  content_type             = "application/x-yaml"
}
`, testAccAWSAppConfigHostedConfigurationVersionSetup())
}

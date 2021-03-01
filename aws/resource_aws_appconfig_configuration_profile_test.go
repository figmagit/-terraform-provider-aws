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

	rName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileName(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckResourceAttrRegionalARN(resourceName, "arn", "appconfig", fmt.Sprintf("application/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAwsAppConfigConfigurationProfile_disappears(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput

	rName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileName(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					testAccCheckAwsAppConfigConfigurationProfileDisappears(&profile),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAwsAppConfigConfigurationProfile_Tags(t *testing.T) {
	var profile appconfig.GetConfigurationProfileOutput

	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_configuration_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigConfigurationProfileTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAWSAppConfigConfigurationProfileTags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAWSAppConfigConfigurationProfileTags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsAppConfigConfigurationProfileExists(resourceName, &profile),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
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

func testAccAWSAppConfigConfigurationProfileName(rName, rDesc string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_configuration_profile" "test" {
  name = %[1]q
  description = %[2]q
}
`, rName, rDesc)
}

func testAccAWSAppConfigConfigurationProfileTags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_configuration_profile" "test" {
  analyzer_name = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccAWSAppConfigConfigurationProfileTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_configuration_profile" "test" {
  analyzer_name = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

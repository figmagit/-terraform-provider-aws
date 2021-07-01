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

func TestAccAWSAppConfigDeploymentStrategy_basic(t *testing.T) {
	var deploymentStrategy appconfig.GetDeploymentStrategyOutput
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("desc")
	resourceName := "aws_appconfig_deployment_strategy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigDeploymentStrategyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigDeploymentStrategyName(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentStrategyExists(resourceName, &deploymentStrategy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigDeploymentStrategyARN(resourceName, &deploymentStrategy),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "description", rDesc),
					resource.TestCheckResourceAttr(resourceName, "growth_factor", "1"),
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

func TestAccAWSAppConfigDeploymentStrategy_disappears(t *testing.T) {
	var deploymentStrategy appconfig.GetDeploymentStrategyOutput

	rName := acctest.RandomWithPrefix("tf-acc-test")
	rDesc := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_deployment_strategy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigDeploymentStrategyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigDeploymentStrategyName(rName, rDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentStrategyExists(resourceName, &deploymentStrategy),
					testAccCheckAWSAppConfigDeploymentStrategyDisappears(&deploymentStrategy),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAWSAppConfigDeploymentStrategy_Tags(t *testing.T) {
	var deploymentStrategy appconfig.GetDeploymentStrategyOutput

	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_deployment_strategy.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigDeploymentStrategyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigDeploymentStrategyTags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentStrategyExists(resourceName, &deploymentStrategy),
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
				Config: testAccAWSAppConfigDeploymentStrategyTags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentStrategyExists(resourceName, &deploymentStrategy),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAWSAppConfigDeploymentStrategyTags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentStrategyExists(resourceName, &deploymentStrategy),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccAWSAppConfigDeploymentStrategy_Growth(t *testing.T) {
	var deploymentStrategy appconfig.GetDeploymentStrategyOutput
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_deployment_strategy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigDeploymentStrategyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigDeploymentStrategyGrowth(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentStrategyExists(resourceName, &deploymentStrategy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigDeploymentStrategyARN(resourceName, &deploymentStrategy),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "growth_factor", "24.5"),
					resource.TestCheckResourceAttr(resourceName, "growth_type", "EXPONENTIAL"),
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
func TestAccAWSAppConfigDeploymentStrategy_ReplicateTo(t *testing.T) {
	var deploymentStrategy appconfig.GetDeploymentStrategyOutput
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_deployment_strategy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigDeploymentStrategyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigDeploymentStrategyReplicateTo(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentStrategyExists(resourceName, &deploymentStrategy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigDeploymentStrategyARN(resourceName, &deploymentStrategy),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "replicate_to", "SSM_DOCUMENT"),
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
func TestAccAWSAppConfigDeploymentStrategy_BakeTime(t *testing.T) {
	var deploymentStrategy appconfig.GetDeploymentStrategyOutput
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_appconfig_deployment_strategy.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAppConfigDeploymentStrategyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSAppConfigDeploymentStrategyBakeTime(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSAppConfigDeploymentStrategyExists(resourceName, &deploymentStrategy),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckAWSAppConfigDeploymentStrategyARN(resourceName, &deploymentStrategy),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "final_bake_time_in_minutes", "45"),
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

func testAccCheckAppConfigDeploymentStrategyDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).appconfigconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_appconfig_deployment_strategy" {
			continue
		}

		input := &appconfig.GetDeploymentStrategyInput{
			DeploymentStrategyId: aws.String(rs.Primary.ID),
		}

		output, err := conn.GetDeploymentStrategy(input)

		if isAWSErr(err, appconfig.ErrCodeResourceNotFoundException, "") {
			continue
		}

		if err != nil {
			return err
		}

		if output != nil {
			return fmt.Errorf("AppConfig DeploymentStrategy (%s) still exists", rs.Primary.ID)
		}
	}

	return nil

}

func testAccCheckAWSAppConfigDeploymentStrategyDisappears(deploymentStrategy *appconfig.GetDeploymentStrategyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).appconfigconn

		input := &appconfig.DeleteDeploymentStrategyInput{
			DeploymentStrategyId: aws.String(*deploymentStrategy.Id),
		}

		_, err := conn.DeleteDeploymentStrategy(input)

		return err
	}
}

func testAccCheckAWSAppConfigDeploymentStrategyExists(resourceName string, deploymentStrategy *appconfig.GetDeploymentStrategyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		conn := testAccProvider.Meta().(*AWSClient).appconfigconn

		input := &appconfig.GetDeploymentStrategyInput{
			DeploymentStrategyId: aws.String(rs.Primary.ID),
		}

		output, err := conn.GetDeploymentStrategy(input)
		if err != nil {
			return err
		}

		*deploymentStrategy = *output

		return nil
	}
}

func testAccCheckAWSAppConfigDeploymentStrategyARN(resourceName string, deploymentStrategy *appconfig.GetDeploymentStrategyOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return testAccCheckResourceAttrRegionalARN(resourceName, "arn", "appconfig", fmt.Sprintf("deploymentstrategy/%s", aws.StringValue(deploymentStrategy.Id)))(s)
	}
}

func testAccAWSAppConfigDeploymentStrategyName(rName, rDesc string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_deployment_strategy" "test" {
  name                           = %[1]q
  description                    = %[2]q
  deployment_duration_in_minutes = 10
  growth_type                    = "LINEAR"
  replicate_to                   = "NONE"
}
`, rName, rDesc)
}

func testAccAWSAppConfigDeploymentStrategyGrowth(rName string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_deployment_strategy" "test" {
  name                           = %[1]q
  deployment_duration_in_minutes = 10
  growth_factor                  = 24.5
  growth_type                    = "EXPONENTIAL"
  replicate_to                   = "NONE"
}
`, rName)
}

func testAccAWSAppConfigDeploymentStrategyReplicateTo(rName string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_deployment_strategy" "test" {
  name                           = %[1]q
  deployment_duration_in_minutes = 10
  growth_type                    = "LINEAR"
  replicate_to                   = "SSM_DOCUMENT"
}
`, rName)
}

func testAccAWSAppConfigDeploymentStrategyBakeTime(rName string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_deployment_strategy" "test" {
  name                           = %[1]q
  deployment_duration_in_minutes = 10
  final_bake_time_in_minutes     = 45
  growth_type                    = "LINEAR"
  replicate_to                   = "NONE"
}
`, rName)
}

func testAccAWSAppConfigDeploymentStrategyTags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_deployment_strategy" "test" {
  name                           = %[1]q
  deployment_duration_in_minutes = 10
  growth_type                    = "LINEAR"
  replicate_to                   = "NONE"

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccAWSAppConfigDeploymentStrategyTags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_appconfig_deployment_strategy" "test" {
  name                           = %[1]q
  deployment_duration_in_minutes = 10
  growth_type                    = "LINEAR"
  replicate_to                   = "NONE"

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}

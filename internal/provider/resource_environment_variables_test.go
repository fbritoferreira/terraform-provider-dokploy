package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnvironmentVariablesResource(t *testing.T) {
	host := os.Getenv("DOKPLOY_HOST")
	apiKey := os.Getenv("DOKPLOY_API_KEY")

	if host == "" || apiKey == "" {
		t.Skip("DOKPLOY_HOST and DOKPLOY_API_KEY must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEnvironmentVariablesResourceConfig("test-env-vars-project", "test-env-vars-env", "test-env-vars-app"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dokploy_environment_variables.test", "variables.ENV1", "value1"),
					resource.TestCheckResourceAttr("dokploy_environment_variables.test", "variables.ENV2", "value2"),
					resource.TestCheckResourceAttrSet("dokploy_environment_variables.test", "id"),
					resource.TestCheckResourceAttrSet("dokploy_environment_variables.test", "application_id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccEnvironmentVariablesResourceConfigUpdated("test-env-vars-project", "test-env-vars-env", "test-env-vars-app"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dokploy_environment_variables.test", "variables.ENV1", "updated_value1"),
					resource.TestCheckResourceAttr("dokploy_environment_variables.test", "variables.ENV3", "value3"),
					resource.TestCheckNoResourceAttr("dokploy_environment_variables.test", "variables.ENV2"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "dokploy_environment_variables.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_env_file"},
			},
		},
	})
}

func testAccEnvironmentVariablesResourceConfig(projectName, envName, appName string) string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

resource "dokploy_project" "test" {
  name        = "%s"
  description = "Test project for environment variables tests"
}

resource "dokploy_environment" "test" {
  project_id = dokploy_project.test.id
  name       = "%s"
}

resource "dokploy_application" "test" {
  project_id     = dokploy_project.test.id
  environment_id = dokploy_environment.test.id
  name           = "%s"
  build_type     = "nixpacks"
  source_type    = "docker"
  docker_image   = "nginx:latest"
}

resource "dokploy_environment_variables" "test" {
  application_id = dokploy_application.test.id
  variables = {
    ENV1 = "value1"
    ENV2 = "value2"
  }
}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"), projectName, envName, appName)
}

func testAccEnvironmentVariablesResourceConfigUpdated(projectName, envName, appName string) string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

resource "dokploy_project" "test" {
  name        = "%s"
  description = "Test project for environment variables tests"
}

resource "dokploy_environment" "test" {
  project_id = dokploy_project.test.id
  name       = "%s"
}

resource "dokploy_application" "test" {
  project_id     = dokploy_project.test.id
  environment_id = dokploy_environment.test.id
  name           = "%s"
  build_type     = "nixpacks"
  source_type    = "docker"
  docker_image   = "nginx:latest"
}

resource "dokploy_environment_variables" "test" {
  application_id = dokploy_application.test.id
  variables = {
    ENV1 = "updated_value1"
    ENV3 = "value3"
  }
}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"), projectName, envName, appName)
}

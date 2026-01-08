package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRedisResource(t *testing.T) {
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
				Config: testAccRedisResourceConfig("test-redis-project", "test-redis-env", "test-redis", "testredisapp"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dokploy_redis.test", "name", "test-redis"),
					resource.TestCheckResourceAttrSet("dokploy_redis.test", "id"),
					resource.TestCheckResourceAttrSet("dokploy_redis.test", "environment_id"),
					resource.TestCheckResourceAttrSet("dokploy_redis.test", "app_name"),
				),
			},
			// Update and Read testing
			{
				Config: testAccRedisResourceConfigWithDescription("test-redis-project", "test-redis-env", "test-redis-updated", "testredisapp", "Updated Redis instance"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dokploy_redis.test", "name", "test-redis-updated"),
					resource.TestCheckResourceAttr("dokploy_redis.test", "description", "Updated Redis instance"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "dokploy_redis.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"database_password"}, // Password not returned by API.
			},
		},
	})
}

func testAccRedisResourceConfig(projectName, envName, redisName, appName string) string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

resource "dokploy_project" "test" {
  name        = "%s"
  description = "Test project for Redis tests"
}

resource "dokploy_environment" "test" {
  project_id = dokploy_project.test.id
  name       = "%s"
}

resource "dokploy_redis" "test" {
  name              = "%s"
  app_name          = "%s"
  database_password = "test_redis_password_123"
  environment_id    = dokploy_environment.test.id
}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"), projectName, envName, redisName, appName)
}

func testAccRedisResourceConfigWithDescription(projectName, envName, redisName, appName string, description string) string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

resource "dokploy_project" "test" {
  name        = "%s"
  description = "Test project for Redis tests"
}

resource "dokploy_environment" "test" {
  project_id = dokploy_project.test.id
  name       = "%s"
}

resource "dokploy_redis" "test" {
  name              = "%s"
  app_name          = "%s"
  database_password = "test_redis_password_123"
  environment_id    = dokploy_environment.test.id
  description       = "%s"
}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"), projectName, envName, redisName, appName, description)
}

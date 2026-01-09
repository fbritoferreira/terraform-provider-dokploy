package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApiKeyResource(t *testing.T) {
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
				Config: testAccApiKeyResourceConfig("test-terraform-api-key"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dokploy_api_key.test", "name", "test-terraform-api-key"),
					resource.TestCheckResourceAttrSet("dokploy_api_key.test", "id"),
					resource.TestCheckResourceAttrSet("dokploy_api_key.test", "key"),
					resource.TestCheckResourceAttrSet("dokploy_api_key.test", "start"),
					resource.TestCheckResourceAttrSet("dokploy_api_key.test", "user_id"),
					resource.TestCheckResourceAttrSet("dokploy_api_key.test", "organization_id"),
					resource.TestCheckResourceAttrSet("dokploy_api_key.test", "created_at"),
					resource.TestCheckResourceAttr("dokploy_api_key.test", "enabled", "true"),
				),
			},
			// ImportState testing - note that key value won't be available after import
			{
				ResourceName:            "dokploy_api_key.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key", "expires_in", "organization_id"}, // Key and org_id are not returned on read
			},
		},
	})
}

func TestAccApiKeyResourceWithExpiry(t *testing.T) {
	host := os.Getenv("DOKPLOY_HOST")
	apiKey := os.Getenv("DOKPLOY_API_KEY")

	if host == "" || apiKey == "" {
		t.Skip("DOKPLOY_HOST and DOKPLOY_API_KEY must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with expiry
			{
				Config: testAccApiKeyResourceConfigWithExpiry("test-terraform-api-key-expiry", 86400),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dokploy_api_key.test", "name", "test-terraform-api-key-expiry"),
					resource.TestCheckResourceAttr("dokploy_api_key.test", "expires_in", "86400"),
					resource.TestCheckResourceAttrSet("dokploy_api_key.test", "expires_at"),
					resource.TestCheckResourceAttrSet("dokploy_api_key.test", "key"),
				),
			},
		},
	})
}

func TestAccApiKeyResourceWithRateLimit(t *testing.T) {
	host := os.Getenv("DOKPLOY_HOST")
	apiKey := os.Getenv("DOKPLOY_API_KEY")

	if host == "" || apiKey == "" {
		t.Skip("DOKPLOY_HOST and DOKPLOY_API_KEY must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with rate limit
			{
				Config: testAccApiKeyResourceConfigWithRateLimit("test-terraform-api-key-ratelimit", 100, 3600000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dokploy_api_key.test", "name", "test-terraform-api-key-ratelimit"),
					resource.TestCheckResourceAttr("dokploy_api_key.test", "rate_limit_enabled", "true"),
					resource.TestCheckResourceAttr("dokploy_api_key.test", "rate_limit_max", "100"),
					resource.TestCheckResourceAttr("dokploy_api_key.test", "rate_limit_time_window", "3600000"),
				),
			},
		},
	})
}

func testAccApiKeyResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

resource "dokploy_api_key" "test" {
  name = "%s"
}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"), name)
}

func testAccApiKeyResourceConfigWithExpiry(name string, expiresIn int) string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

resource "dokploy_api_key" "test" {
  name       = "%s"
  expires_in = %d
}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"), name, expiresIn)
}

func testAccApiKeyResourceConfigWithRateLimit(name string, rateLimitMax int, rateLimitTimeWindow int) string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

resource "dokploy_api_key" "test" {
  name                   = "%s"
  rate_limit_enabled     = true
  rate_limit_max         = %d
  rate_limit_time_window = %d
}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"), name, rateLimitMax, rateLimitTimeWindow)
}

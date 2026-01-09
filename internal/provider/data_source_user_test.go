package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	host := os.Getenv("DOKPLOY_HOST")
	apiKey := os.Getenv("DOKPLOY_API_KEY")

	if host == "" || apiKey == "" {
		t.Skip("DOKPLOY_HOST and DOKPLOY_API_KEY must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.dokploy_user.current", "id"),
					resource.TestCheckResourceAttrSet("data.dokploy_user.current", "user_id"),
					resource.TestCheckResourceAttrSet("data.dokploy_user.current", "member_id"),
					resource.TestCheckResourceAttrSet("data.dokploy_user.current", "organization_id"),
					resource.TestCheckResourceAttrSet("data.dokploy_user.current", "role"),
					resource.TestCheckResourceAttrSet("data.dokploy_user.current", "email"),
					resource.TestCheckResourceAttrSet("data.dokploy_user.current", "created_at"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig() string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

data "dokploy_user" "current" {}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"))
}

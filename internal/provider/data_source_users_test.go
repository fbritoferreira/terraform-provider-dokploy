package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUsersDataSource(t *testing.T) {
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
				Config: testAccUsersDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.dokploy_users.all", "users.#"),
					// At least one user should exist (the current user)
					resource.TestCheckResourceAttrSet("data.dokploy_users.all", "users.0.member_id"),
					resource.TestCheckResourceAttrSet("data.dokploy_users.all", "users.0.user_id"),
					resource.TestCheckResourceAttrSet("data.dokploy_users.all", "users.0.organization_id"),
					resource.TestCheckResourceAttrSet("data.dokploy_users.all", "users.0.email"),
				),
			},
		},
	})
}

func testAccUsersDataSourceConfig() string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

data "dokploy_users" "all" {}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"))
}

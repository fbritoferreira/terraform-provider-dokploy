package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationResource(t *testing.T) {
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
				Config: testAccOrganizationResourceConfig("test-terraform-org"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dokploy_organization.test", "name", "test-terraform-org"),
					resource.TestCheckResourceAttrSet("dokploy_organization.test", "id"),
					resource.TestCheckResourceAttrSet("dokploy_organization.test", "owner_id"),
					resource.TestCheckResourceAttrSet("dokploy_organization.test", "created_at"),
				),
			},
			// Update and Read testing
			{
				Config: testAccOrganizationResourceConfig("test-terraform-org-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dokploy_organization.test", "name", "test-terraform-org-updated"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "dokploy_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOrganizationResourceWithLogo(t *testing.T) {
	host := os.Getenv("DOKPLOY_HOST")
	apiKey := os.Getenv("DOKPLOY_API_KEY")

	if host == "" || apiKey == "" {
		t.Skip("DOKPLOY_HOST and DOKPLOY_API_KEY must be set for acceptance tests")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with logo
			{
				Config: testAccOrganizationResourceConfigWithLogo("test-terraform-org-logo", "https://example.com/logo.png"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dokploy_organization.test", "name", "test-terraform-org-logo"),
					resource.TestCheckResourceAttr("dokploy_organization.test", "logo", "https://example.com/logo.png"),
					resource.TestCheckResourceAttrSet("dokploy_organization.test", "id"),
				),
			},
			// Update logo
			{
				Config: testAccOrganizationResourceConfigWithLogo("test-terraform-org-logo", "https://example.com/new-logo.png"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("dokploy_organization.test", "logo", "https://example.com/new-logo.png"),
				),
			},
		},
	})
}

func TestAccOrganizationsDataSource(t *testing.T) {
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
				Config: testAccOrganizationsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.dokploy_organizations.all", "organizations.#"),
				),
			},
		},
	})
}

func testAccOrganizationResourceConfig(name string) string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

resource "dokploy_organization" "test" {
  name = "%s"
}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"), name)
}

func testAccOrganizationResourceConfigWithLogo(name, logo string) string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

resource "dokploy_organization" "test" {
  name = "%s"
  logo = "%s"
}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"), name, logo)
}

func testAccOrganizationsDataSourceConfig() string {
	return fmt.Sprintf(`
provider "dokploy" {
  host    = "%s"
  api_key = "%s"
}

data "dokploy_organizations" "all" {}
`, os.Getenv("DOKPLOY_HOST"), os.Getenv("DOKPLOY_API_KEY"))
}

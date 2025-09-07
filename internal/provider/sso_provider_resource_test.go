package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/supabase/cli/pkg/api"
	"gopkg.in/h2non/gock.v1"
)

func TestAccSsoProviderResource(t *testing.T) {
	defer gock.Off()

	// Mock API responses for SSO provider CRUD operations
	gock.New("https://api.supabase.com").
		Post("/v1/projects/mayuaycdtijbctgqbycg/config/auth/sso/providers").
		Reply(http.StatusCreated).
		JSON(api.CreateProviderResponse{
			Id:        "provider-123",
			CreatedAt: Ptr("2024-01-01T00:00:00Z"),
			UpdatedAt: Ptr("2024-01-01T00:00:00Z"),
			Domains: &[]api.Domain{
				{Domain: Ptr("example.com")},
			},
		})

	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth/sso/providers/provider-123").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.GetProviderResponse{
			Id:        "provider-123",
			CreatedAt: Ptr("2024-01-01T00:00:00Z"),
			UpdatedAt: Ptr("2024-01-01T00:00:00Z"),
			Domains: &[]api.Domain{
				{Domain: Ptr("example.com")},
			},
		})

	gock.New("https://api.supabase.com").
		Put("/v1/projects/mayuaycdtijbctgqbycg/config/auth/sso/providers/provider-123").
		Reply(http.StatusOK).
		JSON(api.UpdateProviderResponse{
			UpdatedAt: Ptr("2024-01-01T01:00:00Z"),
			Domains: &[]api.Domain{
				{Domain: Ptr("example.com")},
				{Domain: Ptr("example2.com")},
			},
		})

	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth/sso/providers/provider-123").
		Reply(http.StatusOK).
		JSON(api.GetProviderResponse{
			Id:        "provider-123",
			CreatedAt: Ptr("2024-01-01T00:00:00Z"),
			UpdatedAt: Ptr("2024-01-01T01:00:00Z"),
			Domains: &[]api.Domain{
				{Domain: Ptr("example.com")},
				{Domain: Ptr("example2.com")},
			},
		})

	gock.New("https://api.supabase.com").
		Delete("/v1/projects/mayuaycdtijbctgqbycg/config/auth/sso/providers/provider-123").
		Reply(http.StatusOK).
		JSON(api.DeleteProviderResponse{})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSsoProviderResourceConfig("example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_sso_provider.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_sso_provider.test", "type", "saml"),
					resource.TestCheckResourceAttr("supabase_sso_provider.test", "metadata_url", "https://example.com/metadata"),
					resource.TestCheckResourceAttr("supabase_sso_provider.test", "domains.#", "1"),
					resource.TestCheckResourceAttr("supabase_sso_provider.test", "domains.0", "example.com"),
					resource.TestCheckResourceAttr("supabase_sso_provider.test", "id", "provider-123"),
				),
			},
			// Update and Read testing
			{
				Config: testAccSsoProviderResourceConfigUpdate("example.com", "example2.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_sso_provider.test", "domains.#", "2"),
					resource.TestCheckResourceAttr("supabase_sso_provider.test", "domains.0", "example.com"),
					resource.TestCheckResourceAttr("supabase_sso_provider.test", "domains.1", "example2.com"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccSsoProviderResourceConfig(domain string) string {
	return `
resource "supabase_sso_provider" "test" {
  project_ref  = "mayuaycdtijbctgqbycg"
  type         = "saml"
  metadata_url = "https://example.com/metadata"
  domains      = ["` + domain + `"]
}
`
}

func testAccSsoProviderResourceConfigUpdate(domain1, domain2 string) string {
	return `
resource "supabase_sso_provider" "test" {
  project_ref  = "mayuaycdtijbctgqbycg"
  type         = "saml"
  metadata_url = "https://example.com/metadata"
  domains      = ["` + domain1 + `", "` + domain2 + `"]
}
`
}
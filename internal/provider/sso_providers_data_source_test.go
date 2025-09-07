package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/supabase/cli/pkg/api"
	"gopkg.in/h2non/gock.v1"
)

func TestAccSsoProvidersDataSource(t *testing.T) {
	defer gock.Off()

	// Mock API response for listing SSO providers
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth/sso/providers").
		Times(3). // Called multiple times during the test
		Reply(http.StatusOK).
		JSON(api.ListProvidersResponse{
			Items: []api.Provider{
			{
				Id:        "provider-123",
				CreatedAt: Ptr("2024-01-01T00:00:00Z"),
				UpdatedAt: Ptr("2024-01-01T00:00:00Z"),
				Saml: &api.SamlDescriptor{
					EntityId: "https://example.com/entity",
				},
				Domains: &[]api.Domain{
					{Domain: Ptr("example.com")},
					{Domain: Ptr("example2.com")},
				},
			},
			{
				Id:        "provider-456",
				CreatedAt: Ptr("2024-01-02T00:00:00Z"),
				UpdatedAt: Ptr("2024-01-02T00:00:00Z"),
				Saml: &api.SamlDescriptor{
					EntityId: "https://example2.com/entity",
				},
				Domains: &[]api.Domain{
					{Domain: Ptr("example3.com")},
				},
			},
			},
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSsoProvidersDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.supabase_sso_providers.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("data.supabase_sso_providers.test", "providers.#", "2"),
					resource.TestCheckResourceAttr("data.supabase_sso_providers.test", "providers.0.id", "provider-123"),
					resource.TestCheckResourceAttr("data.supabase_sso_providers.test", "providers.0.type", "saml"),
					resource.TestCheckResourceAttr("data.supabase_sso_providers.test", "providers.0.domains.#", "2"),
					resource.TestCheckResourceAttr("data.supabase_sso_providers.test", "providers.0.domains.0", "example.com"),
					resource.TestCheckResourceAttr("data.supabase_sso_providers.test", "providers.0.domains.1", "example2.com"),
					resource.TestCheckResourceAttr("data.supabase_sso_providers.test", "providers.1.id", "provider-456"),
					resource.TestCheckResourceAttr("data.supabase_sso_providers.test", "providers.1.type", "saml"),
					resource.TestCheckResourceAttr("data.supabase_sso_providers.test", "providers.1.domains.#", "1"),
					resource.TestCheckResourceAttr("data.supabase_sso_providers.test", "providers.1.domains.0", "example3.com"),
				),
			},
		},
	})
}

const testAccSsoProvidersDataSourceConfig = `
data "supabase_sso_providers" "test" {
  project_ref = "mayuaycdtijbctgqbycg"
}
`
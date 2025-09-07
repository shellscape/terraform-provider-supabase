package provider

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"gopkg.in/h2non/gock.v1"
)

func TestAccDatabaseWebhookResource(t *testing.T) {
	defer gock.Off()

	// Mock API responses for webhook enable
	gock.New("https://api.supabase.com").
		Post("/v1/projects/mayuaycdtijbctgqbycg/database/webhooks/enable").
		Reply(http.StatusCreated)

	gock.New("https://api.supabase.com").
		Post("/v1/projects/mayuaycdtijbctgqbycg/database/webhooks/enable").
		Reply(http.StatusOK) // For update operation

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDatabaseWebhookResourceConfig(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_database_webhook.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_database_webhook.test", "enabled", "true"),
					resource.TestCheckResourceAttr("supabase_database_webhook.test", "id", "mayuaycdtijbctgqbycg-webhook"),
				),
			},
			// Update testing
			{
				Config: testAccDatabaseWebhookResourceConfig(true), // Keep enabled
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_database_webhook.test", "enabled", "true"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDatabaseWebhookResourceDisabled(t *testing.T) {
	defer gock.Off()

	// Test creating with enabled = false (should be a no-op)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseWebhookResourceConfig(false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_database_webhook.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_database_webhook.test", "enabled", "false"),
					resource.TestCheckResourceAttr("supabase_database_webhook.test", "id", "mayuaycdtijbctgqbycg-webhook"),
				),
			},
		},
	})
}

func testAccDatabaseWebhookResourceConfig(enabled bool) string {
	return `
resource "supabase_database_webhook" "test" {
  project_ref = "mayuaycdtijbctgqbycg"
  enabled     = ` + fmt.Sprintf("%t", enabled) + `
}
`
}
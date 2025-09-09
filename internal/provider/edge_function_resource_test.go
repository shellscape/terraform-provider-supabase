package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/supabase/cli/pkg/api"
	"gopkg.in/h2non/gock.v1"
)

func TestAccEdgeFunctionResource(t *testing.T) {
	// Setup mock API responses - no CLI bundling
	defer gock.OffAll()
	
	// Mock function creation
	gock.New("https://api.supabase.com").
		Post("/v1/projects/mayuaycdtijbctgqbycg/functions").
		Reply(http.StatusCreated).
		JSON(&api.FunctionResponse{
			Id:        "test-func-id",
			Slug:      "test-function",
			Name:      "Test Function",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995200,
		})

	// Mock function read - first two times return original name
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/test-function").
		Times(2).
		Reply(http.StatusOK).
		JSON(&api.FunctionSlugResponse{
			Id:        "test-func-id",
			Slug:      "test-function",
			Name:      "Test Function",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995200,
			Version:   1,
		})

	// Mock function read after update - return updated name
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/test-function").
		Times(2).
		Reply(http.StatusOK).
		JSON(&api.FunctionSlugResponse{
			Id:        "test-func-id",
			Slug:      "test-function",
			Name:      "Updated Function",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995300,
			Version:   2,
		})

	// Mock function update
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/functions/test-function").
		Reply(http.StatusOK).
		JSON(&api.FunctionResponse{
			Id:        "test-func-id",
			Slug:      "test-function", 
			Name:      "Updated Function",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995300,
		})

	// Mock function delete
	gock.New("https://api.supabase.com").
		Delete("/v1/projects/mayuaycdtijbctgqbycg/functions/test-function").
		Reply(http.StatusOK).
		JSON(map[string]string{"message": "Function deleted successfully"})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create testing
			{
				Config: `
resource "supabase_edge_function" "test" {
  project_ref      = "mayuaycdtijbctgqbycg"
  slug             = "test-function"
  name             = "Test Function"
  entrypoint_path  = "/tmp/index.ts"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_edge_function.test", "id", "test-func-id"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "slug", "test-function"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "name", "Test Function"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet("supabase_edge_function.test", "created_at"),
				),
			},
			// Update testing
			{
				Config: `
resource "supabase_edge_function" "test" {
  project_ref      = "mayuaycdtijbctgqbycg"
  slug             = "test-function"
  name             = "Updated Function"
  entrypoint_path  = "/tmp/index.ts"
}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_edge_function.test", "name", "Updated Function"),
				),
			},
		},
	})
}
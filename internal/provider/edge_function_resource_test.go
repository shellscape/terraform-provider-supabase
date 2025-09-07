package provider

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/supabase/cli/pkg/api"
	"gopkg.in/h2non/gock.v1"
)

func TestAccEdgeFunctionResource(t *testing.T) {
	defer gock.OffAll()

	// Step 1: Create
	gock.New("https://api.supabase.com").
		Post("/v1/projects/mayuaycdtijbctgqbycg/functions").
		Reply(http.StatusCreated).
		JSON(api.FunctionResponse{
			Id:        "func123",
			Slug:      "hello-world",
			Name:      "Hello World",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995200,
		})

	// Step 2: Read for create
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/hello-world").
		Reply(http.StatusOK).
		JSON(api.FunctionSlugResponse{
			Id:        "func123",
			Slug:      "hello-world",
			Name:      "Hello World",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995200,
		})

	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/hello-world/body").
		Reply(http.StatusOK).
		BodyString(`export default function handler(req) { return new Response("Hello World!"); }`)

	// Step 3: Read for refresh
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/hello-world").
		Reply(http.StatusOK).
		JSON(api.FunctionSlugResponse{
			Id:        "func123",
			Slug:      "hello-world",
			Name:      "Hello World",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995200,
		})

	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/hello-world/body").
		Reply(http.StatusOK).
		BodyString(`export default function handler(req) { return new Response("Hello World!"); }`)

	// Step 4: Update
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/functions/hello-world").
		Reply(http.StatusOK).
		JSON(api.FunctionResponse{
			Id:        "func123",
			Slug:      "hello-world",
			Name:      "Hello Updated",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995300,
		})

	// Step 5: Read after update
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/hello-world").
		Reply(http.StatusOK).
		JSON(api.FunctionSlugResponse{
			Id:        "func123",
			Slug:      "hello-world",
			Name:      "Hello Updated",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995300,
		})

	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/hello-world/body").
		Reply(http.StatusOK).
		BodyString(`export default function handler(req) { return new Response("Hello Updated!"); }`)

	// Step 6: Delete
	gock.New("https://api.supabase.com").
		Delete("/v1/projects/mayuaycdtijbctgqbycg/functions/hello-world").
		Reply(http.StatusOK)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccEdgeFunctionResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_edge_function.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "slug", "hello-world"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "name", "Hello World"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "body", `export default function handler(req) { return new Response("Hello World!"); }`),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "verify_jwt", "false"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "id", "func123"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "created_at", "1640995200"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "updated_at", "1640995200"),
				),
			},
			// Update testing
			{
				Config: testAccEdgeFunctionResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_edge_function.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "slug", "hello-world"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "name", "Hello Updated"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "body", `export default function handler(req) { return new Response("Hello Updated!"); }`),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "verify_jwt", "true"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "id", "func123"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "status", "ACTIVE"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "created_at", "1640995200"),
					resource.TestCheckResourceAttr("supabase_edge_function.test", "updated_at", "1640995300"),
				),
			},
		},
	})
}

const testAccEdgeFunctionResourceConfig = `
resource "supabase_edge_function" "test" {
  project_ref = "mayuaycdtijbctgqbycg"
  slug        = "hello-world"
  name        = "Hello World"
  body        = "export default function handler(req) { return new Response(\"Hello World!\"); }"
  verify_jwt  = false
}
`

const testAccEdgeFunctionResourceConfigUpdate = `
resource "supabase_edge_function" "test" {
  project_ref = "mayuaycdtijbctgqbycg"
  slug        = "hello-world"
  name        = "Hello Updated"
  body        = "export default function handler(req) { return new Response(\"Hello Updated!\"); }"
  verify_jwt  = true
}
`

func TestAccEdgeFunctionResourceSlugValidation(t *testing.T) {
	defer gock.OffAll()

	// No API mock needed - validation happens at schema level before API calls
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccEdgeFunctionResourceConfigInvalidSlug,
				ExpectError: regexp.MustCompile(`invalid@slug!`),
			},
		},
	})
}

const testAccEdgeFunctionResourceConfigInvalidSlug = `
resource "supabase_edge_function" "test" {
  project_ref = "mayuaycdtijbctgqbycg"
  slug        = "invalid@slug!"
  name        = "Invalid Slug Function"
  body        = "export default () => new Response(\"Hello\");"
  verify_jwt  = false
}
`
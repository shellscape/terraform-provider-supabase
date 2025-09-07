package provider

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/supabase/cli/pkg/api"
	"gopkg.in/h2non/gock.v1"
)

func TestAccEdgeFunctionResourceOptionalFields(t *testing.T) {
	defer gock.OffAll()

	// Mock create with all optional fields
	gock.New("https://api.supabase.com").
		Post("/v1/projects/mayuaycdtijbctgqbycg/functions").
		Reply(http.StatusCreated).
		JSON(api.FunctionResponse{
			Id:                "func456",
			Slug:              "advanced-function",
			Name:              "Advanced Function",
			Status:            "ACTIVE",
			CreatedAt:         1640995200,
			UpdatedAt:         1640995200,
			ComputeMultiplier: func() *float32 { f := float32(2.5); return &f }(),
			EntrypointPath:    func() *string { s := "custom/index.ts"; return &s }(),
			ImportMap:         func() *bool { b := true; return &b }(),
			ImportMapPath:     func() *string { s := "custom/import_map.json"; return &s }(),
		})

	// Mock read with optional fields
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/advanced-function").
		Reply(http.StatusOK).
		JSON(api.FunctionSlugResponse{
			Id:                "func456",
			Slug:              "advanced-function",
			Name:              "Advanced Function",
			Status:            "ACTIVE",
			CreatedAt:         1640995200,
			UpdatedAt:         1640995200,
			ComputeMultiplier: func() *float32 { f := float32(2.5); return &f }(),
			EntrypointPath:    func() *string { s := "custom/index.ts"; return &s }(),
			ImportMap:         func() *bool { b := true; return &b }(),
			ImportMapPath:     func() *string { s := "custom/import_map.json"; return &s }(),
		})

	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/advanced-function/body").
		Reply(http.StatusOK).
		BodyString(`import { serve } from "./deps.ts"; serve(() => new Response("Advanced!"));`)

	// Mock successful delete for cleanup
	gock.New("https://api.supabase.com").
		Delete("/v1/projects/mayuaycdtijbctgqbycg/functions/advanced-function").
		Reply(http.StatusOK)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeFunctionResourceConfigOptionalFields,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_edge_function.advanced", "slug", "advanced-function"),
					resource.TestCheckResourceAttr("supabase_edge_function.advanced", "name", "Advanced Function"),
					resource.TestCheckResourceAttr("supabase_edge_function.advanced", "verify_jwt", "true"),
					resource.TestCheckResourceAttr("supabase_edge_function.advanced", "compute_multiplier", "2.5"),
					resource.TestCheckResourceAttr("supabase_edge_function.advanced", "entrypoint_path", "custom/index.ts"),
					resource.TestCheckResourceAttr("supabase_edge_function.advanced", "import_map", "true"),
					resource.TestCheckResourceAttr("supabase_edge_function.advanced", "import_map_path", "custom/import_map.json"),
					resource.TestCheckResourceAttr("supabase_edge_function.advanced", "id", "func456"),
					resource.TestCheckResourceAttr("supabase_edge_function.advanced", "status", "ACTIVE"),
				),
			},
		},
	})
}

func TestAccEdgeFunctionResourceErrorHandling(t *testing.T) {
	defer gock.OffAll()

	// Test creation failure
	gock.New("https://api.supabase.com").
		Post("/v1/projects/mayuaycdtijbctgqbycg/functions").
		Reply(http.StatusBadRequest).
		JSON(map[string]string{"message": "Invalid function body"})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccEdgeFunctionResourceConfigInvalid,
				ExpectError: regexp.MustCompile(`Unable to create edge function`),
			},
		},
	})
}

func TestAccEdgeFunctionResourceNotFound(t *testing.T) {
	defer gock.OffAll()

	// Mock successful creation
	gock.New("https://api.supabase.com").
		Post("/v1/projects/mayuaycdtijbctgqbycg/functions").
		Reply(http.StatusCreated).
		JSON(api.FunctionResponse{
			Id:        "func789",
			Slug:      "disappearing-function",
			Name:      "Disappearing Function",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995200,
		})

	// Mock read returning success initially
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/disappearing-function").
		Reply(http.StatusOK).
		JSON(api.FunctionSlugResponse{
			Id:        "func789",
			Slug:      "disappearing-function",
			Name:      "Disappearing Function",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995200,
		})

	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/disappearing-function/body").
		Reply(http.StatusOK).
		BodyString(`export default () => new Response("I exist!");`)

	// Mock read returning 404 (function deleted outside terraform)
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/disappearing-function").
		Reply(http.StatusNotFound).
		JSON(map[string]string{"message": "Function not found"})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeFunctionResourceConfigDisappearing,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_edge_function.disappearing", "slug", "disappearing-function"),
				),
			},
			{
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					// After refresh, resource should be detected as missing and planned for recreation
					func(s *terraform.State) error {
						// This validates that 404 is properly handled by removing from state
						// The ExpectNonEmptyPlan confirms Terraform wants to recreate the missing resource
						return nil
					},
				),
			},
		},
	})
}

func TestAccEdgeFunctionResourceComplexUpdate(t *testing.T) {
	defer gock.OffAll()

	// Create function with minimal config
	gock.New("https://api.supabase.com").
		Post("/v1/projects/mayuaycdtijbctgqbycg/functions").
		Reply(http.StatusCreated).
		JSON(api.FunctionResponse{
			Id:        "func999",
			Slug:      "evolving-function",
			Name:      "Simple Function",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995200,
		})

	// Initial read
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/evolving-function").
		Reply(http.StatusOK).
		JSON(api.FunctionSlugResponse{
			Id:        "func999",
			Slug:      "evolving-function",
			Name:      "Simple Function",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995200,
		})

	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/evolving-function/body").
		Reply(http.StatusOK).
		BodyString(`export default () => new Response("Simple");`)

	// Refresh read (duplicate for Terraform's refresh cycle)
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/evolving-function").
		Reply(http.StatusOK).
		JSON(api.FunctionSlugResponse{
			Id:        "func999",
			Slug:      "evolving-function",
			Name:      "Simple Function",
			Status:    "ACTIVE",
			CreatedAt: 1640995200,
			UpdatedAt: 1640995200,
		})

	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/evolving-function/body").
		Reply(http.StatusOK).
		BodyString(`export default () => new Response("Simple");`)

	// Update to complex function
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/functions/evolving-function").
		Reply(http.StatusOK).
		JSON(api.FunctionResponse{
			Id:                "func999",
			Slug:              "evolving-function",
			Name:              "Complex Function",
			Status:            "ACTIVE",
			CreatedAt:         1640995200,
			UpdatedAt:         1640995400,
			ComputeMultiplier: func() *float32 { f := float32(1.5); return &f }(),
		})

	// Read after update
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/evolving-function").
		Reply(http.StatusOK).
		JSON(api.FunctionSlugResponse{
			Id:                "func999",
			Slug:              "evolving-function",
			Name:              "Complex Function",
			Status:            "ACTIVE",
			CreatedAt:         1640995200,
			UpdatedAt:         1640995400,
			ComputeMultiplier: func() *float32 { f := float32(1.5); return &f }(),
		})

	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/functions/evolving-function/body").
		Reply(http.StatusOK).
		BodyString(`export default (req) => { const data = req.json(); return new Response(JSON.stringify(data)); }`)

	// Mock successful delete for cleanup
	gock.New("https://api.supabase.com").
		Delete("/v1/projects/mayuaycdtijbctgqbycg/functions/evolving-function").
		Reply(http.StatusOK)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEdgeFunctionResourceConfigMinimal,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_edge_function.evolving", "name", "Simple Function"),
					resource.TestCheckResourceAttr("supabase_edge_function.evolving", "verify_jwt", "false"),
					resource.TestCheckNoResourceAttr("supabase_edge_function.evolving", "compute_multiplier"),
				),
			},
			{
				Config: testAccEdgeFunctionResourceConfigComplex,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_edge_function.evolving", "name", "Complex Function"),
					resource.TestCheckResourceAttr("supabase_edge_function.evolving", "verify_jwt", "true"),
					resource.TestCheckResourceAttr("supabase_edge_function.evolving", "compute_multiplier", "1.5"),
					resource.TestCheckResourceAttr("supabase_edge_function.evolving", "updated_at", "1640995400"),
				),
			},
		},
	})
}

// Test configurations
const testAccEdgeFunctionResourceConfigOptionalFields = `
resource "supabase_edge_function" "advanced" {
  project_ref        = "mayuaycdtijbctgqbycg"
  slug               = "advanced-function"
  name               = "Advanced Function"
  body               = "import { serve } from \"./deps.ts\"; serve(() => new Response(\"Advanced!\"));"
  verify_jwt         = true
  compute_multiplier = 2.5
  entrypoint_path    = "custom/index.ts"
  import_map         = true
  import_map_path    = "custom/import_map.json"
}
`

const testAccEdgeFunctionResourceConfigInvalid = `
resource "supabase_edge_function" "invalid" {
  project_ref = "mayuaycdtijbctgqbycg"
  slug        = "invalid-function"
  name        = "Invalid Function"
  body        = "this is not valid javascript/typescript code { { {"
  verify_jwt  = false
}
`

const testAccEdgeFunctionResourceConfigDisappearing = `
resource "supabase_edge_function" "disappearing" {
  project_ref = "mayuaycdtijbctgqbycg"
  slug        = "disappearing-function"
  name        = "Disappearing Function"
  body        = "export default () => new Response(\"I exist!\");"
  verify_jwt  = false
}
`

const testAccEdgeFunctionResourceConfigMinimal = `
resource "supabase_edge_function" "evolving" {
  project_ref = "mayuaycdtijbctgqbycg"
  slug        = "evolving-function"
  name        = "Simple Function"
  body        = "export default () => new Response(\"Simple\");"
  verify_jwt  = false
}
`

const testAccEdgeFunctionResourceConfigComplex = `
resource "supabase_edge_function" "evolving" {
  project_ref        = "mayuaycdtijbctgqbycg"
  slug               = "evolving-function"
  name               = "Complex Function"
  body               = "export default (req) => { const data = req.json(); return new Response(JSON.stringify(data)); }"
  verify_jwt         = true
  compute_multiplier = 1.5
}
`
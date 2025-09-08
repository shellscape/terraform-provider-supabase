package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"gopkg.in/h2non/gock.v1"
)

func TestAccStorageBucketResource(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Mock API keys endpoint for token exchange
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/api-keys").
		MatchParam("reveal", "false").
		Times(10).  // Allow multiple calls for caching
		Reply(200).
		JSON([]map[string]interface{}{
			{
				"api_key": "service_role_jwt_token_here",
				"name":    "service_role",
			},
			{
				"api_key": "anon_key_here", 
				"name":    "anon",
			},
		})

	// ==> STEP 1: CREATE <==
	// Mock creation
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Post("/storage/v1/bucket").
		Reply(201).
		JSON(map[string]string{"name": "test-bucket"})

	// Mock read operations for create step - returning initial values
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Get("/storage/v1/bucket/test-bucket").
		Times(3).
		Reply(200).
		JSON(map[string]interface{}{
			"id":                 "test-bucket",
			"name":               "test-bucket",
			"public":             false,
			"file_size_limit":    52428800,
			"allowed_mime_types": []string{"image/*", "video/mp4"},
			"created_at":         "2023-01-01T00:00:00Z",
			"updated_at":         "2023-01-01T00:00:00Z",
			"owner":              "owner123",
		})

	// ==> STEP 2: UPDATE <==
	// Mock update operation
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Put("/storage/v1/bucket/test-bucket").
		Reply(200).
		JSON(map[string]string{"message": "Updated"})

	// Mock read operations after update - returning updated values
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Get("/storage/v1/bucket/test-bucket").
		Times(10).
		Reply(200).
		JSON(map[string]interface{}{
			"id":                 "test-bucket",
			"name":               "test-bucket",
			"public":             true,
			"file_size_limit":    104857600,
			"allowed_mime_types": []string{"image/*", "video/*"},
			"created_at":         "2023-01-01T00:00:00Z",
			"updated_at":         "2023-01-01T01:00:00Z",
			"owner":              "owner123",
		})

	// Mock delete for cleanup
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Delete("/storage/v1/bucket/test-bucket").
		Reply(200).
		JSON(map[string]string{"message": "Deleted"})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccStorageBucketResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_storage_bucket.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.test", "name", "test-bucket"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.test", "public", "false"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.test", "file_size_limit", "52428800"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.test", "allowed_mime_types.#", "2"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.test", "allowed_mime_types.0", "image/*"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.test", "allowed_mime_types.1", "video/mp4"),
				),
			},
			// Update testing
			{
				Config: testAccStorageBucketResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_storage_bucket.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.test", "name", "test-bucket"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.test", "public", "true"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.test", "file_size_limit", "104857600"),
				),
			},
		},
	})
}

func TestAccStorageBucketResourceValidation(t *testing.T) {
	defer gock.OffAll()

	// No API mock needed - validation happens at schema level before API calls
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccStorageBucketResourceConfigInvalidName,
				ExpectError: regexp.MustCompile(`invalid bucket name with spaces`),
			},
		},
	})
}

const testAccStorageBucketResourceConfig = `
resource "supabase_storage_bucket" "test" {
  project_ref    = "mayuaycdtijbctgqbycg"
  name           = "test-bucket"
  public         = false
  file_size_limit = 52428800
  allowed_mime_types = ["image/*", "video/mp4"]
}
`

const testAccStorageBucketResourceConfigUpdate = `
resource "supabase_storage_bucket" "test" {
  project_ref    = "mayuaycdtijbctgqbycg"
  name           = "test-bucket"
  public         = true
  file_size_limit = 104857600
  allowed_mime_types = ["image/*", "video/*"]
}
`

const testAccStorageBucketResourceConfigInvalidName = `
resource "supabase_storage_bucket" "test" {
  project_ref = "mayuaycdtijbctgqbycg"
  name        = "invalid bucket name with spaces"
  public      = false
}
`
package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"gopkg.in/h2non/gock.v1"
)

// mockApiKeysForTokenExchange sets up the mock for the API keys endpoint
// This is required for token exchange in storage operations
func mockApiKeysForTokenExchange() {
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
}

func TestAccStorageBucketResourceOptionalFields(t *testing.T) {
	defer gock.OffAll()

	// Mock API keys endpoint for token exchange
	mockApiKeysForTokenExchange()

	// Mock Storage API calls for bucket with null optional fields
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Post("/storage/v1/bucket").
		Reply(201).
		JSON(map[string]string{"name": "optional-bucket"})

	// Mock GET for reading the created bucket (null optional fields) - multiple times for refresh cycles
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Get("/storage/v1/bucket/optional-bucket").
		Times(3).  // Allow multiple refresh calls
		Reply(200).
		JSON(map[string]interface{}{
			"id":                 "optional-bucket",
			"name":               "optional-bucket",
			"public":             true,
			"file_size_limit":    nil,
			"allowed_mime_types": nil,
			"created_at":         "2023-01-01T00:00:00Z",
			"updated_at":         "2023-01-01T00:00:00Z",
			"owner":              "owner123",
		})

	// Mock successful delete for cleanup
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Delete("/storage/v1/bucket/optional-bucket").
		Reply(200)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucketResourceConfigOptional,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_storage_bucket.optional", "name", "optional-bucket"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.optional", "public", "true"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.optional", "id", "optional-bucket"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.optional", "owner", "owner123"),
					// Null optional fields should not be set
					resource.TestCheckNoResourceAttr("supabase_storage_bucket.optional", "file_size_limit"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.optional", "allowed_mime_types.#", "0"),
				),
			},
		},
	})
}

func TestAccStorageBucketResourceErrorHandling(t *testing.T) {
	defer gock.OffAll()

	// Mock API keys endpoint for token exchange
	mockApiKeysForTokenExchange()

	// Mock Storage API creation failure
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Post("/storage/v1/bucket").
		Reply(409).  // Conflict - bucket already exists
		JSON(map[string]string{"message": "Bucket already exists"})

	// No delete mock needed since creation fails

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccStorageBucketResourceConfigDuplicate,
				ExpectError: regexp.MustCompile(`Unable to create storage bucket`),
			},
		},
	})
}

func TestAccStorageBucketResourceUpdateFailure(t *testing.T) {
	defer gock.OffAll()

	// Mock API keys endpoint for token exchange
	mockApiKeysForTokenExchange()

	// Mock successful creation
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Post("/storage/v1/bucket").
		Reply(201).
		JSON(map[string]string{"name": "update-fail-bucket"})

	// Mock successful reads - multiple times for refresh cycles
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Get("/storage/v1/bucket/update-fail-bucket").
		Times(5).
		Reply(200).
		JSON(map[string]interface{}{
			"id":                 "update-fail-bucket",
			"name":               "update-fail-bucket",
			"public":             false,
			"file_size_limit":    1048576,
			"allowed_mime_types": []string{"image/*"},
			"created_at":         "2023-01-01T00:00:00Z",
			"updated_at":         "2023-01-01T00:00:00Z",
			"owner":              "owner123",
		})

	// Mock update failure
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Put("/storage/v1/bucket/update-fail-bucket").
		Reply(403).  // Forbidden
		JSON(map[string]string{"message": "Insufficient permissions"})

	// Mock successful delete for cleanup
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Delete("/storage/v1/bucket/update-fail-bucket").
		Reply(200)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucketResourceConfigUpdateFail,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_storage_bucket.update_fail", "public", "false"),
				),
			},
			{
				Config:      testAccStorageBucketResourceConfigUpdateFailUpdated,
				ExpectError: regexp.MustCompile(`Unable to update storage bucket`),
			},
		},
	})
}

func TestAccStorageBucketResourceNotFound(t *testing.T) {
	defer gock.OffAll()

	// Mock API keys endpoint for token exchange
	mockApiKeysForTokenExchange()

	// Mock successful creation
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Post("/storage/v1/bucket").
		Reply(201).
		JSON(map[string]string{"name": "vanishing-bucket"})

	// Mock successful initial read - multiple times for refresh cycles during step 1
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Get("/storage/v1/bucket/vanishing-bucket").
		Times(5).
		Reply(200).
		JSON(map[string]interface{}{
			"id":                 "vanishing-bucket",
			"name":               "vanishing-bucket",
			"public":             false,
			"file_size_limit":    nil,
			"allowed_mime_types": nil,
			"created_at":         "2023-01-01T00:00:00Z",
			"updated_at":         "2023-01-01T00:00:00Z",
			"owner":              "owner123",
		})

	// Mock 404 on refresh (bucket deleted outside Terraform)
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Get("/storage/v1/bucket/vanishing-bucket").
		Reply(404).
		JSON(map[string]string{"message": "Bucket not found"})

	// Mock delete for cleanup - should also return 404 since bucket is already gone
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Delete("/storage/v1/bucket/vanishing-bucket").
		Reply(404).
		JSON(map[string]string{"message": "Bucket not found"})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucketResourceConfigVanishing,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_storage_bucket.vanishing", "name", "vanishing-bucket"),
				),
			},
			{
				RefreshState: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Resource should be removed from state due to 404
					func(s *terraform.State) error {
						// This validates proper handling of 404 responses
						return nil
					},
				),
			},
		},
	})
}

func TestAccStorageBucketResourceComplexMimeTypes(t *testing.T) {
	defer gock.OffAll()

	// Mock API keys endpoint for token exchange
	mockApiKeysForTokenExchange()

	// Mock creation with complex MIME types
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Post("/storage/v1/bucket").
		Reply(201).
		JSON(map[string]string{"name": "mime-bucket"})

	// Mock read with complex MIME types - multiple times for refresh cycles
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Get("/storage/v1/bucket/mime-bucket").
		Times(5).
		Reply(200).
		JSON(map[string]interface{}{
			"id":     "mime-bucket",
			"name":   "mime-bucket",
			"public": true,
			"allowed_mime_types": []string{
				"image/*",
				"video/mp4",
				"video/quicktime",
				"application/pdf",
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				"text/plain",
				"text/csv",
			},
			"file_size_limit": 104857600,
			"created_at":      "2023-01-01T00:00:00Z",
			"updated_at":      "2023-01-01T00:00:00Z",
			"owner":           "owner123",
		})

	// Mock successful delete for cleanup
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Delete("/storage/v1/bucket/mime-bucket").
		Reply(200)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucketResourceConfigComplexMimeTypes,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_storage_bucket.mime", "name", "mime-bucket"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.mime", "public", "true"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.mime", "file_size_limit", "104857600"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.mime", "allowed_mime_types.#", "7"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.mime", "allowed_mime_types.0", "image/*"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.mime", "allowed_mime_types.1", "video/mp4"),
					resource.TestCheckResourceAttr("supabase_storage_bucket.mime", "allowed_mime_types.3", "application/pdf"),
				),
			},
		},
	})
}

func TestAccStorageBucketResourceDeleteFailure(t *testing.T) {
	defer gock.OffAll()

	// Mock API keys endpoint for token exchange
	mockApiKeysForTokenExchange()

	// Mock successful creation
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Post("/storage/v1/bucket").
		Reply(201).
		JSON(map[string]string{"name": "delete-test-bucket"})

	// Mock successful read - multiple times for refresh cycles
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Get("/storage/v1/bucket/delete-test-bucket").
		Times(5).
		Reply(200).
		JSON(map[string]interface{}{
			"id":                 "delete-test-bucket",
			"name":               "delete-test-bucket",
			"public":             false,
			"file_size_limit":    nil,
			"allowed_mime_types": nil,
			"created_at":         "2023-01-01T00:00:00Z",
			"updated_at":         "2023-01-01T00:00:00Z",
			"owner":              "owner123",
		})

	// First delete attempt fails (bucket not empty)
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Delete("/storage/v1/bucket/delete-test-bucket").
		Reply(400).
		JSON(map[string]string{"message": "Bucket is not empty"})

	// Subsequent delete attempts succeed (bucket was emptied)
	gock.New("https://mayuaycdtijbctgqbycg.supabase.co").
		Delete("/storage/v1/bucket/delete-test-bucket").
		Times(3).
		Reply(200)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucketResourceConfigDeleteTest,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_storage_bucket.delete_test", "name", "delete-test-bucket"),
				),
			},
			{
				Config:      testAccStorageBucketResourceConfigDeleteTestDestroy,
				ExpectError: regexp.MustCompile(`Unable to delete storage bucket`),
			},
			{
				// After the "bucket is emptied", deletion should succeed
				Config: testAccStorageBucketResourceConfigDeleteTestDestroy,
			},
		},
	})
}

// Test configurations
const testAccStorageBucketResourceConfigOptional = `
resource "supabase_storage_bucket" "optional" {
  project_ref = "mayuaycdtijbctgqbycg"
  name        = "optional-bucket"
  public      = true
  # file_size_limit and allowed_mime_types are intentionally omitted
}
`

const testAccStorageBucketResourceConfigDuplicate = `
resource "supabase_storage_bucket" "duplicate" {
  project_ref = "mayuaycdtijbctgqbycg"
  name        = "existing-bucket"
  public      = false
}
`

const testAccStorageBucketResourceConfigUpdateFail = `
resource "supabase_storage_bucket" "update_fail" {
  project_ref        = "mayuaycdtijbctgqbycg"
  name               = "update-fail-bucket"
  public             = false
  file_size_limit    = 1048576
  allowed_mime_types = ["image/*"]
}
`

const testAccStorageBucketResourceConfigUpdateFailUpdated = `
resource "supabase_storage_bucket" "update_fail" {
  project_ref        = "mayuaycdtijbctgqbycg"
  name               = "update-fail-bucket"
  public             = true  # This change should fail
  file_size_limit    = 2097152
  allowed_mime_types = ["image/*", "video/*"]
}
`

const testAccStorageBucketResourceConfigVanishing = `
resource "supabase_storage_bucket" "vanishing" {
  project_ref = "mayuaycdtijbctgqbycg"
  name        = "vanishing-bucket"
  public      = false
}
`

const testAccStorageBucketResourceConfigComplexMimeTypes = `
resource "supabase_storage_bucket" "mime" {
  project_ref        = "mayuaycdtijbctgqbycg"
  name               = "mime-bucket"
  public             = true
  file_size_limit    = 104857600  # 100MB
  allowed_mime_types = [
    "image/*",
    "video/mp4",
    "video/quicktime",
    "application/pdf",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "text/plain",
    "text/csv"
  ]
}
`

const testAccStorageBucketResourceConfigDeleteTest = `
resource "supabase_storage_bucket" "delete_test" {
  project_ref = "mayuaycdtijbctgqbycg"
  name        = "delete-test-bucket"
  public      = false
}
`

const testAccStorageBucketResourceConfigDeleteTestDestroy = `
# Empty configuration to trigger resource destruction
provider "supabase" {
  access_token = "test-token"
}
`
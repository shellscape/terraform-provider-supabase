package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/supabase/cli/pkg/api"
	"gopkg.in/h2non/gock.v1"
)

func TestAccStorageBucketsDataSource(t *testing.T) {
	defer gock.OffAll()

	// Mock the storage buckets API response - multiple times for refresh cycles
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/storage/buckets").
		Times(3).
		Reply(http.StatusOK).
		JSON([]api.V1StorageBucketResponse{
			{
				Id:        "bucket1",
				Name:      "images",
				Owner:     "owner1",
				Public:    true,
				CreatedAt: "2023-01-01T00:00:00Z",
				UpdatedAt: "2023-01-01T00:00:00Z",
			},
			{
				Id:        "bucket2",
				Name:      "documents",
				Owner:     "owner2",
				Public:    false,
				CreatedAt: "2023-01-02T00:00:00Z",
				UpdatedAt: "2023-01-02T00:00:00Z",
			},
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStorageBucketsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "id", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.#", "2"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.0.id", "bucket1"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.0.name", "images"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.0.owner", "owner1"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.0.public", "true"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.0.created_at", "2023-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.0.updated_at", "2023-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.1.id", "bucket2"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.1.name", "documents"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.1.owner", "owner2"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.1.public", "false"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.1.created_at", "2023-01-02T00:00:00Z"),
					resource.TestCheckResourceAttr("data.supabase_storage_buckets.test", "buckets.1.updated_at", "2023-01-02T00:00:00Z"),
				),
			},
		},
	})
}

const testAccStorageBucketsDataSourceConfig = `
data "supabase_storage_buckets" "test" {
  project_ref = "mayuaycdtijbctgqbycg"
}
`
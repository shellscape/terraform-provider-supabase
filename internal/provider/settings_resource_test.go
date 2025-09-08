// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/supabase/cli/pkg/api"
	"gopkg.in/h2non/gock.v1"
)

func TestAccSettingsResource(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// ==> INITIAL CREATE OPERATIONS (Step 1) <==
	// Initial PUT operations for create
	gock.New("https://api.supabase.com").
		Put("/v1/projects/mayuaycdtijbctgqbycg/config/database/postgres").
		Reply(http.StatusOK).
		JSON(api.PostgresConfigResponse{
			StatementTimeout: Ptr("10s"),
		})
	gock.New("https://api.supabase.com").
		Post("/v1/projects/mayuaycdtijbctgqbycg/network-restrictions/apply").
		Reply(http.StatusCreated).
		JSON(api.NetworkRestrictionsResponse{
			Config: api.NetworkRestrictionsRequest{
				DbAllowedCidrs:   Ptr([]string{"0.0.0.0/0"}),
				DbAllowedCidrsV6: Ptr([]string{"::/0"}),
			},
		})
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/postgrest").
		Reply(http.StatusOK).
		JSON(api.V1PostgrestConfigResponse{
			DbExtraSearchPath: "public,extensions",
			DbSchema:          "public,storage,graphql_public",
			MaxRows:           1000,
		})
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			SiteUrl:           Ptr("http://localhost:3000"),
			UriAllowList:      Ptr("http://localhost:3000/auth/callback"),
			MailerOtpExp:      3600,
			MfaPhoneOtpLength: 6,
			SmsOtpLength:      6,
		})
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/storage").
		Reply(http.StatusOK).
		JSON(api.StorageConfigResponse{
			FileSizeLimit: 52428800,
			Features: api.StorageFeatures{
				ImageTransformation: api.StorageFeatureImageTransformation{
					Enabled: true,
				},
			},
		})
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/database/pooler").
		Reply(http.StatusOK).
		JSON(api.UpdateSupavisorConfigResponse{
			DefaultPoolSize: Ptr(20),
			PoolMode:        api.UpdateSupavisorConfigResponsePoolModeSession,
		})

	// Read operations during step 1 - exactly 2 rounds of reads (8 total)
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/database/postgres").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.PostgresConfigResponse{
			StatementTimeout: Ptr("10s"),
		})
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/network-restrictions").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.NetworkRestrictionsResponse{
			Config: api.NetworkRestrictionsRequest{
				DbAllowedCidrs:   Ptr([]string{"0.0.0.0/0"}),
				DbAllowedCidrsV6: Ptr([]string{"::/0"}),
			},
		})
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/postgrest").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.V1PostgrestConfigResponse{
			DbExtraSearchPath: "public,extensions",
			DbSchema:          "public,storage,graphql_public",
			MaxRows:           1000,
		})
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			SiteUrl:           Ptr("http://localhost:3000"),
			UriAllowList:      Ptr("http://localhost:3000/auth/callback"),
			MailerOtpExp:      3600,
			MfaPhoneOtpLength: 6,
			SmsOtpLength:      6,
		})
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/storage").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.StorageConfigResponse{
			FileSizeLimit: 52428800,
			Features: api.StorageFeatures{
				ImageTransformation: api.StorageFeatureImageTransformation{
					Enabled: true,
				},
			},
		})
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/database/pooler").
		Times(2).
		Reply(http.StatusOK).
		JSON([]api.SupavisorConfigResponse{
			{
				DefaultPoolSize:  Ptr(20),
				PoolMode:         api.SupavisorConfigResponsePoolModeSession,
				ConnectionString: "postgresql://postgres:password@db.mayuaycdtijbctgqbycg.supabase.co:5432/postgres",
				DatabaseType:     api.PRIMARY,
				DbHost:          "db.mayuaycdtijbctgqbycg.supabase.co",
				DbName:          "postgres",
				DbPort:          5432,
				DbUser:          "postgres",
				Identifier:      "primary",
				IsUsingScramAuth: true,
				MaxClientConn:   Ptr(200),
			},
		})

	// ==> UPDATE OPERATIONS (Step 2) <==
	// Update operations
	gock.New("https://api.supabase.com").
		Put("/v1/projects/mayuaycdtijbctgqbycg/config/database/postgres").
		Reply(http.StatusOK).
		JSON(api.PostgresConfigResponse{
			StatementTimeout: Ptr("20s"),
			MaxConnections:   Ptr(200),
		})
	gock.New("https://api.supabase.com").
		Post("/v1/projects/mayuaycdtijbctgqbycg/network-restrictions/apply").
		Reply(http.StatusCreated).
		JSON(api.NetworkRestrictionsResponse{
			Config: api.NetworkRestrictionsRequest{
				DbAllowedCidrs:   Ptr([]string{"8.8.8.0/24"}),
				DbAllowedCidrsV6: nil,
			},
		})
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/postgrest").
		Reply(http.StatusOK).
		JSON(api.V1PostgrestConfigResponse{
			DbExtraSearchPath: "public,extensions",
			DbSchema:          "public,storage,graphql_public",
			MaxRows:           2000,
		})
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			SiteUrl:           Ptr("http://localhost:3001"),
			UriAllowList:      Ptr("http://localhost:3001/auth/callback,https://app.example.com/callback"),
			MailerOtpExp:      7200,
			MfaPhoneOtpLength: 8,
			SmsOtpLength:      8,
		})
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/storage").
		Reply(http.StatusOK).
		JSON(api.StorageConfigResponse{
			FileSizeLimit: 104857600,
			Features: api.StorageFeatures{
				ImageTransformation: api.StorageFeatureImageTransformation{
					Enabled: false,
				},
			},
		})
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/database/pooler").
		Reply(http.StatusOK).
		JSON(api.UpdateSupavisorConfigResponse{
			DefaultPoolSize: Ptr(40),
			PoolMode:        api.UpdateSupavisorConfigResponsePoolModeSession,
		})

	// Read operations after update and during step 2 refresh - exactly 1 round (4 total)
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/database/postgres").
		Reply(http.StatusOK).
		JSON(api.PostgresConfigResponse{
			StatementTimeout: Ptr("20s"),
			MaxConnections:   Ptr(200),
		})
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/network-restrictions").
		Reply(http.StatusOK).
		JSON(api.NetworkRestrictionsResponse{
			Config: api.NetworkRestrictionsRequest{
				DbAllowedCidrs:   Ptr([]string{"8.8.8.0/24"}),
				DbAllowedCidrsV6: nil,
			},
		})
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/postgrest").
		Reply(http.StatusOK).
		JSON(api.V1PostgrestConfigResponse{
			DbExtraSearchPath: "public,extensions",
			DbSchema:          "public,storage,graphql_public",
			MaxRows:           2000,
		})
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			SiteUrl:           Ptr("http://localhost:3001"),
			UriAllowList:      Ptr("http://localhost:3001/auth/callback,https://app.example.com/callback"),
			MailerOtpExp:      7200,
			MfaPhoneOtpLength: 8,
			SmsOtpLength:      8,
		})
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/storage").
		Reply(http.StatusOK).
		JSON(api.StorageConfigResponse{
			FileSizeLimit: 104857600,
			Features: api.StorageFeatures{
				ImageTransformation: api.StorageFeatureImageTransformation{
					Enabled: false,
				},
			},
		})
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/database/pooler").
		Reply(http.StatusOK).
		JSON([]api.SupavisorConfigResponse{
			{
				DefaultPoolSize:  Ptr(40),
				PoolMode:         api.SupavisorConfigResponsePoolModeSession,
				ConnectionString: "postgresql://postgres:password@db.mayuaycdtijbctgqbycg.supabase.co:5432/postgres",
				DatabaseType:     api.PRIMARY,
				DbHost:          "db.mayuaycdtijbctgqbycg.supabase.co",
				DbName:          "postgres",
				DbPort:          5432,
				DbUser:          "postgres",
				Identifier:      "primary",
				IsUsingScramAuth: true,
				MaxClientConn:   Ptr(200),
			},
		})



	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSettingsResourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "database.statement_timeout", "10s"),
					resource.TestCheckResourceAttr("supabase_settings.test", "network.db_allowed_cidrs.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("supabase_settings.test", "network.db_allowed_cidrs_v6.0", "::/0"),
					resource.TestCheckResourceAttr("supabase_settings.test", "api.db_extra_search_path", "public,extensions"),
					resource.TestCheckResourceAttr("supabase_settings.test", "api.db_schema", "public,storage,graphql_public"),
					resource.TestCheckResourceAttr("supabase_settings.test", "api.max_rows", "1000"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.site_url", "http://localhost:3000"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.uri_allow_list", "http://localhost:3000/auth/callback"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "3600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "storage.file_size_limit", "52428800"),
					resource.TestCheckResourceAttr("supabase_settings.test", "storage.features.image_transformation.enabled", "true"),
					resource.TestCheckResourceAttr("supabase_settings.test", "pooler.default_pool_size", "20"),
					resource.TestCheckResourceAttr("supabase_settings.test", "id", "mayuaycdtijbctgqbycg"),
				),
			},
			// Update testing
			{
				Config: testAccSettingsResourceConfigUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "database.statement_timeout", "20s"),
					resource.TestCheckResourceAttr("supabase_settings.test", "database.max_connections", "200"),
					resource.TestCheckResourceAttr("supabase_settings.test", "network.db_allowed_cidrs.0", "8.8.8.0/24"),
					resource.TestCheckResourceAttr("supabase_settings.test", "api.max_rows", "2000"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.site_url", "http://localhost:3001"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.uri_allow_list", "http://localhost:3001/auth/callback,https://app.example.com/callback"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "7200"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "8"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "8"),
					resource.TestCheckResourceAttr("supabase_settings.test", "storage.file_size_limit", "104857600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "storage.features.image_transformation.enabled", "false"),
					resource.TestCheckResourceAttr("supabase_settings.test", "pooler.default_pool_size", "40"),
				),
			},
		},
	})
}

func TestAccSettingsResourceValidation(t *testing.T) {
	defer gock.OffAll()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccSettingsResourceConfigInvalidCIDR,
				ExpectError: regexp.MustCompile(`Invalid CIDR`),
			},
			{
				Config:      testAccSettingsResourceConfigPrivateIP,
				ExpectError: regexp.MustCompile(`Private IP not allowed`),
			},
		},
	})
}

const testAccSettingsResourceConfig = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  database = {
    statement_timeout = "10s"
  }

  network = {
    db_allowed_cidrs = ["0.0.0.0/0"]
    db_allowed_cidrs_v6 = ["::/0"]
  }

  api = {
    db_extra_search_path = "public,extensions"
    db_schema = "public,storage,graphql_public"
    max_rows = 1000
  }

  auth = {
    site_url = "http://localhost:3000"
    uri_allow_list = "http://localhost:3000/auth/callback"
    mailer_otp_exp = 3600
    mfa_phone_otp_length = 6
    sms_otp_length = 6
  }

  storage = {
    file_size_limit = 52428800
    features = {
      image_transformation = {
        enabled = true
      }
    }
  }

  pooler = {
    default_pool_size = 20
  }
}
`

const testAccSettingsResourceConfigUpdate = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  database = {
    statement_timeout = "20s"
    max_connections = 200
  }

  network = {
    db_allowed_cidrs = ["8.8.8.0/24"]
  }

  api = {
    db_extra_search_path = "public,extensions"
    db_schema = "public,storage,graphql_public"
    max_rows = 2000
  }

  auth = {
    site_url = "http://localhost:3001"
    uri_allow_list = "http://localhost:3001/auth/callback,https://app.example.com/callback"
    mailer_otp_exp = 7200
    mfa_phone_otp_length = 8
    sms_otp_length = 8
  }

  storage = {
    file_size_limit = 104857600
    features = {
      image_transformation = {
        enabled = false
      }
    }
  }

  pooler = {
    default_pool_size = 40
  }
}
`

const testAccSettingsResourceConfigInvalidCIDR = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  network = {
    db_allowed_cidrs = ["invalid-cidr"]
  }
}
`

const testAccSettingsResourceConfigPrivateIP = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  network = {
    db_allowed_cidrs = ["10.0.0.0/8"]
  }
}
`

func TestAccSettingsResourceExternalGithub(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Create operations
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			ExternalGithubEnabled:  Ptr(true),
			ExternalGithubClientId: Ptr("github_client_id_123"),
			MailerOtpExp:           3600,
			MfaPhoneOtpLength:      6,
			SmsOtpLength:           6,
		})

	// Read operations
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			ExternalGithubEnabled:  Ptr(true),
			ExternalGithubClientId: Ptr("github_client_id_123"),
			MailerOtpExp:           3600,
			MfaPhoneOtpLength:      6,
			SmsOtpLength:           6,
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigExternalGithub,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "3600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.external_github.enabled", "true"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.external_github.client_id", "github_client_id_123"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.external_github.secret", "github_secret_456"),
				),
			},
		},
	})
}

const testAccSettingsResourceConfigExternalGithub = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  auth = {
    mailer_otp_exp = 3600
    mfa_phone_otp_length = 6
    sms_otp_length = 6
    external_github = {
      enabled = true
      client_id = "github_client_id_123"
      secret = "github_secret_456"
    }
  }
}
`

func TestAccSettingsResourceImageTransformationEnabled(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Test that users can set the enabled property for image transformation
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/storage").
		Reply(http.StatusOK).
		JSON(api.StorageConfigResponse{
			FileSizeLimit: 52428800,
			Features: api.StorageFeatures{
				ImageTransformation: api.StorageFeatureImageTransformation{
					Enabled: true,
				},
			},
		})

	// Read operations
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/storage").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.StorageConfigResponse{
			FileSizeLimit: 52428800,
			Features: api.StorageFeatures{
				ImageTransformation: api.StorageFeatureImageTransformation{
					Enabled: true,
				},
			},
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigImageTransformation,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "storage.file_size_limit", "52428800"),
					resource.TestCheckResourceAttr("supabase_settings.test", "storage.features.image_transformation.enabled", "true"),
				),
			},
		},
	})
}

const testAccSettingsResourceConfigImageTransformation = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  storage = {
    file_size_limit = 52428800
    features = {
      image_transformation = {
        enabled = true
      }
    }
  }
}
`

func TestAccSettingsResourceS3Protocol(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Test that users can set the s3Protocol enabled property
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/storage").
		Reply(http.StatusOK).
		JSON(api.StorageConfigResponse{
			FileSizeLimit: 52428800,
			Features: api.StorageFeatures{
				ImageTransformation: api.StorageFeatureImageTransformation{
					Enabled: false,
				},
			},
		})

	// Read operations
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/storage").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.StorageConfigResponse{
			FileSizeLimit: 52428800,
			Features: api.StorageFeatures{
				ImageTransformation: api.StorageFeatureImageTransformation{
					Enabled: false,
				},
			},
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigS3Protocol,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "storage.file_size_limit", "52428800"),
					resource.TestCheckResourceAttr("supabase_settings.test", "storage.features.s3_protocol.enabled", "true"),
				),
			},
		},
	})
}

const testAccSettingsResourceConfigS3Protocol = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  storage = {
    file_size_limit = 52428800
    features = {
      s3_protocol = {
        enabled = true
      }
    }
  }
}
`

func TestAccSettingsResourceExternalGithubClientIdDirect(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Test that users can set the external_github_client_id as direct property
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			ExternalGithubClientId: Ptr("direct_github_client_id_456"),
			MailerOtpExp:           3600,
			MfaPhoneOtpLength:      6,
			SmsOtpLength:           6,
		})

	// Read operations
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			ExternalGithubClientId: Ptr("direct_github_client_id_456"),
			MailerOtpExp:           3600,
			MfaPhoneOtpLength:      6,
			SmsOtpLength:           6,
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigExternalGithubClientIdDirect,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.external_github_client_id", "direct_github_client_id_456"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "3600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "6"),
				),
			},
		},
	})
}

const testAccSettingsResourceConfigExternalGithubClientIdDirect = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  auth = {
    mailer_otp_exp = 3600
    mfa_phone_otp_length = 6
    sms_otp_length = 6
    external_github_client_id = "direct_github_client_id_456"
  }
}
`

func TestAccSettingsResourceExternalGithubEnabledDirect(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Test that users can set the external_github_enabled as direct property
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			ExternalGithubEnabled: Ptr(true),
			MailerOtpExp:          3600,
			MfaPhoneOtpLength:     6,
			SmsOtpLength:          6,
		})

	// Read operations
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			ExternalGithubEnabled: Ptr(true),
			MailerOtpExp:          3600,
			MfaPhoneOtpLength:     6,
			SmsOtpLength:          6,
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigExternalGithubEnabledDirect,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.external_github_enabled", "true"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "3600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "6"),
				),
			},
		},
	})
}

const testAccSettingsResourceConfigExternalGithubEnabledDirect = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  auth = {
    mailer_otp_exp = 3600
    mfa_phone_otp_length = 6
    sms_otp_length = 6
    external_github_enabled = true
  }
}
`

func TestAccSettingsResourceExternalGoogleEnabledDirect(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Test that users can set the external_google_enabled as direct property
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			ExternalGoogleEnabled: Ptr(true),
			MailerOtpExp:          3600,
			MfaPhoneOtpLength:     6,
			SmsOtpLength:          6,
		})

	// Read operations
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			ExternalGoogleEnabled: Ptr(true),
			MailerOtpExp:          3600,
			MfaPhoneOtpLength:     6,
			SmsOtpLength:          6,
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigExternalGoogleEnabledDirect,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.external_google_enabled", "true"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "3600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "6"),
				),
			},
		},
	})
}

const testAccSettingsResourceConfigExternalGoogleEnabledDirect = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  auth = {
    mailer_otp_exp = 3600
    mfa_phone_otp_length = 6
    sms_otp_length = 6
    external_google_enabled = true
  }
}
`

func TestAccSettingsResourceExternalGoogleClientIdDirect(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Test that users can set the external_google_client_id as direct property
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			ExternalGoogleClientId: Ptr("direct_google_client_id_123"),
			MailerOtpExp:           3600,
			MfaPhoneOtpLength:      6,
			SmsOtpLength:           6,
		})

	// Read operations
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			ExternalGoogleClientId: Ptr("direct_google_client_id_123"),
			MailerOtpExp:           3600,
			MfaPhoneOtpLength:      6,
			SmsOtpLength:           6,
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigExternalGoogleClientIdDirect,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.external_google_client_id", "direct_google_client_id_123"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "3600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "6"),
				),
			},
		},
	})
}

const testAccSettingsResourceConfigExternalGoogleClientIdDirect = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  auth = {
    mailer_otp_exp = 3600
    mfa_phone_otp_length = 6
    sms_otp_length = 6
    external_google_client_id = "direct_google_client_id_123"
  }
}
`

func TestAccSettingsResourceMailerAutoconfirm(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Test that users can set mailer_autoconfirm
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			MailerAutoconfirm: Ptr(true),
			MailerOtpExp:      3600,
			MfaPhoneOtpLength: 6,
			SmsOtpLength:      6,
		})

	// Read operations
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			MailerAutoconfirm: Ptr(true),
			MailerOtpExp:      3600,
			MfaPhoneOtpLength: 6,
			SmsOtpLength:      6,
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigMailerAutoconfirm,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_autoconfirm", "true"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "3600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "6"),
				),
			},
		},
	})
}

const testAccSettingsResourceConfigMailerAutoconfirm = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  auth = {
    mailer_otp_exp = 3600
    mfa_phone_otp_length = 6
    sms_otp_length = 6
    mailer_autoconfirm = true
  }
}
`

func TestAccSettingsResourcePasswordMinLength(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Test that users can set password_min_length
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			PasswordMinLength: Ptr(8),
			MailerOtpExp:      3600,
			MfaPhoneOtpLength: 6,
			SmsOtpLength:      6,
		})

	// Read operations
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			PasswordMinLength: Ptr(8),
			MailerOtpExp:      3600,
			MfaPhoneOtpLength: 6,
			SmsOtpLength:      6,
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigPasswordMinLength,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.password_min_length", "8"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "3600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "6"),
				),
			},
		},
	})
}

const testAccSettingsResourceConfigPasswordMinLength = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  auth = {
    mailer_otp_exp = 3600
    mfa_phone_otp_length = 6
    sms_otp_length = 6
    password_min_length = 8
  }
}
`

func TestAccSettingsResourceHookSendEmailEnabled(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Test that users can set hook_send_email_enabled
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			HookSendEmailEnabled: Ptr(true),
			MailerOtpExp:         3600,
			MfaPhoneOtpLength:    6,
			SmsOtpLength:         6,
		})

	// Read operations
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			HookSendEmailEnabled: Ptr(true),
			MailerOtpExp:         3600,
			MfaPhoneOtpLength:    6,
			SmsOtpLength:         6,
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigHookSendEmailEnabled,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.hook_send_email_enabled", "true"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "3600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "6"),
				),
			},
		},
	})
}

const testAccSettingsResourceConfigHookSendEmailEnabled = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  auth = {
    mailer_otp_exp = 3600
    mfa_phone_otp_length = 6
    sms_otp_length = 6
    hook_send_email_enabled = true
  }
}
`

func TestAccSettingsResourceJwtExp(t *testing.T) {
	defer gock.OffAll()
	gock.Observe(gock.DumpRequest)

	// Test that users can set jwt_exp
	gock.New("https://api.supabase.com").
		Patch("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			JwtExp:            Ptr(7200),
			MailerOtpExp:      3600,
			MfaPhoneOtpLength: 6,
			SmsOtpLength:      6,
		})

	// Read operations
	gock.New("https://api.supabase.com").
		Get("/v1/projects/mayuaycdtijbctgqbycg/config/auth").
		Times(2).
		Reply(http.StatusOK).
		JSON(api.AuthConfigResponse{
			JwtExp:            Ptr(7200),
			MailerOtpExp:      3600,
			MfaPhoneOtpLength: 6,
			SmsOtpLength:      6,
		})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingsResourceConfigJwtExp,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("supabase_settings.test", "project_ref", "mayuaycdtijbctgqbycg"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.jwt_exp", "7200"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "3600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "6"),
				),
			},
		},
	})
}

const testAccSettingsResourceConfigJwtExp = `
resource "supabase_settings" "test" {
  project_ref = "mayuaycdtijbctgqbycg"

  auth = {
    mailer_otp_exp = 3600
    mfa_phone_otp_length = 6
    sms_otp_length = 6
    jwt_exp = 7200
  }
}
`

// Note: Ptr function is available from utils.go
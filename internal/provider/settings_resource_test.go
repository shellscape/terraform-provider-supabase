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
			MailerOtpExp:      3600,
			MfaPhoneOtpLength: 6,
			SmsOtpLength:      6,
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
			MailerOtpExp:      3600,
			MfaPhoneOtpLength: 6,
			SmsOtpLength:      6,
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
			MailerOtpExp:      7200,
			MfaPhoneOtpLength: 8,
			SmsOtpLength:      8,
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
			MailerOtpExp:      7200,
			MfaPhoneOtpLength: 8,
			SmsOtpLength:      8,
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
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "3600"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "6"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "6"),
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
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mailer_otp_exp", "7200"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.mfa_phone_otp_length", "8"),
					resource.TestCheckResourceAttr("supabase_settings.test", "auth.sms_otp_length", "8"),
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
    mailer_otp_exp = 3600
    mfa_phone_otp_length = 6
    sms_otp_length = 6
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
    mailer_otp_exp = 7200
    mfa_phone_otp_length = 8
    sms_otp_length = 8
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

// Note: Ptr function is available from utils.go
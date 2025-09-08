// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package settings

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthLocalConfig represents local authentication settings (email/password, JWT)
type AuthLocalConfig struct {
	// Basic settings
	ApiMaxRequestDuration types.Int64 `tfsdk:"api_max_request_duration"`
	DbMaxPoolSize         types.Int64 `tfsdk:"db_max_pool_size"`
	DisableSignup         types.Bool  `tfsdk:"disable_signup"`

	// Password settings
	PasswordMinLength          types.Int64  `tfsdk:"password_min_length"`
	PasswordHibpEnabled        types.Bool   `tfsdk:"password_hibp_enabled"`
	PasswordRequiredCharacters types.String `tfsdk:"password_required_characters"`

	// JWT settings
	JwtExp types.Int64 `tfsdk:"jwt_exp"`

	// Site settings
	SiteUrl types.String `tfsdk:"site_url"`

	// URI settings
	UriAllowList types.String `tfsdk:"uri_allow_list"`
}

func GetAuthLocalSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"api_max_request_duration": schema.Int64Attribute{
			MarkdownDescription: "Maximum request duration in seconds",
			Optional:            true,
		},
		"db_max_pool_size": schema.Int64Attribute{
			MarkdownDescription: "Maximum database connection pool size",
			Optional:            true,
		},
		"disable_signup": schema.BoolAttribute{
			MarkdownDescription: "Disable new user signups",
			Optional:            true,
		},
		"password_min_length": schema.Int64Attribute{
			MarkdownDescription: "Minimum password length",
			Optional:            true,
		},
		"password_hibp_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable Have I Been Pwned password validation",
			Optional:            true,
		},
		"password_required_characters": schema.StringAttribute{
			MarkdownDescription: "Required character types in passwords (e.g., lower, upper, number, special)",
			Optional:            true,
		},
		"jwt_exp": schema.Int64Attribute{
			MarkdownDescription: "JWT token expiration time in seconds",
			Optional:            true,
		},
		"site_url": schema.StringAttribute{
			MarkdownDescription: "Site URL for redirects and email links",
			Optional:            true,
		},
		"uri_allow_list": schema.StringAttribute{
			MarkdownDescription: "Comma-separated list of allowed redirect URIs",
			Optional:            true,
		},
	}
}
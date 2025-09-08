// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package settings

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthSecurityConfig represents security and rate limiting settings
type AuthSecurityConfig struct {
	// Security settings
	SecurityCaptchaEnabled  types.Bool   `tfsdk:"security_captcha_enabled"`
	SecurityCaptchaProvider types.String `tfsdk:"security_captcha_provider"`
	SecurityCaptchaSecret   types.String `tfsdk:"security_captcha_secret"`

	// Rate limiting settings
	RateLimitAnonymousUsers types.Int64 `tfsdk:"rate_limit_anonymous_users"`
	RateLimitEmailSent      types.Int64 `tfsdk:"rate_limit_email_sent"`
	RateLimitOtp            types.Int64 `tfsdk:"rate_limit_otp"`
	RateLimitSmsSent        types.Int64 `tfsdk:"rate_limit_sms_sent"`
	RateLimitTokenRefresh   types.Int64 `tfsdk:"rate_limit_token_refresh"`
	RateLimitVerify         types.Int64 `tfsdk:"rate_limit_verify"`

	// Additional security settings
	RefreshTokenRotationEnabled                   types.Bool  `tfsdk:"refresh_token_rotation_enabled"`
	SecurityManualLinkingEnabled                  types.Bool  `tfsdk:"security_manual_linking_enabled"`
	SecurityRefreshTokenReuseInterval             types.Int64 `tfsdk:"security_refresh_token_reuse_interval"`
	SecurityUpdatePasswordRequireReauthentication types.Bool  `tfsdk:"security_update_password_require_reauthentication"`

	// Session management settings
	SessionsInactivityTimeout types.Int64  `tfsdk:"sessions_inactivity_timeout"`
	SessionsSinglePerUser     types.Bool   `tfsdk:"sessions_single_per_user"`
	SessionsTags              types.String `tfsdk:"sessions_tags"`
	SessionsTimebox           types.Int64  `tfsdk:"sessions_timebox"`

	// SAML settings
	SamlAllowEncryptedAssertions types.Bool   `tfsdk:"saml_allow_encrypted_assertions"`
	SamlEnabled                  types.Bool   `tfsdk:"saml_enabled"`
	SamlExternalUrl              types.String `tfsdk:"saml_external_url"`
}

func GetAuthSecuritySchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		// Security settings
		"security_captcha_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable CAPTCHA for authentication",
			Optional:            true,
		},
		"security_captcha_provider": schema.StringAttribute{
			MarkdownDescription: "CAPTCHA provider (hcaptcha, recaptcha, turnstile)",
			Optional:            true,
		},
		"security_captcha_secret": schema.StringAttribute{
			MarkdownDescription: "CAPTCHA provider secret key",
			Optional:            true,
			Sensitive:           true,
		},

		// Rate limiting settings
		"rate_limit_anonymous_users": schema.Int64Attribute{
			MarkdownDescription: "Rate limit for anonymous users per hour",
			Optional:            true,
		},
		"rate_limit_email_sent": schema.Int64Attribute{
			MarkdownDescription: "Rate limit for emails sent per hour",
			Optional:            true,
		},
		"rate_limit_otp": schema.Int64Attribute{
			MarkdownDescription: "Rate limit for OTP requests per hour",
			Optional:            true,
		},
		"rate_limit_sms_sent": schema.Int64Attribute{
			MarkdownDescription: "Rate limit for SMS sent per hour",
			Optional:            true,
		},
		"rate_limit_token_refresh": schema.Int64Attribute{
			MarkdownDescription: "Rate limit for token refresh requests per hour",
			Optional:            true,
		},
		"rate_limit_verify": schema.Int64Attribute{
			MarkdownDescription: "Rate limit for verification requests per hour",
			Optional:            true,
		},

		// Additional security settings
		"refresh_token_rotation_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable refresh token rotation",
			Optional:            true,
		},
		"security_manual_linking_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable manual account linking",
			Optional:            true,
		},
		"security_refresh_token_reuse_interval": schema.Int64Attribute{
			MarkdownDescription: "Refresh token reuse interval in seconds",
			Optional:            true,
		},
		"security_update_password_require_reauthentication": schema.BoolAttribute{
			MarkdownDescription: "Require reauthentication for password updates",
			Optional:            true,
		},

		// Session management settings
		"sessions_inactivity_timeout": schema.Int64Attribute{
			MarkdownDescription: "Session inactivity timeout in seconds",
			Optional:            true,
		},
		"sessions_single_per_user": schema.BoolAttribute{
			MarkdownDescription: "Allow only one session per user",
			Optional:            true,
		},
		"sessions_tags": schema.StringAttribute{
			MarkdownDescription: "Session tags for categorization",
			Optional:            true,
		},
		"sessions_timebox": schema.Int64Attribute{
			MarkdownDescription: "Session timebox duration in seconds",
			Optional:            true,
		},

		// SAML settings
		"saml_allow_encrypted_assertions": schema.BoolAttribute{
			MarkdownDescription: "Allow encrypted SAML assertions",
			Optional:            true,
		},
		"saml_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable SAML authentication",
			Optional:            true,
		},
		"saml_external_url": schema.StringAttribute{
			MarkdownDescription: "External SAML URL",
			Optional:            true,
		},
	}
}
package settings

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthMfaConfig represents multi-factor authentication configuration
type AuthMfaConfig struct {
	// MFA settings
	MfaMaxEnrolledFactors    types.Int64  `tfsdk:"mfa_max_enrolled_factors"`
	MfaPhoneEnrollEnabled    types.Bool   `tfsdk:"mfa_phone_enroll_enabled"`
	MfaPhoneMaxFrequency     types.Int64  `tfsdk:"mfa_phone_max_frequency"`
	MfaPhoneOtpLength        types.Int64  `tfsdk:"mfa_phone_otp_length"`
	MfaPhoneTemplate         types.String `tfsdk:"mfa_phone_template"`
	MfaPhoneVerifyEnabled    types.Bool   `tfsdk:"mfa_phone_verify_enabled"`
	MfaTotpEnrollEnabled     types.Bool   `tfsdk:"mfa_totp_enroll_enabled"`
	MfaTotpVerifyEnabled     types.Bool   `tfsdk:"mfa_totp_verify_enabled"`
	MfaWebAuthnEnrollEnabled types.Bool   `tfsdk:"mfa_web_authn_enroll_enabled"`
	MfaWebAuthnVerifyEnabled types.Bool   `tfsdk:"mfa_web_authn_verify_enabled"`
}

func GetAuthMfaSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"mfa_max_enrolled_factors": schema.Int64Attribute{
			MarkdownDescription: "Maximum number of MFA factors a user can enroll",
			Optional:            true,
		},
		"mfa_phone_enroll_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable phone MFA enrollment",
			Optional:            true,
		},
		"mfa_phone_max_frequency": schema.Int64Attribute{
			MarkdownDescription: "Maximum phone MFA verification attempts per hour",
			Optional:            true,
		},
		"mfa_phone_otp_length": schema.Int64Attribute{
			MarkdownDescription: "Phone MFA OTP code length",
			Optional:            true,
		},
		"mfa_phone_template": schema.StringAttribute{
			MarkdownDescription: "Phone MFA SMS message template",
			Optional:            true,
		},
		"mfa_phone_verify_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable phone MFA verification",
			Optional:            true,
		},
		"mfa_totp_enroll_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable TOTP MFA enrollment",
			Optional:            true,
		},
		"mfa_totp_verify_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable TOTP MFA verification",
			Optional:            true,
		},
		"mfa_web_authn_enroll_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable WebAuthn MFA enrollment",
			Optional:            true,
		},
		"mfa_web_authn_verify_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable WebAuthn MFA verification",
			Optional:            true,
		},
	}
}
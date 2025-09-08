package settings

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthHooksConfig represents webhook and hook configuration
type AuthHooksConfig struct {
	// Hook/Webhook settings
	HookCustomAccessTokenEnabled           types.Bool   `tfsdk:"hook_custom_access_token_enabled"`
	HookCustomAccessTokenSecrets           types.String `tfsdk:"hook_custom_access_token_secrets"`
	HookCustomAccessTokenUri               types.String `tfsdk:"hook_custom_access_token_uri"`
	HookMfaVerificationAttemptEnabled      types.Bool   `tfsdk:"hook_mfa_verification_attempt_enabled"`
	HookMfaVerificationAttemptSecrets      types.String `tfsdk:"hook_mfa_verification_attempt_secrets"`
	HookMfaVerificationAttemptUri          types.String `tfsdk:"hook_mfa_verification_attempt_uri"`
	HookPasswordVerificationAttemptEnabled types.Bool   `tfsdk:"hook_password_verification_attempt_enabled"`
	HookPasswordVerificationAttemptSecrets types.String `tfsdk:"hook_password_verification_attempt_secrets"`
	HookPasswordVerificationAttemptUri     types.String `tfsdk:"hook_password_verification_attempt_uri"`
	HookSendEmailEnabled                   types.Bool   `tfsdk:"hook_send_email_enabled"`
	HookSendEmailSecrets                   types.String `tfsdk:"hook_send_email_secrets"`
	HookSendEmailUri                       types.String `tfsdk:"hook_send_email_uri"`
	HookSendSmsEnabled                     types.Bool   `tfsdk:"hook_send_sms_enabled"`
	HookSendSmsSecrets                     types.String `tfsdk:"hook_send_sms_secrets"`
	HookSendSmsUri                         types.String `tfsdk:"hook_send_sms_uri"`
}

func GetAuthHooksSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"hook_custom_access_token_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable custom access token hook",
			Optional:            true,
		},
		"hook_custom_access_token_secrets": schema.StringAttribute{
			MarkdownDescription: "Custom access token hook secrets",
			Optional:            true,
			Sensitive:           true,
		},
		"hook_custom_access_token_uri": schema.StringAttribute{
			MarkdownDescription: "Custom access token hook URI",
			Optional:            true,
		},
		"hook_mfa_verification_attempt_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable MFA verification attempt hook",
			Optional:            true,
		},
		"hook_mfa_verification_attempt_secrets": schema.StringAttribute{
			MarkdownDescription: "MFA verification attempt hook secrets",
			Optional:            true,
			Sensitive:           true,
		},
		"hook_mfa_verification_attempt_uri": schema.StringAttribute{
			MarkdownDescription: "MFA verification attempt hook URI",
			Optional:            true,
		},
		"hook_password_verification_attempt_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable password verification attempt hook",
			Optional:            true,
		},
		"hook_password_verification_attempt_secrets": schema.StringAttribute{
			MarkdownDescription: "Password verification attempt hook secrets",
			Optional:            true,
			Sensitive:           true,
		},
		"hook_password_verification_attempt_uri": schema.StringAttribute{
			MarkdownDescription: "Password verification attempt hook URI",
			Optional:            true,
		},
		"hook_send_email_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable send email hook",
			Optional:            true,
		},
		"hook_send_email_secrets": schema.StringAttribute{
			MarkdownDescription: "Send email hook secrets",
			Optional:            true,
			Sensitive:           true,
		},
		"hook_send_email_uri": schema.StringAttribute{
			MarkdownDescription: "Send email hook URI",
			Optional:            true,
		},
		"hook_send_sms_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable send SMS hook",
			Optional:            true,
		},
		"hook_send_sms_secrets": schema.StringAttribute{
			MarkdownDescription: "Send SMS hook secrets",
			Optional:            true,
			Sensitive:           true,
		},
		"hook_send_sms_uri": schema.StringAttribute{
			MarkdownDescription: "Send SMS hook URI",
			Optional:            true,
		},
	}
}
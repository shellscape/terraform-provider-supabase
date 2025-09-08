// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package settings

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthMailerConfig represents mailer and SMTP configuration
type AuthMailerConfig struct {
	// Email settings
	MailerAutoconfirm                 types.Bool  `tfsdk:"mailer_autoconfirm"`
	MailerAllowUnverifiedEmailSignIns types.Bool  `tfsdk:"mailer_allow_unverified_email_sign_ins"`
	MailerSecureEmailChangeEnabled    types.Bool  `tfsdk:"mailer_secure_email_change_enabled"`
	MailerOtpExp                      types.Int64 `tfsdk:"mailer_otp_exp"`
	MailerOtpLength                   types.Int64 `tfsdk:"mailer_otp_length"`

	// SMTP settings
	SmtpAdminEmail   types.String `tfsdk:"smtp_admin_email"`
	SmtpHost         types.String `tfsdk:"smtp_host"`
	SmtpMaxFrequency types.Int64  `tfsdk:"smtp_max_frequency"`
	SmtpPort         types.Int64  `tfsdk:"smtp_port"`
	SmtpSenderName   types.String `tfsdk:"smtp_sender_name"`
	SmtpUser         types.String `tfsdk:"smtp_user"`
	SmtpPass         types.String `tfsdk:"smtp_pass"`

	// Mailer templates and subjects
	MailerSubjectsConfirmation             types.String `tfsdk:"mailer_subjects_confirmation"`
	MailerSubjectsEmailChange              types.String `tfsdk:"mailer_subjects_email_change"`
	MailerSubjectsInvite                   types.String `tfsdk:"mailer_subjects_invite"`
	MailerSubjectsMagicLink                types.String `tfsdk:"mailer_subjects_magic_link"`
	MailerSubjectsReauthentication         types.String `tfsdk:"mailer_subjects_reauthentication"`
	MailerSubjectsRecovery                 types.String `tfsdk:"mailer_subjects_recovery"`
	MailerTemplatesConfirmationContent     types.String `tfsdk:"mailer_templates_confirmation_content"`
	MailerTemplatesEmailChangeContent      types.String `tfsdk:"mailer_templates_email_change_content"`
	MailerTemplatesInviteContent           types.String `tfsdk:"mailer_templates_invite_content"`
	MailerTemplatesMagicLinkContent        types.String `tfsdk:"mailer_templates_magic_link_content"`
	MailerTemplatesReauthenticationContent types.String `tfsdk:"mailer_templates_reauthentication_content"`
	MailerTemplatesRecoveryContent         types.String `tfsdk:"mailer_templates_recovery_content"`
}

func GetAuthMailerSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		// Email settings
		"mailer_autoconfirm": schema.BoolAttribute{
			MarkdownDescription: "Automatically confirm user emails",
			Optional:            true,
		},
		"mailer_allow_unverified_email_sign_ins": schema.BoolAttribute{
			MarkdownDescription: "Allow sign-ins with unverified emails",
			Optional:            true,
		},
		"mailer_secure_email_change_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable secure email change process",
			Optional:            true,
		},
		"mailer_otp_exp": schema.Int64Attribute{
			MarkdownDescription: "Email OTP expiration time in seconds",
			Optional:            true,
		},
		"mailer_otp_length": schema.Int64Attribute{
			MarkdownDescription: "Email OTP length",
			Optional:            true,
		},

		// SMTP settings
		"smtp_admin_email": schema.StringAttribute{
			MarkdownDescription: "SMTP admin email address",
			Optional:            true,
		},
		"smtp_host": schema.StringAttribute{
			MarkdownDescription: "SMTP server hostname",
			Optional:            true,
		},
		"smtp_max_frequency": schema.Int64Attribute{
			MarkdownDescription: "Maximum SMTP send frequency per hour",
			Optional:            true,
		},
		"smtp_port": schema.Int64Attribute{
			MarkdownDescription: "SMTP server port",
			Optional:            true,
		},
		"smtp_sender_name": schema.StringAttribute{
			MarkdownDescription: "SMTP sender display name",
			Optional:            true,
		},
		"smtp_user": schema.StringAttribute{
			MarkdownDescription: "SMTP username",
			Optional:            true,
		},
		"smtp_pass": schema.StringAttribute{
			MarkdownDescription: "SMTP password",
			Optional:            true,
			Sensitive:           true,
		},

		// Mailer subjects
		"mailer_subjects_confirmation": schema.StringAttribute{
			MarkdownDescription: "Email confirmation subject template",
			Optional:            true,
		},
		"mailer_subjects_email_change": schema.StringAttribute{
			MarkdownDescription: "Email change subject template",
			Optional:            true,
		},
		"mailer_subjects_invite": schema.StringAttribute{
			MarkdownDescription: "User invite subject template",
			Optional:            true,
		},
		"mailer_subjects_magic_link": schema.StringAttribute{
			MarkdownDescription: "Magic link subject template",
			Optional:            true,
		},
		"mailer_subjects_reauthentication": schema.StringAttribute{
			MarkdownDescription: "Reauthentication subject template",
			Optional:            true,
		},
		"mailer_subjects_recovery": schema.StringAttribute{
			MarkdownDescription: "Password recovery subject template",
			Optional:            true,
		},

		// Mailer templates
		"mailer_templates_confirmation_content": schema.StringAttribute{
			MarkdownDescription: "Email confirmation content template",
			Optional:            true,
		},
		"mailer_templates_email_change_content": schema.StringAttribute{
			MarkdownDescription: "Email change content template",
			Optional:            true,
		},
		"mailer_templates_invite_content": schema.StringAttribute{
			MarkdownDescription: "User invite content template",
			Optional:            true,
		},
		"mailer_templates_magic_link_content": schema.StringAttribute{
			MarkdownDescription: "Magic link content template",
			Optional:            true,
		},
		"mailer_templates_reauthentication_content": schema.StringAttribute{
			MarkdownDescription: "Reauthentication content template",
			Optional:            true,
		},
		"mailer_templates_recovery_content": schema.StringAttribute{
			MarkdownDescription: "Password recovery content template",
			Optional:            true,
		},
	}
}
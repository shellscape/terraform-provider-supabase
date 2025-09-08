package settings

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthSmsConfig represents SMS and phone authentication configuration
type AuthSmsConfig struct {
	// SMS settings
	SmsProvider  types.String `tfsdk:"sms_provider"`
	SmsOtpLength types.Int64  `tfsdk:"sms_otp_length"`

	// Additional SMS settings
	SmsAutoconfirm       types.Bool   `tfsdk:"sms_autoconfirm"`
	SmsMaxFrequency      types.Int64  `tfsdk:"sms_max_frequency"`
	SmsOtpExp            types.Int64  `tfsdk:"sms_otp_exp"`
	SmsTemplate          types.String `tfsdk:"sms_template"`
	SmsTestOtp           types.String `tfsdk:"sms_test_otp"`
	SmsTestOtpValidUntil types.String `tfsdk:"sms_test_otp_valid_until"`

	// SMS Provider specific settings
	SmsMessagebirdAccessKey          types.String `tfsdk:"sms_messagebird_access_key"`
	SmsMessagebirdOriginator         types.String `tfsdk:"sms_messagebird_originator"`
	SmsTextlocalApiKey               types.String `tfsdk:"sms_textlocal_api_key"`
	SmsTextlocalSender               types.String `tfsdk:"sms_textlocal_sender"`
	SmsTwilioAccountSid              types.String `tfsdk:"sms_twilio_account_sid"`
	SmsTwilioAuthToken               types.String `tfsdk:"sms_twilio_auth_token"`
	SmsTwilioContentSid              types.String `tfsdk:"sms_twilio_content_sid"`
	SmsTwilioMessageServiceSid       types.String `tfsdk:"sms_twilio_message_service_sid"`
	SmsTwilioVerifyAccountSid        types.String `tfsdk:"sms_twilio_verify_account_sid"`
	SmsTwilioVerifyAuthToken         types.String `tfsdk:"sms_twilio_verify_auth_token"`
	SmsTwilioVerifyMessageServiceSid types.String `tfsdk:"sms_twilio_verify_message_service_sid"`
	SmsVonageApiKey                  types.String `tfsdk:"sms_vonage_api_key"`
	SmsVonageApiSecret               types.String `tfsdk:"sms_vonage_api_secret"`
	SmsVonageFrom                    types.String `tfsdk:"sms_vonage_from"`
}

func GetAuthSmsSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		// Basic SMS settings
		"sms_provider": schema.StringAttribute{
			MarkdownDescription: "SMS provider (twilio, messagebird, textlocal, vonage)",
			Optional:            true,
		},
		"sms_otp_length": schema.Int64Attribute{
			MarkdownDescription: "SMS OTP code length",
			Optional:            true,
		},
		"sms_autoconfirm": schema.BoolAttribute{
			MarkdownDescription: "Automatically confirm SMS OTP",
			Optional:            true,
		},
		"sms_max_frequency": schema.Int64Attribute{
			MarkdownDescription: "Maximum SMS send frequency per hour",
			Optional:            true,
		},
		"sms_otp_exp": schema.Int64Attribute{
			MarkdownDescription: "SMS OTP expiration time in seconds",
			Optional:            true,
		},
		"sms_template": schema.StringAttribute{
			MarkdownDescription: "SMS message template",
			Optional:            true,
		},
		"sms_test_otp": schema.StringAttribute{
			MarkdownDescription: "Test SMS OTP code for development",
			Optional:            true,
			Sensitive:           true,
		},
		"sms_test_otp_valid_until": schema.StringAttribute{
			MarkdownDescription: "Test SMS OTP valid until timestamp",
			Optional:            true,
		},

		// MessageBird settings
		"sms_messagebird_access_key": schema.StringAttribute{
			MarkdownDescription: "MessageBird access key",
			Optional:            true,
			Sensitive:           true,
		},
		"sms_messagebird_originator": schema.StringAttribute{
			MarkdownDescription: "MessageBird originator/sender ID",
			Optional:            true,
		},

		// Textlocal settings
		"sms_textlocal_api_key": schema.StringAttribute{
			MarkdownDescription: "Textlocal API key",
			Optional:            true,
			Sensitive:           true,
		},
		"sms_textlocal_sender": schema.StringAttribute{
			MarkdownDescription: "Textlocal sender name",
			Optional:            true,
		},

		// Twilio settings
		"sms_twilio_account_sid": schema.StringAttribute{
			MarkdownDescription: "Twilio account SID",
			Optional:            true,
			Sensitive:           true,
		},
		"sms_twilio_auth_token": schema.StringAttribute{
			MarkdownDescription: "Twilio auth token",
			Optional:            true,
			Sensitive:           true,
		},
		"sms_twilio_content_sid": schema.StringAttribute{
			MarkdownDescription: "Twilio content SID",
			Optional:            true,
		},
		"sms_twilio_message_service_sid": schema.StringAttribute{
			MarkdownDescription: "Twilio message service SID",
			Optional:            true,
		},
		"sms_twilio_verify_account_sid": schema.StringAttribute{
			MarkdownDescription: "Twilio verify account SID",
			Optional:            true,
			Sensitive:           true,
		},
		"sms_twilio_verify_auth_token": schema.StringAttribute{
			MarkdownDescription: "Twilio verify auth token",
			Optional:            true,
			Sensitive:           true,
		},
		"sms_twilio_verify_message_service_sid": schema.StringAttribute{
			MarkdownDescription: "Twilio verify message service SID",
			Optional:            true,
		},

		// Vonage settings
		"sms_vonage_api_key": schema.StringAttribute{
			MarkdownDescription: "Vonage API key",
			Optional:            true,
			Sensitive:           true,
		},
		"sms_vonage_api_secret": schema.StringAttribute{
			MarkdownDescription: "Vonage API secret",
			Optional:            true,
			Sensitive:           true,
		},
		"sms_vonage_from": schema.StringAttribute{
			MarkdownDescription: "Vonage sender number or name",
			Optional:            true,
		},
	}
}
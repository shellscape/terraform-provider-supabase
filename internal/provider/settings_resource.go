// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/supabase/cli/pkg/api"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SettingsResource{}
var _ resource.ResourceWithImportState = &SettingsResource{}

func NewSettingsResource() resource.Resource {
	return &SettingsResource{}
}

// SettingsResource defines the resource implementation.
type SettingsResource struct {
	client *api.ClientWithResponses
}

// SettingsResourceModel describes the resource data model.
type SettingsResourceModel struct {
	ProjectRef types.String    `tfsdk:"project_ref"`
	Database   *DatabaseConfig `tfsdk:"database"`
	Pooler     *PoolerConfig   `tfsdk:"pooler"`
	Network    *NetworkConfig  `tfsdk:"network"`
	Storage    *StorageConfig  `tfsdk:"storage"`
	Auth       *AuthConfig     `tfsdk:"auth"`
	Api        *ApiConfig      `tfsdk:"api"`
	Id         types.String    `tfsdk:"id"`
}

// DatabaseConfig represents PostgreSQL database configuration
type DatabaseConfig struct {
	EffectiveCacheSize            types.String `tfsdk:"effective_cache_size"`
	LogicalDecodingWorkMem        types.String `tfsdk:"logical_decoding_work_mem"`
	MaintenanceWorkMem            types.String `tfsdk:"maintenance_work_mem"`
	MaxConnections                types.Int64  `tfsdk:"max_connections"`
	MaxLocksPerTransaction        types.Int64  `tfsdk:"max_locks_per_transaction"`
	MaxParallelMaintenanceWorkers types.Int64  `tfsdk:"max_parallel_maintenance_workers"`
	MaxParallelWorkers            types.Int64  `tfsdk:"max_parallel_workers"`
	MaxParallelWorkersPerGather   types.Int64  `tfsdk:"max_parallel_workers_per_gather"`
	MaxReplicationSlots           types.Int64  `tfsdk:"max_replication_slots"`
	MaxSlotWalKeepSize            types.String `tfsdk:"max_slot_wal_keep_size"`
	MaxStandbyArchiveDelay        types.String `tfsdk:"max_standby_archive_delay"`
	MaxStandbyStreamingDelay      types.String `tfsdk:"max_standby_streaming_delay"`
	MaxWalSenders                 types.Int64  `tfsdk:"max_wal_senders"`
	MaxWalSize                    types.String `tfsdk:"max_wal_size"`
	MaxWorkerProcesses            types.Int64  `tfsdk:"max_worker_processes"`
	RestartDatabase               types.Bool   `tfsdk:"restart_database"`
	SessionReplicationRole        types.String `tfsdk:"session_replication_role"`
	SharedBuffers                 types.String `tfsdk:"shared_buffers"`
	StatementTimeout              types.String `tfsdk:"statement_timeout"`
	TrackCommitTimestamp          types.Bool   `tfsdk:"track_commit_timestamp"`
	WalKeepSize                   types.String `tfsdk:"wal_keep_size"`
	WalSenderTimeout              types.String `tfsdk:"wal_sender_timeout"`
	WorkMem                       types.String `tfsdk:"work_mem"`
}

// PoolerConfig represents connection pooler configuration
type PoolerConfig struct {
	DefaultPoolSize types.Int64 `tfsdk:"default_pool_size"`
}

// NetworkConfig represents network restrictions configuration
type NetworkConfig struct {
	DbAllowedCidrs   []types.String `tfsdk:"db_allowed_cidrs"`
	DbAllowedCidrsV6 []types.String `tfsdk:"db_allowed_cidrs_v6"`
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	FileSizeLimit types.Int64         `tfsdk:"file_size_limit"`
	Features      *StorageFeatures    `tfsdk:"features"`
}

// StorageFeatures represents storage feature flags
type StorageFeatures struct {
	ImageTransformation types.Bool `tfsdk:"image_transformation"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	ApiMaxRequestDuration         types.Int64  `tfsdk:"api_max_request_duration"`
	DbMaxPoolSize                 types.Int64  `tfsdk:"db_max_pool_size"`
	DisableSignup                 types.Bool   `tfsdk:"disable_signup"`
	ExternalAnonymousUsersEnabled types.Bool   `tfsdk:"external_anonymous_users_enabled"`
	ExternalEmailEnabled          types.Bool   `tfsdk:"external_email_enabled"`
	
	// External providers
	ExternalApple     *ExternalProviderConfig `tfsdk:"external_apple"`
	ExternalAzure     *ExternalProviderConfig `tfsdk:"external_azure"`
	ExternalBitbucket *ExternalProviderConfig `tfsdk:"external_bitbucket"`
	ExternalDiscord   *ExternalProviderConfig `tfsdk:"external_discord"`
	ExternalFacebook  *ExternalProviderConfig `tfsdk:"external_facebook"`
	ExternalFigma     *ExternalProviderConfig `tfsdk:"external_figma"`
	ExternalGithub    *ExternalProviderConfig `tfsdk:"external_github"`
	ExternalGitlab    *ExternalProviderConfig `tfsdk:"external_gitlab"`
	ExternalGoogle    *ExternalProviderConfig `tfsdk:"external_google"`
	ExternalKakao     *ExternalProviderConfig `tfsdk:"external_kakao"`
	ExternalKeycloak  *ExternalProviderConfig `tfsdk:"external_keycloak"`
	ExternalLinkedinOidc *ExternalProviderConfig `tfsdk:"external_linkedin_oidc"`
	ExternalNotion    *ExternalProviderConfig `tfsdk:"external_notion"`
	ExternalSlack     *ExternalProviderConfig `tfsdk:"external_slack"`
	ExternalSlackOidc *ExternalProviderConfig `tfsdk:"external_slack_oidc"`
	ExternalSpotify   *ExternalProviderConfig `tfsdk:"external_spotify"`
	ExternalTwitch    *ExternalProviderConfig `tfsdk:"external_twitch"`
	ExternalTwitter   *ExternalProviderConfig `tfsdk:"external_twitter"`
	ExternalWorkos    *ExternalProviderConfig `tfsdk:"external_workos"`
	ExternalZoom      *ExternalProviderConfig `tfsdk:"external_zoom"`

	// SMS settings
	SmsProvider types.String `tfsdk:"sms_provider"`
	SmsOtpLength types.Int64 `tfsdk:"sms_otp_length"`
	
	// Email settings
	MailerOtpExp types.Int64 `tfsdk:"mailer_otp_exp"`
	SmtpHost     types.String `tfsdk:"smtp_host"`
	SmtpPort     types.Int64  `tfsdk:"smtp_port"`
	SmtpUser     types.String `tfsdk:"smtp_user"`
	SmtpPass     types.String `tfsdk:"smtp_pass"`

	// Security settings
	SecurityCaptchaEnabled  types.Bool   `tfsdk:"security_captcha_enabled"`
	SecurityCaptchaProvider types.String `tfsdk:"security_captcha_provider"`
	SecurityCaptchaSecret   types.String `tfsdk:"security_captcha_secret"`

	// Site settings
	SiteUrl types.String `tfsdk:"site_url"`

	// MFA settings
	MfaPhoneOtpLength types.Int64 `tfsdk:"mfa_phone_otp_length"`
}

// ExternalProviderConfig represents external OAuth provider configuration
type ExternalProviderConfig struct {
	Enabled              types.Bool   `tfsdk:"enabled"`
	ClientId             types.String `tfsdk:"client_id"`
	Secret               types.String `tfsdk:"secret"`
	Url                  types.String `tfsdk:"url"`
	AdditionalClientIds  types.String `tfsdk:"additional_client_ids"`
}


// ApiConfig represents PostgREST API configuration
type ApiConfig struct {
	DbExtraSearchPath types.String `tfsdk:"db_extra_search_path"`
	DbPool            types.Int64  `tfsdk:"db_pool"`
	DbSchema          types.String `tfsdk:"db_schema"`
	MaxRows           types.Int64  `tfsdk:"max_rows"`
}

func (r *SettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings"
}

func (r *SettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Settings resource",

		Attributes: map[string]schema.Attribute{
			"project_ref": schema.StringAttribute{
				MarkdownDescription: "Project reference ID",
				Required:            true,
			},
			"database": schema.SingleNestedAttribute{
				MarkdownDescription: "Database settings as [structured configuration](https://api.supabase.com/api/v1#/v1-update-postgres-config)",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"effective_cache_size": schema.StringAttribute{
						MarkdownDescription: "Amount of memory available for disk caching by the OS and within the database itself",
						Optional:            true,
					},
					"logical_decoding_work_mem": schema.StringAttribute{
						MarkdownDescription: "Memory used for logical decoding",
						Optional:            true,
					},
					"maintenance_work_mem": schema.StringAttribute{
						MarkdownDescription: "Maximum amount of memory to be used by maintenance operations",
						Optional:            true,
					},
					"max_connections": schema.Int64Attribute{
						MarkdownDescription: "Maximum number of concurrent connections to the database server",
						Optional:            true,
					},
					"max_locks_per_transaction": schema.Int64Attribute{
						MarkdownDescription: "Maximum number of locks per transaction",
						Optional:            true,
					},
					"max_parallel_maintenance_workers": schema.Int64Attribute{
						MarkdownDescription: "Maximum number of parallel maintenance workers",
						Optional:            true,
					},
					"max_parallel_workers": schema.Int64Attribute{
						MarkdownDescription: "Maximum number of parallel worker processes",
						Optional:            true,
					},
					"max_parallel_workers_per_gather": schema.Int64Attribute{
						MarkdownDescription: "Maximum number of parallel workers per Gather node",
						Optional:            true,
					},
					"max_replication_slots": schema.Int64Attribute{
						MarkdownDescription: "Maximum number of replication slots",
						Optional:            true,
					},
					"max_slot_wal_keep_size": schema.StringAttribute{
						MarkdownDescription: "Maximum size of WAL files that replication slots are allowed to retain",
						Optional:            true,
					},
					"max_standby_archive_delay": schema.StringAttribute{
						MarkdownDescription: "Maximum delay before canceling queries when a hot standby server is processing archived WAL data",
						Optional:            true,
					},
					"max_standby_streaming_delay": schema.StringAttribute{
						MarkdownDescription: "Maximum delay before canceling queries when a hot standby server is processing streamed WAL data",
						Optional:            true,
					},
					"max_wal_senders": schema.Int64Attribute{
						MarkdownDescription: "Maximum number of WAL sender processes",
						Optional:            true,
					},
					"max_wal_size": schema.StringAttribute{
						MarkdownDescription: "Maximum size to let the WAL grow during automatic checkpoints",
						Optional:            true,
					},
					"max_worker_processes": schema.Int64Attribute{
						MarkdownDescription: "Maximum number of background worker processes",
						Optional:            true,
					},
					"restart_database": schema.BoolAttribute{
						MarkdownDescription: "Whether to restart the database to apply configuration changes",
						Optional:            true,
					},
					"session_replication_role": schema.StringAttribute{
						MarkdownDescription: "Controls firing of replication-related triggers and rules (origin, replica, local)",
						Optional:            true,
					},
					"shared_buffers": schema.StringAttribute{
						MarkdownDescription: "Amount of memory the database server uses for shared memory buffers",
						Optional:            true,
					},
					"statement_timeout": schema.StringAttribute{
						MarkdownDescription: "Maximum allowed duration of any statement",
						Optional:            true,
					},
					"track_commit_timestamp": schema.BoolAttribute{
						MarkdownDescription: "Whether to track commit time stamps of transactions",
						Optional:            true,
					},
					"wal_keep_size": schema.StringAttribute{
						MarkdownDescription: "Minimum size to retain in the pg_wal directory",
						Optional:            true,
					},
					"wal_sender_timeout": schema.StringAttribute{
						MarkdownDescription: "Maximum time to wait for WAL replication",
						Optional:            true,
					},
					"work_mem": schema.StringAttribute{
						MarkdownDescription: "Amount of memory to be used by internal sort operations and hash tables",
						Optional:            true,
					},
				},
			},
			"pooler": schema.SingleNestedAttribute{
				MarkdownDescription: "Connection pooler settings",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"default_pool_size": schema.Int64Attribute{
						MarkdownDescription: "Default connection pool size",
						Optional:            true,
					},
				},
			},
			"network": schema.SingleNestedAttribute{
				MarkdownDescription: "Network restrictions settings",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"db_allowed_cidrs": schema.ListAttribute{
						MarkdownDescription: "List of allowed IPv4 CIDR blocks for database access",
						Optional:            true,
						ElementType:         types.StringType,
					},
					"db_allowed_cidrs_v6": schema.ListAttribute{
						MarkdownDescription: "List of allowed IPv6 CIDR blocks for database access",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"storage": schema.SingleNestedAttribute{
				MarkdownDescription: "Storage configuration settings",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"file_size_limit": schema.Int64Attribute{
						MarkdownDescription: "Maximum file size limit in bytes",
						Optional:            true,
					},
					"features": schema.SingleNestedAttribute{
						MarkdownDescription: "Storage feature flags",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"image_transformation": schema.BoolAttribute{
								MarkdownDescription: "Enable image transformation features",
								Optional:            true,
							},
						},
					},
				},
			},
			"auth": schema.SingleNestedAttribute{
				MarkdownDescription: "Auth settings as [structured configuration](https://api.supabase.com/api/v1#/v1-update-auth-service-config)",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"api_max_request_duration": schema.Int64Attribute{
						MarkdownDescription: "Maximum request duration for auth API in seconds",
						Optional:            true,
					},
					"db_max_pool_size": schema.Int64Attribute{
						MarkdownDescription: "Maximum database connection pool size for auth service",
						Optional:            true,
					},
					"disable_signup": schema.BoolAttribute{
						MarkdownDescription: "Disable new user signups",
						Optional:            true,
					},
					"external_anonymous_users_enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable anonymous users for external authentication",
						Optional:            true,
					},
					"external_email_enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable email/password authentication",
						Optional:            true,
					},
					"sms_provider": schema.StringAttribute{
						MarkdownDescription: "SMS provider (twilio, messagebird, textlocal, vonage)",
						Optional:            true,
					},
					"sms_otp_length": schema.Int64Attribute{
						MarkdownDescription: "Length of SMS OTP codes",
						Optional:            true,
					},
					"mailer_otp_exp": schema.Int64Attribute{
						MarkdownDescription: "Email OTP expiration time in seconds",
						Optional:            true,
					},
					"smtp_host": schema.StringAttribute{
						MarkdownDescription: "SMTP server hostname",
						Optional:            true,
					},
					"smtp_port": schema.Int64Attribute{
						MarkdownDescription: "SMTP server port",
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
					"security_captcha_enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable CAPTCHA for authentication",
						Optional:            true,
					},
					"security_captcha_provider": schema.StringAttribute{
						MarkdownDescription: "CAPTCHA provider",
						Optional:            true,
					},
					"security_captcha_secret": schema.StringAttribute{
						MarkdownDescription: "CAPTCHA secret key",
						Optional:            true,
						Sensitive:           true,
					},
					"site_url": schema.StringAttribute{
						MarkdownDescription: "Site URL for redirects and email links",
						Optional:            true,
					},
					"mfa_phone_otp_length": schema.Int64Attribute{
						MarkdownDescription: "Length of MFA phone OTP codes",
						Optional:            true,
					},
					
					// External OAuth Providers
					"external_apple": schema.SingleNestedAttribute{
						MarkdownDescription: "Apple OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Apple"),
					},
					"external_azure": schema.SingleNestedAttribute{
						MarkdownDescription: "Azure OAuth provider configuration", 
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Azure"),
					},
					"external_bitbucket": schema.SingleNestedAttribute{
						MarkdownDescription: "Bitbucket OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Bitbucket"),
					},
					"external_discord": schema.SingleNestedAttribute{
						MarkdownDescription: "Discord OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Discord"),
					},
					"external_facebook": schema.SingleNestedAttribute{
						MarkdownDescription: "Facebook OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Facebook"),
					},
					"external_figma": schema.SingleNestedAttribute{
						MarkdownDescription: "Figma OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Figma"),
					},
					"external_github": schema.SingleNestedAttribute{
						MarkdownDescription: "GitHub OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("GitHub"),
					},
					"external_gitlab": schema.SingleNestedAttribute{
						MarkdownDescription: "GitLab OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("GitLab"),
					},
					"external_google": schema.SingleNestedAttribute{
						MarkdownDescription: "Google OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Google"),
					},
					"external_kakao": schema.SingleNestedAttribute{
						MarkdownDescription: "Kakao OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Kakao"),
					},
					"external_keycloak": schema.SingleNestedAttribute{
						MarkdownDescription: "Keycloak OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Keycloak"),
					},
					"external_linkedin_oidc": schema.SingleNestedAttribute{
						MarkdownDescription: "LinkedIn OIDC OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("LinkedIn OIDC"),
					},
					"external_notion": schema.SingleNestedAttribute{
						MarkdownDescription: "Notion OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Notion"),
					},
					"external_slack": schema.SingleNestedAttribute{
						MarkdownDescription: "Slack OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Slack"),
					},
					"external_slack_oidc": schema.SingleNestedAttribute{
						MarkdownDescription: "Slack OIDC OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Slack OIDC"),
					},
					"external_spotify": schema.SingleNestedAttribute{
						MarkdownDescription: "Spotify OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Spotify"),
					},
					"external_twitch": schema.SingleNestedAttribute{
						MarkdownDescription: "Twitch OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Twitch"),
					},
					"external_twitter": schema.SingleNestedAttribute{
						MarkdownDescription: "Twitter OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Twitter"),
					},
					"external_workos": schema.SingleNestedAttribute{
						MarkdownDescription: "WorkOS OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("WorkOS"),
					},
					"external_zoom": schema.SingleNestedAttribute{
						MarkdownDescription: "Zoom OAuth provider configuration",
						Optional:            true,
						Attributes: getExternalProviderSchemaAttributes("Zoom"),
					},
				},
			},
			"api": schema.SingleNestedAttribute{
				MarkdownDescription: "API settings as [structured configuration](https://api.supabase.com/api/v1#/v1-update-postgrest-service-config)",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"db_extra_search_path": schema.StringAttribute{
						MarkdownDescription: "Extra search path for database schemas",
						Optional:            true,
					},
					"db_pool": schema.Int64Attribute{
						MarkdownDescription: "Database connection pool size",
						Optional:            true,
					},
					"db_schema": schema.StringAttribute{
						MarkdownDescription: "Database schemas to expose via PostgREST",
						Optional:            true,
					},
					"max_rows": schema.Int64Attribute{
						MarkdownDescription: "Maximum number of rows returned in a single request",
						Optional:            true,
					},
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Project identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *SettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*SupabaseProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *SupabaseProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = providerData.ManagementClient
}

func (r *SettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Database != nil {
		resp.Diagnostics.Append(r.updateDatabaseConfig(ctx, &data)...)
	}
	if data.Network != nil {
		resp.Diagnostics.Append(r.updateNetworkConfig(ctx, &data)...)
	}
	if data.Api != nil {
		resp.Diagnostics.Append(r.updateApiConfig(ctx, &data)...)
	}
	if data.Auth != nil {
		resp.Diagnostics.Append(r.updateAuthConfig(ctx, &data)...)
	}
	if data.Storage != nil {
		resp.Diagnostics.Append(r.updateStorageConfig(ctx, &data)...)
	}
	if data.Pooler != nil {
		resp.Diagnostics.Append(r.updatePoolerConfig(ctx, &data)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = data.ProjectRef

	tflog.Trace(ctx, "created a resource")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SettingsResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Database != nil {
		resp.Diagnostics.Append(r.readDatabaseConfig(ctx, &data)...)
	}
	if data.Network != nil {
		resp.Diagnostics.Append(r.readNetworkConfig(ctx, &data)...)
	}
	if data.Api != nil {
		resp.Diagnostics.Append(r.readApiConfig(ctx, &data)...)
	}
	if data.Auth != nil {
		resp.Diagnostics.Append(r.readAuthConfig(ctx, &data)...)
	}
	if data.Storage != nil {
		resp.Diagnostics.Append(r.readStorageConfig(ctx, &data)...)
	}
	if data.Pooler != nil {
		resp.Diagnostics.Append(r.readPoolerConfig(ctx, &data)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	data.ProjectRef = data.Id
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Database != nil {
		resp.Diagnostics.Append(r.updateDatabaseConfig(ctx, &data)...)
	}
	if data.Network != nil {
		resp.Diagnostics.Append(r.updateNetworkConfig(ctx, &data)...)
	}
	if data.Api != nil {
		resp.Diagnostics.Append(r.updateApiConfig(ctx, &data)...)
	}
	if data.Auth != nil {
		resp.Diagnostics.Append(r.updateAuthConfig(ctx, &data)...)
	}
	if data.Storage != nil {
		resp.Diagnostics.Append(r.updateStorageConfig(ctx, &data)...)
	}
	if data.Pooler != nil {
		resp.Diagnostics.Append(r.updatePoolerConfig(ctx, &data)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Simply fallthrough since there is no API to delete / reset settings.
}

func (r *SettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	data := SettingsResourceModel{Id: types.StringValue(req.ID)}

	// Initialize empty configs for import
	data.Database = &DatabaseConfig{}
	data.Network = &NetworkConfig{}
	data.Api = &ApiConfig{}
	data.Auth = &AuthConfig{}
	data.Storage = &StorageConfig{}
	data.Pooler = &PoolerConfig{}

	// Read all configs from API when importing
	resp.Diagnostics.Append(r.readDatabaseConfig(ctx, &data)...)
	resp.Diagnostics.Append(r.readNetworkConfig(ctx, &data)...)
	resp.Diagnostics.Append(r.readApiConfig(ctx, &data)...)
	resp.Diagnostics.Append(r.readAuthConfig(ctx, &data)...)
	resp.Diagnostics.Append(r.readStorageConfig(ctx, &data)...)
	resp.Diagnostics.Append(r.readPoolerConfig(ctx, &data)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Conversion and CRUD helper functions for each config type will be implemented next...

func (r *SettingsResource) readApiConfig(ctx context.Context, state *SettingsResourceModel) diag.Diagnostics {
	httpResp, err := r.client.V1GetPostgrestServiceConfigWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read api settings: %s", err))}
	}

	switch httpResp.StatusCode() {
	case http.StatusNotFound, http.StatusNotAcceptable:
		return nil
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read api settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	if state.Api == nil {
		state.Api = &ApiConfig{}
	}

	state.Api.DbExtraSearchPath = types.StringValue(httpResp.JSON200.DbExtraSearchPath)
	state.Api.DbSchema = types.StringValue(httpResp.JSON200.DbSchema)
	if httpResp.JSON200.DbPool != nil {
		state.Api.DbPool = types.Int64Value(int64(*httpResp.JSON200.DbPool))
	}
	state.Api.MaxRows = types.Int64Value(int64(httpResp.JSON200.MaxRows))

	return nil
}

func (r *SettingsResource) updateApiConfig(ctx context.Context, plan *SettingsResourceModel) diag.Diagnostics {
	body := api.UpdatePostgrestConfigBody{}

	if !plan.Api.DbExtraSearchPath.IsNull() {
		body.DbExtraSearchPath = plan.Api.DbExtraSearchPath.ValueStringPointer()
	}
	if !plan.Api.DbSchema.IsNull() {
		body.DbSchema = plan.Api.DbSchema.ValueStringPointer()
	}
	if !plan.Api.DbPool.IsNull() {
		val := int(plan.Api.DbPool.ValueInt64())
		body.DbPool = &val
	}
	if !plan.Api.MaxRows.IsNull() {
		val := int(plan.Api.MaxRows.ValueInt64())
		body.MaxRows = &val
	}

	httpResp, err := r.client.V1UpdatePostgrestServiceConfigWithResponse(ctx, plan.ProjectRef.ValueString(), body)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update api settings: %s", err))}
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update api settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	return nil
}

func (r *SettingsResource) readAuthConfig(ctx context.Context, state *SettingsResourceModel) diag.Diagnostics {
	httpResp, err := r.client.V1GetAuthServiceConfigWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read auth settings: %s", err))}
	}

	switch httpResp.StatusCode() {
	case http.StatusNotFound, http.StatusNotAcceptable:
		return nil
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read auth settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	if state.Auth == nil {
		state.Auth = &AuthConfig{}
	}

	resp := httpResp.JSON200

	if resp.ApiMaxRequestDuration != nil {
		state.Auth.ApiMaxRequestDuration = types.Int64Value(int64(*resp.ApiMaxRequestDuration))
	}
	if resp.DbMaxPoolSize != nil {
		state.Auth.DbMaxPoolSize = types.Int64Value(int64(*resp.DbMaxPoolSize))
	}
	if resp.DisableSignup != nil {
		state.Auth.DisableSignup = types.BoolValue(*resp.DisableSignup)
	}
	if resp.ExternalAnonymousUsersEnabled != nil {
		state.Auth.ExternalAnonymousUsersEnabled = types.BoolValue(*resp.ExternalAnonymousUsersEnabled)
	}
	if resp.ExternalEmailEnabled != nil {
		state.Auth.ExternalEmailEnabled = types.BoolValue(*resp.ExternalEmailEnabled)
	}
	if resp.SiteUrl != nil {
		state.Auth.SiteUrl = types.StringValue(*resp.SiteUrl)
	}
	state.Auth.MailerOtpExp = types.Int64Value(int64(resp.MailerOtpExp))
	state.Auth.MfaPhoneOtpLength = types.Int64Value(int64(resp.MfaPhoneOtpLength))
	state.Auth.SmsOtpLength = types.Int64Value(int64(resp.SmsOtpLength))
	
	// Handle SMTP settings
	if resp.SmtpHost != nil {
		state.Auth.SmtpHost = types.StringValue(*resp.SmtpHost)
	}
	if resp.SmtpPort != nil {
		state.Auth.SmtpPort = types.Int64Value(parseInt(*resp.SmtpPort))
	}
	if resp.SmtpUser != nil {
		state.Auth.SmtpUser = types.StringValue(*resp.SmtpUser)
	}
	
	// Handle external providers - only populate if they exist in config
	state.Auth.ExternalApple = readExternalProvider(state.Auth.ExternalApple, resp.ExternalAppleEnabled, resp.ExternalAppleClientId, "", resp.ExternalAppleAdditionalClientIds, nil)
	state.Auth.ExternalAzure = readExternalProvider(state.Auth.ExternalAzure, resp.ExternalAzureEnabled, resp.ExternalAzureClientId, "", nil, resp.ExternalAzureUrl)
	state.Auth.ExternalBitbucket = readExternalProvider(state.Auth.ExternalBitbucket, resp.ExternalBitbucketEnabled, resp.ExternalBitbucketClientId, "", nil, nil)
	state.Auth.ExternalDiscord = readExternalProvider(state.Auth.ExternalDiscord, resp.ExternalDiscordEnabled, resp.ExternalDiscordClientId, "", nil, nil)
	state.Auth.ExternalFacebook = readExternalProvider(state.Auth.ExternalFacebook, resp.ExternalFacebookEnabled, resp.ExternalFacebookClientId, "", nil, nil)
	state.Auth.ExternalFigma = readExternalProvider(state.Auth.ExternalFigma, resp.ExternalFigmaEnabled, resp.ExternalFigmaClientId, "", nil, nil)
	state.Auth.ExternalGithub = readExternalProvider(state.Auth.ExternalGithub, resp.ExternalGithubEnabled, resp.ExternalGithubClientId, "", nil, nil)
	state.Auth.ExternalGitlab = readExternalProvider(state.Auth.ExternalGitlab, resp.ExternalGitlabEnabled, resp.ExternalGitlabClientId, "", nil, resp.ExternalGitlabUrl)
	state.Auth.ExternalGoogle = readExternalProvider(state.Auth.ExternalGoogle, resp.ExternalGoogleEnabled, resp.ExternalGoogleClientId, "", resp.ExternalGoogleAdditionalClientIds, nil)
	state.Auth.ExternalKakao = readExternalProvider(state.Auth.ExternalKakao, resp.ExternalKakaoEnabled, resp.ExternalKakaoClientId, "", nil, nil)
	state.Auth.ExternalKeycloak = readExternalProvider(state.Auth.ExternalKeycloak, resp.ExternalKeycloakEnabled, resp.ExternalKeycloakClientId, "", nil, resp.ExternalKeycloakUrl)
	state.Auth.ExternalLinkedinOidc = readExternalProvider(state.Auth.ExternalLinkedinOidc, resp.ExternalLinkedinOidcEnabled, resp.ExternalLinkedinOidcClientId, "", nil, nil)
	state.Auth.ExternalNotion = readExternalProvider(state.Auth.ExternalNotion, resp.ExternalNotionEnabled, resp.ExternalNotionClientId, "", nil, nil)
	state.Auth.ExternalSlack = readExternalProvider(state.Auth.ExternalSlack, resp.ExternalSlackEnabled, resp.ExternalSlackClientId, "", nil, nil)
	state.Auth.ExternalSlackOidc = readExternalProvider(state.Auth.ExternalSlackOidc, resp.ExternalSlackOidcEnabled, resp.ExternalSlackOidcClientId, "", nil, nil)
	state.Auth.ExternalSpotify = readExternalProvider(state.Auth.ExternalSpotify, resp.ExternalSpotifyEnabled, resp.ExternalSpotifyClientId, "", nil, nil)
	state.Auth.ExternalTwitch = readExternalProvider(state.Auth.ExternalTwitch, resp.ExternalTwitchEnabled, resp.ExternalTwitchClientId, "", nil, nil)
	state.Auth.ExternalTwitter = readExternalProvider(state.Auth.ExternalTwitter, resp.ExternalTwitterEnabled, resp.ExternalTwitterClientId, "", nil, nil)
	state.Auth.ExternalWorkos = readExternalProvider(state.Auth.ExternalWorkos, resp.ExternalWorkosEnabled, resp.ExternalWorkosClientId, "", nil, resp.ExternalWorkosUrl)
	state.Auth.ExternalZoom = readExternalProvider(state.Auth.ExternalZoom, resp.ExternalZoomEnabled, resp.ExternalZoomClientId, "", nil, nil)

	return nil
}

func (r *SettingsResource) updateAuthConfig(ctx context.Context, plan *SettingsResourceModel) diag.Diagnostics {
	body := api.UpdateAuthConfigBody{}

	if !plan.Auth.ApiMaxRequestDuration.IsNull() {
		val := int(plan.Auth.ApiMaxRequestDuration.ValueInt64())
		body.ApiMaxRequestDuration = &val
	}
	if !plan.Auth.DbMaxPoolSize.IsNull() {
		val := int(plan.Auth.DbMaxPoolSize.ValueInt64())
		body.DbMaxPoolSize = &val
	}
	if !plan.Auth.DisableSignup.IsNull() {
		body.DisableSignup = plan.Auth.DisableSignup.ValueBoolPointer()
	}
	if !plan.Auth.ExternalAnonymousUsersEnabled.IsNull() {
		body.ExternalAnonymousUsersEnabled = plan.Auth.ExternalAnonymousUsersEnabled.ValueBoolPointer()
	}
	if !plan.Auth.ExternalEmailEnabled.IsNull() {
		body.ExternalEmailEnabled = plan.Auth.ExternalEmailEnabled.ValueBoolPointer()
	}
	if !plan.Auth.SiteUrl.IsNull() {
		body.SiteUrl = plan.Auth.SiteUrl.ValueStringPointer()
	}
	if !plan.Auth.SmtpHost.IsNull() {
		body.SmtpHost = plan.Auth.SmtpHost.ValueStringPointer()
	}
	if !plan.Auth.SmtpPort.IsNull() {
		val := fmt.Sprintf("%d", plan.Auth.SmtpPort.ValueInt64())
		body.SmtpPort = &val
	}
	if !plan.Auth.SmtpUser.IsNull() {
		body.SmtpUser = plan.Auth.SmtpUser.ValueStringPointer()
	}
	if !plan.Auth.SmtpPass.IsNull() {
		body.SmtpPass = plan.Auth.SmtpPass.ValueStringPointer()
	}
	
	// Handle external providers
	updateExternalProvider(plan.Auth.ExternalApple, &body.ExternalAppleEnabled, &body.ExternalAppleClientId, &body.ExternalAppleSecret, &body.ExternalAppleAdditionalClientIds, nil)
	updateExternalProvider(plan.Auth.ExternalAzure, &body.ExternalAzureEnabled, &body.ExternalAzureClientId, &body.ExternalAzureSecret, nil, &body.ExternalAzureUrl)
	updateExternalProvider(plan.Auth.ExternalBitbucket, &body.ExternalBitbucketEnabled, &body.ExternalBitbucketClientId, &body.ExternalBitbucketSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalDiscord, &body.ExternalDiscordEnabled, &body.ExternalDiscordClientId, &body.ExternalDiscordSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalFacebook, &body.ExternalFacebookEnabled, &body.ExternalFacebookClientId, &body.ExternalFacebookSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalFigma, &body.ExternalFigmaEnabled, &body.ExternalFigmaClientId, &body.ExternalFigmaSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalGithub, &body.ExternalGithubEnabled, &body.ExternalGithubClientId, &body.ExternalGithubSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalGitlab, &body.ExternalGitlabEnabled, &body.ExternalGitlabClientId, &body.ExternalGitlabSecret, nil, &body.ExternalGitlabUrl)
	updateExternalProvider(plan.Auth.ExternalGoogle, &body.ExternalGoogleEnabled, &body.ExternalGoogleClientId, &body.ExternalGoogleSecret, &body.ExternalGoogleAdditionalClientIds, nil)
	updateExternalProvider(plan.Auth.ExternalKakao, &body.ExternalKakaoEnabled, &body.ExternalKakaoClientId, &body.ExternalKakaoSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalKeycloak, &body.ExternalKeycloakEnabled, &body.ExternalKeycloakClientId, &body.ExternalKeycloakSecret, nil, &body.ExternalKeycloakUrl)
	updateExternalProvider(plan.Auth.ExternalLinkedinOidc, &body.ExternalLinkedinOidcEnabled, &body.ExternalLinkedinOidcClientId, &body.ExternalLinkedinOidcSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalNotion, &body.ExternalNotionEnabled, &body.ExternalNotionClientId, &body.ExternalNotionSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalSlack, &body.ExternalSlackEnabled, &body.ExternalSlackClientId, &body.ExternalSlackSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalSlackOidc, &body.ExternalSlackOidcEnabled, &body.ExternalSlackOidcClientId, &body.ExternalSlackOidcSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalSpotify, &body.ExternalSpotifyEnabled, &body.ExternalSpotifyClientId, &body.ExternalSpotifySecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalTwitch, &body.ExternalTwitchEnabled, &body.ExternalTwitchClientId, &body.ExternalTwitchSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalTwitter, &body.ExternalTwitterEnabled, &body.ExternalTwitterClientId, &body.ExternalTwitterSecret, nil, nil)
	updateExternalProvider(plan.Auth.ExternalWorkos, &body.ExternalWorkosEnabled, &body.ExternalWorkosClientId, &body.ExternalWorkosSecret, nil, &body.ExternalWorkosUrl)
	updateExternalProvider(plan.Auth.ExternalZoom, &body.ExternalZoomEnabled, &body.ExternalZoomClientId, &body.ExternalZoomSecret, nil, nil)

	httpResp, err := r.client.V1UpdateAuthServiceConfigWithResponse(ctx, plan.ProjectRef.ValueString(), body)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update auth settings: %s", err))}
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update auth settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	return nil
}

func (r *SettingsResource) readDatabaseConfig(ctx context.Context, state *SettingsResourceModel) diag.Diagnostics {
	httpResp, err := r.client.V1GetPostgresConfigWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read database settings: %s", err))}
	}

	switch httpResp.StatusCode() {
	case http.StatusNotFound, http.StatusNotAcceptable:
		return nil
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read database settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	if state.Database == nil {
		state.Database = &DatabaseConfig{}
	}

	resp := httpResp.JSON200

	if resp.EffectiveCacheSize != nil {
		state.Database.EffectiveCacheSize = types.StringValue(*resp.EffectiveCacheSize)
	}
	if resp.LogicalDecodingWorkMem != nil {
		state.Database.LogicalDecodingWorkMem = types.StringValue(*resp.LogicalDecodingWorkMem)
	}
	if resp.MaintenanceWorkMem != nil {
		state.Database.MaintenanceWorkMem = types.StringValue(*resp.MaintenanceWorkMem)
	}
	if resp.MaxConnections != nil {
		state.Database.MaxConnections = types.Int64Value(int64(*resp.MaxConnections))
	}
	if resp.StatementTimeout != nil {
		state.Database.StatementTimeout = types.StringValue(*resp.StatementTimeout)
	}
	if resp.SharedBuffers != nil {
		state.Database.SharedBuffers = types.StringValue(*resp.SharedBuffers)
	}
	if resp.WorkMem != nil {
		state.Database.WorkMem = types.StringValue(*resp.WorkMem)
	}

	return nil
}

func (r *SettingsResource) updateDatabaseConfig(ctx context.Context, plan *SettingsResourceModel) diag.Diagnostics {
	body := api.UpdatePostgresConfigBody{}

	if !plan.Database.EffectiveCacheSize.IsNull() {
		body.EffectiveCacheSize = plan.Database.EffectiveCacheSize.ValueStringPointer()
	}
	if !plan.Database.LogicalDecodingWorkMem.IsNull() {
		body.LogicalDecodingWorkMem = plan.Database.LogicalDecodingWorkMem.ValueStringPointer()
	}
	if !plan.Database.MaintenanceWorkMem.IsNull() {
		body.MaintenanceWorkMem = plan.Database.MaintenanceWorkMem.ValueStringPointer()
	}
	if !plan.Database.MaxConnections.IsNull() {
		val := int(plan.Database.MaxConnections.ValueInt64())
		body.MaxConnections = &val
	}
	if !plan.Database.StatementTimeout.IsNull() {
		body.StatementTimeout = plan.Database.StatementTimeout.ValueStringPointer()
	}
	if !plan.Database.SharedBuffers.IsNull() {
		body.SharedBuffers = plan.Database.SharedBuffers.ValueStringPointer()
	}
	if !plan.Database.WorkMem.IsNull() {
		body.WorkMem = plan.Database.WorkMem.ValueStringPointer()
	}
	if !plan.Database.RestartDatabase.IsNull() {
		body.RestartDatabase = plan.Database.RestartDatabase.ValueBoolPointer()
	}

	httpResp, err := r.client.V1UpdatePostgresConfigWithResponse(ctx, plan.ProjectRef.ValueString(), body)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update database settings: %s", err))}
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update database settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	return nil
}

func (r *SettingsResource) readNetworkConfig(ctx context.Context, state *SettingsResourceModel) diag.Diagnostics {
	httpResp, err := r.client.V1GetNetworkRestrictionsWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read network settings: %s", err))}
	}

	switch httpResp.StatusCode() {
	case http.StatusNotFound, http.StatusNotAcceptable:
		return nil
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read network settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	if state.Network == nil {
		state.Network = &NetworkConfig{}
	}

	// Initialize slices based on API response
	if v4 := httpResp.JSON200.Config.DbAllowedCidrs; v4 != nil {
		state.Network.DbAllowedCidrs = []types.String{}
		for _, cidr := range *v4 {
			state.Network.DbAllowedCidrs = append(state.Network.DbAllowedCidrs, types.StringValue(cidr))
		}
	} else {
		state.Network.DbAllowedCidrs = nil
	}
	
	if v6 := httpResp.JSON200.Config.DbAllowedCidrsV6; v6 != nil {
		state.Network.DbAllowedCidrsV6 = []types.String{}
		for _, cidr := range *v6 {
			state.Network.DbAllowedCidrsV6 = append(state.Network.DbAllowedCidrsV6, types.StringValue(cidr))
		}
	} else {
		state.Network.DbAllowedCidrsV6 = nil
	}

	return nil
}

func (r *SettingsResource) updateNetworkConfig(ctx context.Context, plan *SettingsResourceModel) diag.Diagnostics {
	body := api.NetworkRestrictionsRequest{
		DbAllowedCidrs:   &[]string{},
		DbAllowedCidrsV6: &[]string{},
	}

	for _, cidr := range plan.Network.DbAllowedCidrs {
		cidrStr := cidr.ValueString()
		ip, _, err := net.ParseCIDR(cidrStr)
		if err != nil {
			return diag.Diagnostics{diag.NewErrorDiagnostic("Validation Error", fmt.Sprintf("Invalid CIDR: %s", cidrStr))}
		}
		if ip.IsPrivate() {
			return diag.Diagnostics{diag.NewErrorDiagnostic("Validation Error", fmt.Sprintf("Private IP not allowed: %s", cidrStr))}
		}
		if ip.To4() != nil {
			*body.DbAllowedCidrs = append(*body.DbAllowedCidrs, cidrStr)
		}
	}

	for _, cidr := range plan.Network.DbAllowedCidrsV6 {
		cidrStr := cidr.ValueString()
		ip, _, err := net.ParseCIDR(cidrStr)
		if err != nil {
			return diag.Diagnostics{diag.NewErrorDiagnostic("Validation Error", fmt.Sprintf("Invalid CIDR: %s", cidrStr))}
		}
		if ip.IsPrivate() {
			return diag.Diagnostics{diag.NewErrorDiagnostic("Validation Error", fmt.Sprintf("Private IP not allowed: %s", cidrStr))}
		}
		if ip.To4() == nil {
			*body.DbAllowedCidrsV6 = append(*body.DbAllowedCidrsV6, cidrStr)
		}
	}

	httpResp, err := r.client.V1UpdateNetworkRestrictionsWithResponse(ctx, plan.ProjectRef.ValueString(), body)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update network settings: %s", err))}
	}

	if httpResp.JSON201 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update network settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	return nil
}

// Placeholder implementations for storage and pooler - these would need actual API endpoints
func (r *SettingsResource) readStorageConfig(ctx context.Context, state *SettingsResourceModel) diag.Diagnostics {
	// TODO: Implement when API endpoints are available
	return nil
}

func (r *SettingsResource) updateStorageConfig(ctx context.Context, plan *SettingsResourceModel) diag.Diagnostics {
	// TODO: Implement when API endpoints are available
	return nil
}

func (r *SettingsResource) readPoolerConfig(ctx context.Context, state *SettingsResourceModel) diag.Diagnostics {
	// TODO: Implement when API endpoints are available
	return nil
}

func (r *SettingsResource) updatePoolerConfig(ctx context.Context, plan *SettingsResourceModel) diag.Diagnostics {
	// TODO: Implement when API endpoints are available
	return nil
}

// Helper function to parse string to int64 for SMTP port
func parseInt(s string) int64 {
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0 // default to 0 on parse error
	}
	return val
}

// Helper function to generate schema attributes for external OAuth providers
func getExternalProviderSchemaAttributes(providerName string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			MarkdownDescription: fmt.Sprintf("Enable %s OAuth provider", providerName),
			Optional:            true,
		},
		"client_id": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("%s OAuth client ID", providerName),
			Optional:            true,
		},
		"secret": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("%s OAuth client secret", providerName),
			Optional:            true,
			Sensitive:           true,
		},
		"url": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("%s OAuth server URL (for self-hosted providers)", providerName),
			Optional:            true,
		},
		"additional_client_ids": schema.StringAttribute{
			MarkdownDescription: fmt.Sprintf("Additional %s client IDs (comma-separated)", providerName),
			Optional:            true,
		},
	}
}

// Helper function to read external provider configuration from API response
func readExternalProvider(existingConfig *ExternalProviderConfig, enabled *bool, clientId *string, secret string, additionalClientIds *string, url *string) *ExternalProviderConfig {
	// Only create config if it exists in Terraform plan or if provider is enabled in API
	if existingConfig == nil && (enabled == nil || !*enabled) {
		return nil
	}
	
	config := &ExternalProviderConfig{}
	if enabled != nil {
		config.Enabled = types.BoolValue(*enabled)
	}
	if clientId != nil {
		config.ClientId = types.StringValue(*clientId)
	}
	// Note: API doesn't return secrets, so we keep the planned secret
	if existingConfig != nil && !existingConfig.Secret.IsNull() {
		config.Secret = existingConfig.Secret
	}
	if additionalClientIds != nil {
		config.AdditionalClientIds = types.StringValue(*additionalClientIds)
	}
	if url != nil {
		config.Url = types.StringValue(*url)
	}
	
	return config
}

// Helper function to update external provider configuration in API request body
func updateExternalProvider(config *ExternalProviderConfig, enabled **bool, clientId **string, secret **string, additionalClientIds **string, url **string) {
	if config == nil {
		return
	}
	
	if !config.Enabled.IsNull() {
		*enabled = config.Enabled.ValueBoolPointer()
	}
	if !config.ClientId.IsNull() {
		*clientId = config.ClientId.ValueStringPointer()
	}
	if !config.Secret.IsNull() {
		*secret = config.Secret.ValueStringPointer()
	}
	if additionalClientIds != nil && !config.AdditionalClientIds.IsNull() {
		*additionalClientIds = config.AdditionalClientIds.ValueStringPointer()
	}
	if url != nil && !config.Url.IsNull() {
		*url = config.Url.ValueStringPointer()
	}
}
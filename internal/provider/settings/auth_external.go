package settings

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AuthExternalConfig represents external OAuth provider settings
type AuthExternalConfig struct {
	// External providers
	ExternalApple            *ExternalProviderConfig `tfsdk:"external_apple"`
	ExternalAppleClientId    types.String            `tfsdk:"external_apple_client_id"`
	ExternalAppleEnabled     types.Bool              `tfsdk:"external_apple_enabled"`
	ExternalAzure            *ExternalProviderConfig `tfsdk:"external_azure"`
	ExternalAzureClientId    types.String            `tfsdk:"external_azure_client_id"`
	ExternalAzureEnabled     types.Bool              `tfsdk:"external_azure_enabled"`
	ExternalBitbucket        *ExternalProviderConfig `tfsdk:"external_bitbucket"`
	ExternalDiscord          *ExternalProviderConfig `tfsdk:"external_discord"`
	ExternalDiscordClientId  types.String            `tfsdk:"external_discord_client_id"`
	ExternalDiscordEnabled   types.Bool              `tfsdk:"external_discord_enabled"`
	ExternalFacebook         *ExternalProviderConfig `tfsdk:"external_facebook"`
	ExternalFacebookClientId types.String            `tfsdk:"external_facebook_client_id"`
	ExternalFacebookEnabled  types.Bool              `tfsdk:"external_facebook_enabled"`
	ExternalFigma            *ExternalProviderConfig `tfsdk:"external_figma"`
	ExternalGithub           *ExternalProviderConfig `tfsdk:"external_github"`
	ExternalGithubClientId   types.String            `tfsdk:"external_github_client_id"`
	ExternalGithubEnabled    types.Bool              `tfsdk:"external_github_enabled"`
	ExternalGitlab           *ExternalProviderConfig `tfsdk:"external_gitlab"`
	ExternalGoogle           *ExternalProviderConfig `tfsdk:"external_google"`
	ExternalGoogleClientId   types.String            `tfsdk:"external_google_client_id"`
	ExternalGoogleEnabled    types.Bool              `tfsdk:"external_google_enabled"`
	ExternalKakao            *ExternalProviderConfig `tfsdk:"external_kakao"`
	ExternalKeycloak         *ExternalProviderConfig `tfsdk:"external_keycloak"`
	ExternalLinkedinOidc     *ExternalProviderConfig `tfsdk:"external_linkedin_oidc"`
	ExternalNotion           *ExternalProviderConfig `tfsdk:"external_notion"`
	ExternalSlack            *ExternalProviderConfig `tfsdk:"external_slack"`
	ExternalSlackOidc        *ExternalProviderConfig `tfsdk:"external_slack_oidc"`
	ExternalSpotify          *ExternalProviderConfig `tfsdk:"external_spotify"`
	ExternalTwitch           *ExternalProviderConfig `tfsdk:"external_twitch"`
	ExternalTwitter          *ExternalProviderConfig `tfsdk:"external_twitter"`
	ExternalWorkos           *ExternalProviderConfig `tfsdk:"external_workos"`
	ExternalZoom             *ExternalProviderConfig `tfsdk:"external_zoom"`

	// Additional external provider properties
	ExternalGoogleSkipNonceCheck types.Bool `tfsdk:"external_google_skip_nonce_check"`
	ExternalPhoneEnabled         types.Bool `tfsdk:"external_phone_enabled"`
	ExternalAnonymousUsersEnabled types.Bool  `tfsdk:"external_anonymous_users_enabled"`
	ExternalEmailEnabled          types.Bool  `tfsdk:"external_email_enabled"`
}

// ExternalProviderConfig represents external OAuth provider configuration
type ExternalProviderConfig struct {
	Enabled             types.Bool   `tfsdk:"enabled"`
	ClientId            types.String `tfsdk:"client_id"`
	Secret              types.String `tfsdk:"secret"`
	RedirectUri         types.String `tfsdk:"redirect_uri"`
	Url                 types.String `tfsdk:"url"`
	AdditionalClientIds types.String `tfsdk:"additional_client_ids"`
}

func GetExternalProviderSchemaAttributes(providerName string) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable " + providerName + " provider",
			Optional:            true,
		},
		"client_id": schema.StringAttribute{
			MarkdownDescription: providerName + " OAuth application client ID",
			Optional:            true,
			Sensitive:           true,
		},
		"secret": schema.StringAttribute{
			MarkdownDescription: providerName + " OAuth application secret",
			Optional:            true,
			Sensitive:           true,
		},
		"redirect_uri": schema.StringAttribute{
			MarkdownDescription: providerName + " OAuth redirect URI",
			Optional:            true,
		},
		"url": schema.StringAttribute{
			MarkdownDescription: providerName + " OAuth server URL",
			Optional:            true,
		},
		"additional_client_ids": schema.StringAttribute{
			MarkdownDescription: "Additional " + providerName + " client IDs",
			Optional:            true,
		},
	}
}

func GetAuthExternalSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"external_anonymous_users_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable anonymous users",
			Optional:            true,
		},
		"external_email_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable email/password authentication",
			Optional:            true,
		},
		"external_phone_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable phone authentication",
			Optional:            true,
		},
		"external_google_skip_nonce_check": schema.BoolAttribute{
			MarkdownDescription: "Skip nonce check for Google OAuth",
			Optional:            true,
		},

		// Direct provider config
		"external_apple": schema.SingleNestedAttribute{
			MarkdownDescription: "Apple OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Apple"),
		},
		"external_apple_client_id": schema.StringAttribute{
			MarkdownDescription: "Apple OAuth client ID (direct)",
			Optional:            true,
			Sensitive:           true,
		},
		"external_apple_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable Apple OAuth (direct)",
			Optional:            true,
		},
		"external_azure": schema.SingleNestedAttribute{
			MarkdownDescription: "Azure OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Azure"),
		},
		"external_azure_client_id": schema.StringAttribute{
			MarkdownDescription: "Azure OAuth client ID (direct)",
			Optional:            true,
			Sensitive:           true,
		},
		"external_azure_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable Azure OAuth (direct)",
			Optional:            true,
		},
		"external_bitbucket": schema.SingleNestedAttribute{
			MarkdownDescription: "Bitbucket OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Bitbucket"),
		},
		"external_discord": schema.SingleNestedAttribute{
			MarkdownDescription: "Discord OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Discord"),
		},
		"external_discord_client_id": schema.StringAttribute{
			MarkdownDescription: "Discord OAuth client ID (direct)",
			Optional:            true,
			Sensitive:           true,
		},
		"external_discord_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable Discord OAuth (direct)",
			Optional:            true,
		},
		"external_facebook": schema.SingleNestedAttribute{
			MarkdownDescription: "Facebook OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Facebook"),
		},
		"external_facebook_client_id": schema.StringAttribute{
			MarkdownDescription: "Facebook OAuth client ID (direct)",
			Optional:            true,
			Sensitive:           true,
		},
		"external_facebook_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable Facebook OAuth (direct)",
			Optional:            true,
		},
		"external_figma": schema.SingleNestedAttribute{
			MarkdownDescription: "Figma OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Figma"),
		},
		"external_github": schema.SingleNestedAttribute{
			MarkdownDescription: "GitHub OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("GitHub"),
		},
		"external_github_client_id": schema.StringAttribute{
			MarkdownDescription: "GitHub OAuth client ID (direct)",
			Optional:            true,
			Sensitive:           true,
		},
		"external_github_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable GitHub OAuth (direct)",
			Optional:            true,
		},
		"external_gitlab": schema.SingleNestedAttribute{
			MarkdownDescription: "GitLab OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("GitLab"),
		},
		"external_google": schema.SingleNestedAttribute{
			MarkdownDescription: "Google OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Google"),
		},
		"external_google_client_id": schema.StringAttribute{
			MarkdownDescription: "Google OAuth client ID (direct)",
			Optional:            true,
			Sensitive:           true,
		},
		"external_google_enabled": schema.BoolAttribute{
			MarkdownDescription: "Enable Google OAuth (direct)",
			Optional:            true,
		},
		"external_kakao": schema.SingleNestedAttribute{
			MarkdownDescription: "Kakao OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Kakao"),
		},
		"external_keycloak": schema.SingleNestedAttribute{
			MarkdownDescription: "Keycloak OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Keycloak"),
		},
		"external_linkedin_oidc": schema.SingleNestedAttribute{
			MarkdownDescription: "LinkedIn OIDC OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("LinkedIn OIDC"),
		},
		"external_notion": schema.SingleNestedAttribute{
			MarkdownDescription: "Notion OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Notion"),
		},
		"external_slack": schema.SingleNestedAttribute{
			MarkdownDescription: "Slack OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Slack"),
		},
		"external_slack_oidc": schema.SingleNestedAttribute{
			MarkdownDescription: "Slack OIDC OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Slack OIDC"),
		},
		"external_spotify": schema.SingleNestedAttribute{
			MarkdownDescription: "Spotify OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Spotify"),
		},
		"external_twitch": schema.SingleNestedAttribute{
			MarkdownDescription: "Twitch OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Twitch"),
		},
		"external_twitter": schema.SingleNestedAttribute{
			MarkdownDescription: "Twitter OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Twitter"),
		},
		"external_workos": schema.SingleNestedAttribute{
			MarkdownDescription: "WorkOS OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("WorkOS"),
		},
		"external_zoom": schema.SingleNestedAttribute{
			MarkdownDescription: "Zoom OAuth configuration",
			Optional:            true,
			Attributes:          GetExternalProviderSchemaAttributes("Zoom"),
		},
	}
}
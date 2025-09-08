// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package settings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/supabase/cli/pkg/api"
)

// AuthConfig represents the complete auth configuration  
type AuthConfig struct {
	// Embed all auth config types without tfsdk tags (Go embedding, not Terraform)
	AuthExternalConfig
	AuthLocalConfig
	AuthSecurityConfig
	AuthMailerConfig
	AuthSmsConfig
	AuthMfaConfig
	AuthHooksConfig
}

// GetAuthSchemaAttributes returns auth schema attributes
func GetAuthSchemaAttributes() map[string]schema.Attribute {
	attrs := make(map[string]schema.Attribute)
	
	// Merge all schema attributes from the separate files
	for k, v := range GetAuthExternalSchemaAttributes() {
		attrs[k] = v
	}
	for k, v := range GetAuthLocalSchemaAttributes() {
		attrs[k] = v
	}
	for k, v := range GetAuthSecuritySchemaAttributes() {
		attrs[k] = v
	}
	for k, v := range GetAuthMailerSchemaAttributes() {
		attrs[k] = v
	}
	for k, v := range GetAuthSmsSchemaAttributes() {
		attrs[k] = v
	}
	for k, v := range GetAuthMfaSchemaAttributes() {
		attrs[k] = v
	}
	for k, v := range GetAuthHooksSchemaAttributes() {
		attrs[k] = v
	}

	return attrs
}

// ReadAuthConfig reads auth configuration from the API and populates the state
func ReadAuthConfig(ctx context.Context, client *api.ClientWithResponses, state *SettingsResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	resp, err := client.V1GetAuthServiceConfigWithResponse(ctx, state.ProjectRef.ValueString())
	if err != nil {
		diags.AddError(
			"Error Reading Auth Config",
			"Could not read auth config, unexpected error: "+err.Error(),
		)
		return diags
	}

	if resp.StatusCode() != 200 {
		diags.AddError(
			"Error Reading Auth Config",
			fmt.Sprintf("Received status %d: %s", resp.StatusCode(), string(resp.Body)),
		)
		return diags
	}

	if resp.JSON200 == nil {
		diags.AddError("Error Reading Auth Config", "No response body")
		return diags
	}

	tflog.Trace(ctx, "Read auth config from API")

	// Initialize auth config if nil
	if state.Auth == nil {
		state.Auth = &AuthConfig{}
	}

	// Map API response to state - this would be a comprehensive mapping
	// For brevity, showing key examples of each type:

	// Map API response to state fields (embedded structs, so access directly)
	if resp.JSON200.ExternalGithubEnabled != nil {
		state.Auth.ExternalGithubEnabled = types.BoolPointerValue(resp.JSON200.ExternalGithubEnabled)
	}
	if resp.JSON200.ExternalGithubClientId != nil {
		state.Auth.ExternalGithubClientId = types.StringPointerValue(resp.JSON200.ExternalGithubClientId)
	}
	
	if resp.JSON200.DisableSignup != nil {
		state.Auth.DisableSignup = types.BoolPointerValue(resp.JSON200.DisableSignup)
	}
	if resp.JSON200.JwtExp != nil {
		val := int64(*resp.JSON200.JwtExp)
		state.Auth.JwtExp = types.Int64PointerValue(&val)
	}

	return diags
}

// UpdateAuthConfig updates auth configuration via the API
func UpdateAuthConfig(ctx context.Context, client *api.ClientWithResponses, plan *SettingsResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	if plan.Auth == nil {
		return diags
	}

	// Build the update request body
	body := api.UpdateAuthConfigBody{}

	// Map Terraform plan to API request body (embedded structs, so access directly)
	if !plan.Auth.DisableSignup.IsNull() && !plan.Auth.DisableSignup.IsUnknown() {
		val := plan.Auth.DisableSignup.ValueBool()
		body.DisableSignup = &val
	}
	if !plan.Auth.JwtExp.IsNull() && !plan.Auth.JwtExp.IsUnknown() {
		val := int(plan.Auth.JwtExp.ValueInt64())
		body.JwtExp = &val
	}

	// External providers
	if !plan.Auth.ExternalGithubEnabled.IsNull() && !plan.Auth.ExternalGithubEnabled.IsUnknown() {
		val := plan.Auth.ExternalGithubEnabled.ValueBool()
		body.ExternalGithubEnabled = &val
	}
	if !plan.Auth.ExternalGithubClientId.IsNull() && !plan.Auth.ExternalGithubClientId.IsUnknown() {
		val := plan.Auth.ExternalGithubClientId.ValueString()
		body.ExternalGithubClientId = &val
	}

	resp, err := client.V1UpdateAuthServiceConfigWithResponse(ctx, plan.ProjectRef.ValueString(), body)
	if err != nil {
		diags.AddError(
			"Error Updating Auth Config",
			"Could not update auth config, unexpected error: "+err.Error(),
		)
		return diags
	}

	if resp.StatusCode() != 200 {
		diags.AddError(
			"Error Updating Auth Config",
			fmt.Sprintf("Received status %d: %s", resp.StatusCode(), string(resp.Body)),
		)
		return diags
	}

	tflog.Trace(ctx, "Updated auth config via API")

	return diags
}
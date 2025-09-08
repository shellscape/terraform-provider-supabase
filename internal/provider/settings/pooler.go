// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package settings

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/supabase/cli/pkg/api"
)

// PoolerConfig represents connection pooler configuration
type PoolerConfig struct {
	DefaultPoolSize types.Int64 `tfsdk:"default_pool_size"`
}

func GetPoolerSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"default_pool_size": schema.Int64Attribute{
			MarkdownDescription: "Default connection pool size",
			Optional:            true,
		},
	}
}

// ReadPoolerConfig reads pooler configuration from the API
func ReadPoolerConfig(ctx context.Context, client *api.ClientWithResponses, state *SettingsResourceModel) diag.Diagnostics {
	httpResp, err := client.V1GetSupavisorConfigWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read pooler settings: %s", err))}
	}

	switch httpResp.StatusCode() {
	case http.StatusNotFound, http.StatusNotAcceptable:
		return nil
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read pooler settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	if state.Pooler == nil {
		state.Pooler = &PoolerConfig{}
	}

	// The API returns an array of configurations, typically we want the first one
	if len(*httpResp.JSON200) > 0 {
		config := (*httpResp.JSON200)[0]
		if config.DefaultPoolSize != nil {
			state.Pooler.DefaultPoolSize = types.Int64Value(int64(*config.DefaultPoolSize))
		}
	}

	return nil
}

// UpdatePoolerConfig updates pooler configuration via the API
func UpdatePoolerConfig(ctx context.Context, client *api.ClientWithResponses, plan *SettingsResourceModel) diag.Diagnostics {
	body := api.UpdateSupavisorConfigBody{}

	if !plan.Pooler.DefaultPoolSize.IsNull() {
		val := int(plan.Pooler.DefaultPoolSize.ValueInt64())
		body.DefaultPoolSize = &val
	}

	httpResp, err := client.V1UpdateSupavisorConfigWithResponse(ctx, plan.ProjectRef.ValueString(), body)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update pooler settings: %s", err))}
	}

	if httpResp.StatusCode() < 200 || httpResp.StatusCode() >= 300 {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update pooler settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	return nil
}

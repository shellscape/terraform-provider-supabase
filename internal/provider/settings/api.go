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

// ApiConfig represents PostgREST API configuration
type ApiConfig struct {
	DbExtraSearchPath types.String `tfsdk:"db_extra_search_path"`
	DbPool            types.Int64  `tfsdk:"db_pool"`
	DbSchema          types.String `tfsdk:"db_schema"`
	MaxRows           types.Int64  `tfsdk:"max_rows"`
}

func GetApiSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
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
	}
}

// ReadApiConfig reads API configuration from the API
func ReadApiConfig(ctx context.Context, client *api.ClientWithResponses, state *SettingsResourceModel) diag.Diagnostics {
	httpResp, err := client.V1GetPostgrestServiceConfigWithResponse(ctx, state.Id.ValueString())
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

// UpdateApiConfig updates API configuration via the API
func UpdateApiConfig(ctx context.Context, client *api.ClientWithResponses, plan *SettingsResourceModel) diag.Diagnostics {
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

	httpResp, err := client.V1UpdatePostgrestServiceConfigWithResponse(ctx, plan.ProjectRef.ValueString(), body)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update api settings: %s", err))}
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update api settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	return nil
}

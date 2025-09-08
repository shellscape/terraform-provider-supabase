// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package settings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/supabase/cli/pkg/api"
)

// SupabaseProviderData defines provider data structure
type SupabaseProviderData struct {
	ManagementClient *api.ClientWithResponses
	AccessToken      string
}

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
				Attributes:          GetDatabaseSchemaAttributes(),
			},
			"pooler": schema.SingleNestedAttribute{
				MarkdownDescription: "Connection pooler settings",
				Optional:            true,
				Attributes:          GetPoolerSchemaAttributes(),
			},
			"network": schema.SingleNestedAttribute{
				MarkdownDescription: "Network restrictions settings",
				Optional:            true,
				Attributes:          GetNetworkSchemaAttributes(),
			},
			"storage": schema.SingleNestedAttribute{
				MarkdownDescription: "Storage configuration settings",
				Optional:            true,
				Attributes:          GetStorageSchemaAttributes(),
			},
			"auth": schema.SingleNestedAttribute{
				MarkdownDescription: "Auth settings as [structured configuration](https://api.supabase.com/api/v1#/v1-update-auth-service-config)",
				Optional:            true,
				Attributes:          GetAuthSchemaAttributes(),
			},
			"api": schema.SingleNestedAttribute{
				MarkdownDescription: "API settings as [structured configuration](https://api.supabase.com/api/v1#/v1-update-postgrest-service-config)",
				Optional:            true,
				Attributes:          GetApiSchemaAttributes(),
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
		resp.Diagnostics.Append(UpdateDatabaseConfig(ctx, r.client, &data)...)
	}
	if data.Network != nil {
		resp.Diagnostics.Append(UpdateNetworkConfig(ctx, r.client, &data)...)
	}
	if data.Api != nil {
		resp.Diagnostics.Append(UpdateApiConfig(ctx, r.client, &data)...)
	}
	if data.Auth != nil {
		resp.Diagnostics.Append(UpdateAuthConfig(ctx, r.client, &data)...)
	}
	if data.Storage != nil {
		resp.Diagnostics.Append(UpdateStorageConfig(ctx, r.client, &data)...)
	}
	if data.Pooler != nil {
		resp.Diagnostics.Append(UpdatePoolerConfig(ctx, r.client, &data)...)
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
		resp.Diagnostics.Append(ReadDatabaseConfig(ctx, r.client, &data)...)
	}
	if data.Network != nil {
		resp.Diagnostics.Append(ReadNetworkConfig(ctx, r.client, &data)...)
	}
	if data.Api != nil {
		resp.Diagnostics.Append(ReadApiConfig(ctx, r.client, &data)...)
	}
	if data.Auth != nil {
		resp.Diagnostics.Append(ReadAuthConfig(ctx, r.client, &data)...)
	}
	if data.Storage != nil {
		resp.Diagnostics.Append(ReadStorageConfig(ctx, r.client, &data)...)
	}
	if data.Pooler != nil {
		resp.Diagnostics.Append(ReadPoolerConfig(ctx, r.client, &data)...)
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
		resp.Diagnostics.Append(UpdateDatabaseConfig(ctx, r.client, &data)...)
	}
	if data.Network != nil {
		resp.Diagnostics.Append(UpdateNetworkConfig(ctx, r.client, &data)...)
	}
	if data.Api != nil {
		resp.Diagnostics.Append(UpdateApiConfig(ctx, r.client, &data)...)
	}
	if data.Auth != nil {
		resp.Diagnostics.Append(UpdateAuthConfig(ctx, r.client, &data)...)
	}
	if data.Storage != nil {
		resp.Diagnostics.Append(UpdateStorageConfig(ctx, r.client, &data)...)
	}
	if data.Pooler != nil {
		resp.Diagnostics.Append(UpdatePoolerConfig(ctx, r.client, &data)...)
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
	resp.Diagnostics.Append(ReadDatabaseConfig(ctx, r.client, &data)...)
	resp.Diagnostics.Append(ReadNetworkConfig(ctx, r.client, &data)...)
	resp.Diagnostics.Append(ReadApiConfig(ctx, r.client, &data)...)
	resp.Diagnostics.Append(ReadAuthConfig(ctx, r.client, &data)...)
	resp.Diagnostics.Append(ReadStorageConfig(ctx, r.client, &data)...)
	resp.Diagnostics.Append(ReadPoolerConfig(ctx, r.client, &data)...)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}



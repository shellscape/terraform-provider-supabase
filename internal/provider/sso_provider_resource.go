package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/supabase/cli/pkg/api"
)

var _ resource.Resource = &SsoProviderResource{}

func NewSsoProviderResource() resource.Resource {
	return &SsoProviderResource{}
}

type SsoProviderResource struct {
	client *api.ClientWithResponses
}

type SsoProviderResourceModel struct {
	ProjectRef       types.String   `tfsdk:"project_ref"`
	Id               types.String   `tfsdk:"id"`
	Type             types.String   `tfsdk:"type"`
	MetadataUrl      types.String   `tfsdk:"metadata_url"`
	MetadataXml      types.String   `tfsdk:"metadata_xml"`
	Domains          []types.String `tfsdk:"domains"`
	AttributeMapping types.String   `tfsdk:"attribute_mapping"`
	CreatedAt        types.String   `tfsdk:"created_at"`
	UpdatedAt        types.String   `tfsdk:"updated_at"`
}

func (r *SsoProviderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_provider"
}

func (r *SsoProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Supabase SSO Provider resource",
		Attributes: map[string]schema.Attribute{
			"project_ref": schema.StringAttribute{
				MarkdownDescription: "Project reference",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Provider ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Provider type (saml)",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metadata_url": schema.StringAttribute{
				MarkdownDescription: "SAML metadata URL",
				Optional:            true,
			},
			"metadata_xml": schema.StringAttribute{
				MarkdownDescription: "SAML metadata XML",
				Optional:            true,
				Sensitive:           true,
			},
			"domains": schema.ListAttribute{
				MarkdownDescription: "List of domains for this provider",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"attribute_mapping": schema.StringAttribute{
				MarkdownDescription: "JSON string of attribute mapping",
				Optional:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Creation timestamp",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Last update timestamp",
				Computed:            true,
			},
		},
	}
}

func (r *SsoProviderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SsoProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SsoProviderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := api.CreateProviderBody{
		Type: api.CreateProviderBodyType(data.Type.ValueString()),
	}

	if !data.MetadataUrl.IsNull() {
		body.MetadataUrl = data.MetadataUrl.ValueStringPointer()
	}

	if !data.MetadataXml.IsNull() {
		body.MetadataXml = data.MetadataXml.ValueStringPointer()
	}

	if len(data.Domains) > 0 {
		domains := make([]string, len(data.Domains))
		for i, domain := range data.Domains {
			domains[i] = domain.ValueString()
		}
		body.Domains = &domains
	}

	httpResp, err := r.client.V1CreateASsoProviderWithResponse(ctx, data.ProjectRef.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create SSO provider: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusCreated {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create SSO provider, got status %d: %s", httpResp.StatusCode(), httpResp.Body))
		return
	}

	if httpResp.JSON201 == nil {
		resp.Diagnostics.AddError("API Error", "Empty response from create SSO provider")
		return
	}

	provider := httpResp.JSON201
	data.Id = types.StringValue(provider.Id)

	if provider.CreatedAt != nil {
		data.CreatedAt = types.StringValue(*provider.CreatedAt)
	} else {
		data.CreatedAt = types.StringValue("")
	}

	if provider.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(*provider.UpdatedAt)
	} else {
		data.UpdatedAt = types.StringValue("")
	}

	if provider.Domains != nil {
		data.Domains = make([]types.String, len(*provider.Domains))
		for i, domain := range *provider.Domains {
			if domain.Domain != nil {
				data.Domains[i] = types.StringValue(*domain.Domain)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SsoProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SsoProviderResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.V1GetASsoProviderWithResponse(ctx, data.ProjectRef.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SSO provider: %s", err))
		return
	}

	if httpResp.StatusCode() == http.StatusNotFound {
		resp.State.RemoveResource(ctx)
		return
	}

	if httpResp.JSON200 == nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read SSO provider, got status %d: %s", httpResp.StatusCode(), httpResp.Body))
		return
	}

	provider := httpResp.JSON200
	data.Id = types.StringValue(provider.Id)

	if provider.CreatedAt != nil {
		data.CreatedAt = types.StringValue(*provider.CreatedAt)
	} else {
		data.CreatedAt = types.StringValue("")
	}

	if provider.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(*provider.UpdatedAt)
	} else {
		data.UpdatedAt = types.StringValue("")
	}

	if provider.Domains != nil {
		data.Domains = make([]types.String, len(*provider.Domains))
		for i, domain := range *provider.Domains {
			if domain.Domain != nil {
				data.Domains[i] = types.StringValue(*domain.Domain)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SsoProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SsoProviderResourceModel
	var state SsoProviderResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	body := api.UpdateProviderBody{}

	if !data.MetadataUrl.IsNull() {
		body.MetadataUrl = data.MetadataUrl.ValueStringPointer()
	}

	if !data.MetadataXml.IsNull() {
		body.MetadataXml = data.MetadataXml.ValueStringPointer()
	}

	if len(data.Domains) > 0 {
		domains := make([]string, len(data.Domains))
		for i, domain := range data.Domains {
			domains[i] = domain.ValueString()
		}
		body.Domains = &domains
	}

	httpResp, err := r.client.V1UpdateASsoProviderWithResponse(ctx, data.ProjectRef.ValueString(), data.Id.ValueString(), body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update SSO provider: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update SSO provider, got status %d: %s", httpResp.StatusCode(), httpResp.Body))
		return
	}

	if httpResp.JSON200 == nil {
		resp.Diagnostics.AddError("API Error", "Empty response from update SSO provider")
		return
	}

	// Preserve immutable fields from current state
	data.Id = state.Id
	data.CreatedAt = state.CreatedAt

	provider := httpResp.JSON200
	if provider.UpdatedAt != nil {
		data.UpdatedAt = types.StringValue(*provider.UpdatedAt)
	} else {
		data.UpdatedAt = types.StringValue("")
	}

	if provider.Domains != nil {
		data.Domains = make([]types.String, len(*provider.Domains))
		for i, domain := range *provider.Domains {
			if domain.Domain != nil {
				data.Domains[i] = types.StringValue(*domain.Domain)
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SsoProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SsoProviderResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.V1DeleteASsoProviderWithResponse(ctx, data.ProjectRef.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete SSO provider: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK && httpResp.StatusCode() != http.StatusNotFound {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to delete SSO provider, got status %d: %s", httpResp.StatusCode(), httpResp.Body))
		return
	}
}
package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/supabase/cli/pkg/api"
)

var _ datasource.DataSource = &SsoProvidersDataSource{}

func NewSsoProvidersDataSource() datasource.DataSource {
	return &SsoProvidersDataSource{}
}

type SsoProvidersDataSource struct {
	client *api.ClientWithResponses
}

type SsoProvidersDataSourceModel struct {
	ProjectRef types.String            `tfsdk:"project_ref"`
	Providers  []SsoProviderDataSource `tfsdk:"providers"`
}

type SsoProviderDataSource struct {
	Id        types.String   `tfsdk:"id"`
	Type      types.String   `tfsdk:"type"`
	Domains   []types.String `tfsdk:"domains"`
	CreatedAt types.String   `tfsdk:"created_at"`
	UpdatedAt types.String   `tfsdk:"updated_at"`
}

func (d *SsoProvidersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sso_providers"
}

func (d *SsoProvidersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve a list of SSO providers for a Supabase project",
		Attributes: map[string]schema.Attribute{
			"project_ref": schema.StringAttribute{
				MarkdownDescription: "Project reference",
				Required:            true,
			},
			"providers": schema.ListNestedAttribute{
				MarkdownDescription: "List of SSO providers",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Provider ID",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Provider type",
							Computed:            true,
						},
						"domains": schema.ListAttribute{
							MarkdownDescription: "List of domains for this provider",
							Computed:            true,
							ElementType:         types.StringType,
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
				},
			},
		},
	}
}

func (d *SsoProvidersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*SupabaseProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *SupabaseProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = providerData.ManagementClient
}

func (d *SsoProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SsoProvidersDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := d.client.V1ListAllSsoProviderWithResponse(ctx, data.ProjectRef.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read SSO providers: %s", err))
		return
	}

	if httpResp.StatusCode() != http.StatusOK {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to read SSO providers, got status %d: %s", httpResp.StatusCode(), httpResp.Body))
		return
	}

	if httpResp.JSON200 == nil {
		resp.Diagnostics.AddError("API Error", "Empty response from list SSO providers")
		return
	}

	providers := httpResp.JSON200.Items
	data.Providers = make([]SsoProviderDataSource, len(providers))

	for i, provider := range providers {
		data.Providers[i].Id = types.StringValue(provider.Id)
		
		if provider.Saml != nil {
			data.Providers[i].Type = types.StringValue("saml")
		}

		if provider.CreatedAt != nil {
			data.Providers[i].CreatedAt = types.StringValue(*provider.CreatedAt)
		}

		if provider.UpdatedAt != nil {
			data.Providers[i].UpdatedAt = types.StringValue(*provider.UpdatedAt)
		}

		if provider.Domains != nil {
			data.Providers[i].Domains = make([]types.String, len(*provider.Domains))
			for j, domain := range *provider.Domains {
				if domain.Domain != nil {
					data.Providers[i].Domains[j] = types.StringValue(*domain.Domain)
				}
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
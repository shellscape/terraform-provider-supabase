package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shellscape/terraform-provider-supabase/internal/provider/settings"
	"github.com/supabase/cli/pkg/api"
)

var (
	_ datasource.DataSource              = &StorageBucketsDataSource{}
	_ datasource.DataSourceWithConfigure = &StorageBucketsDataSource{}
)

func NewStorageBucketsDataSource() datasource.DataSource {
	return &StorageBucketsDataSource{}
}

type StorageBucketsDataSource struct {
	client *api.ClientWithResponses
}

type StorageBucketsDataSourceModel struct {
	ProjectRef types.String                `tfsdk:"project_ref"`
	Buckets    []StorageBucketModel        `tfsdk:"buckets"`
	Id         types.String                `tfsdk:"id"`
}

type StorageBucketModel struct {
	Id        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Owner     types.String `tfsdk:"owner"`
	Public    types.Bool   `tfsdk:"public"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func (d *StorageBucketsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_buckets"
}

func (d *StorageBucketsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `

Retrieves a list of storage buckets for a Supabase project.

## Example Usage

~~~hcl
data "supabase_storage_buckets" "all" {
  project_ref = "abcdefghijklmnopqrst"
}
~~~
`,
		Attributes: map[string]schema.Attribute{
			"project_ref": schema.StringAttribute{
				MarkdownDescription: "Project reference ID",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Same as project_ref",
				Computed:            true,
			},
			"buckets": schema.ListNestedAttribute{
				MarkdownDescription: "List of storage buckets",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Bucket ID",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Bucket name",
							Computed:            true,
						},
						"owner": schema.StringAttribute{
							MarkdownDescription: "Bucket owner ID",
							Computed:            true,
						},
						"public": schema.BoolAttribute{
							MarkdownDescription: "Whether the bucket is public",
							Computed:            true,
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "When the bucket was created",
							Computed:            true,
						},
						"updated_at": schema.StringAttribute{
							MarkdownDescription: "When the bucket was last updated",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *StorageBucketsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*settings.SupabaseProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *settings.SupabaseProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = providerData.ManagementClient
}

func (d *StorageBucketsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data StorageBucketsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to list all buckets
	response, err := d.client.V1ListAllBucketsWithResponse(ctx, data.ProjectRef.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to list storage buckets, got error: %s", err))
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"API Error",
			fmt.Sprintf("Unable to list storage buckets, got status %d", response.StatusCode()),
		)
		return
	}

	if response.JSON200 == nil {
		resp.Diagnostics.AddError("API Error", "Empty response from storage buckets API")
		return
	}

	// Convert API response to data model
	data.Id = data.ProjectRef
	data.Buckets = make([]StorageBucketModel, 0, len(*response.JSON200))

	for _, bucket := range *response.JSON200 {
		bucketModel := StorageBucketModel{
			Id:        types.StringValue(bucket.Id),
			Name:      types.StringValue(bucket.Name),
			Owner:     types.StringValue(bucket.Owner),
			Public:    types.BoolValue(bucket.Public),
			CreatedAt: types.StringValue(bucket.CreatedAt),
			UpdatedAt: types.StringValue(bucket.UpdatedAt),
		}
		data.Buckets = append(data.Buckets, bucketModel)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
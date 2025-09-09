package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/shellscape/terraform-provider-supabase/internal/provider/settings"
	"github.com/supabase/cli/pkg/fetcher"
	"github.com/supabase/cli/pkg/storage"
)

var (
	_ resource.Resource              = &StorageBucketResource{}
	_ resource.ResourceWithConfigure = &StorageBucketResource{}
)

func NewStorageBucketResource() resource.Resource {
	return &StorageBucketResource{}
}

type StorageBucketResource struct {
	providerData *settings.SupabaseProviderData
	storageClient *storage.StorageAPI
}

type StorageBucketResourceModel struct {
	ProjectRef       types.String `tfsdk:"project_ref"`
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Public           types.Bool   `tfsdk:"public"`
	FileSizeLimit    types.Int64  `tfsdk:"file_size_limit"`
	AllowedMimeTypes types.List   `tfsdk:"allowed_mime_types"`
	CreatedAt        types.String `tfsdk:"created_at"`
	UpdatedAt        types.String `tfsdk:"updated_at"`
	Owner            types.String `tfsdk:"owner"`
}




func (r *StorageBucketResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_bucket"
}

func (r *StorageBucketResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `

Manages a Supabase storage bucket.

Refer to the [Supabase Storage documentation](https://supabase.com/docs/guides/storage) for more information.

## Example Usage

~~~hcl
resource "supabase_storage_bucket" "example" {
  project_ref    = "abcdefghijklmnopqrst"
  name           = "my-bucket"
  public         = false
  file_size_limit = 52428800  # 50MB
  allowed_mime_types = ["image/*", "video/mp4"]
}
~~~
`,
		Attributes: map[string]schema.Attribute{
			"project_ref": schema.StringAttribute{
				MarkdownDescription: "Project reference ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Bucket ID",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Bucket name (must be unique within project). Must contain only lowercase letters, numbers, dots, and hyphens.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z0-9.-]+$`),
						"bucket name must contain only lowercase letters, numbers, dots, and hyphens",
					),
					stringvalidator.LengthBetween(3, 63),
				},
			},
			"public": schema.BoolAttribute{
				MarkdownDescription: "Whether the bucket is publicly accessible",
				Required:            true,
			},
			"file_size_limit": schema.Int64Attribute{
				MarkdownDescription: "Maximum file size in bytes (null for no limit)",
				Optional:            true,
			},
			"allowed_mime_types": schema.ListAttribute{
				MarkdownDescription: "Allowed MIME types (null for no restriction). Use wildcards like 'image/*'",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "When the bucket was created",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "When the bucket was last updated",
				Computed:            true,
			},
			"owner": schema.StringAttribute{
				MarkdownDescription: "Bucket owner ID",
				Computed:            true,
			},
		},
	}
}

func (r *StorageBucketResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*settings.SupabaseProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *settings.SupabaseProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.providerData = providerData
}

func (r *StorageBucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data StorageBucketResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get storage client
	storageClient, err := r.getStorageClient(ctx, data.ProjectRef.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Storage Client Error", fmt.Sprintf("Unable to create storage client: %s", err))
		return
	}

	// Prepare the create request
	createReq := storage.CreateBucketRequest{
		Id:     data.Name.ValueString(),
		Name:   data.Name.ValueString(),
		Public: &[]bool{data.Public.ValueBool()}[0],
	}

	if !data.FileSizeLimit.IsNull() {
		createReq.FileSizeLimit = data.FileSizeLimit.ValueInt64()
	}

	if !data.AllowedMimeTypes.IsNull() {
		var mimeTypes []string
		resp.Diagnostics.Append(data.AllowedMimeTypes.ElementsAs(ctx, &mimeTypes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.AllowedMimeTypes = mimeTypes
	}

	// Create the bucket via CLI
	_, err = storageClient.CreateBucket(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Storage API Error", fmt.Sprintf("Unable to create storage bucket: %s", err))
		return
	}

	// Read back the created bucket to get all computed fields
	bucket, err := r.getBucket(ctx, storageClient, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Storage API Error", fmt.Sprintf("Unable to read created storage bucket: %s", err))
		return
	}

	// Update the data model with response values
	r.updateDataFromBucket(&data, bucket)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageBucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data StorageBucketResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get storage client
	storageClient, err := r.getStorageClient(ctx, data.ProjectRef.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Storage Client Error", fmt.Sprintf("Unable to create storage client: %s", err))
		return
	}

	// Get the bucket details
	bucket, err := r.getBucket(ctx, storageClient, data.Name.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			// Bucket no longer exists
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Storage API Error", fmt.Sprintf("Unable to read storage bucket: %s", err))
		return
	}

	// Update the data model with response values
	r.updateDataFromBucket(&data, bucket)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageBucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data StorageBucketResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get storage client
	storageClient, err := r.getStorageClient(ctx, data.ProjectRef.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Storage Client Error", fmt.Sprintf("Unable to create storage client: %s", err))
		return
	}

	// Prepare the update request
	updateReq := storage.UpdateBucketRequest{
		Id:     data.Name.ValueString(),
		Public: &[]bool{data.Public.ValueBool()}[0],
	}

	if !data.FileSizeLimit.IsNull() {
		updateReq.FileSizeLimit = data.FileSizeLimit.ValueInt64()
	}

	if !data.AllowedMimeTypes.IsNull() {
		var mimeTypes []string
		resp.Diagnostics.Append(data.AllowedMimeTypes.ElementsAs(ctx, &mimeTypes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.AllowedMimeTypes = mimeTypes
	}

	// Update the bucket via CLI
	_, err = storageClient.UpdateBucket(ctx, updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Storage API Error", fmt.Sprintf("Unable to update storage bucket: %s", err))
		return
	}

	// Read back the updated bucket to get all computed fields
	bucket, err := r.getBucket(ctx, storageClient, data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Storage API Error", fmt.Sprintf("Unable to read updated storage bucket: %s", err))
		return
	}

	// Update the data model with response values
	r.updateDataFromBucket(&data, bucket)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *StorageBucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StorageBucketResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get storage client
	storageClient, err := r.getStorageClient(ctx, data.ProjectRef.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Storage Client Error", fmt.Sprintf("Unable to create storage client: %s", err))
		return
	}

	// Delete the bucket via CLI
	_, err = storageClient.DeleteBucket(ctx, data.Name.ValueString())
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError("Storage API Error", fmt.Sprintf("Unable to delete storage bucket: %s", err))
		return
	}
}

// Helper methods for CLI-based storage operations
func (r *StorageBucketResource) getStorageClient(ctx context.Context, projectRef string) (*storage.StorageAPI, error) {
	// Get service role token using same pattern as existing implementation
	serviceRoleToken, err := r.providerData.TokenManager.GetServiceRoleToken(ctx, projectRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get service role token: %w", err)
	}

	// Create storage client similar to how CLI does it
	storageURL := fmt.Sprintf("https://%s.supabase.co", projectRef)
	client := &storage.StorageAPI{
		Fetcher: fetcher.NewFetcher(
			storageURL,
			fetcher.WithBearerToken(serviceRoleToken),
			fetcher.WithUserAgent("terraform-provider-supabase"),
		),
	}

	return client, nil
}

func (r *StorageBucketResource) getBucket(ctx context.Context, storageClient *storage.StorageAPI, bucketId string) (*storage.BucketResponse, error) {
	buckets, err := storageClient.ListBuckets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}

	for _, bucket := range buckets {
		if bucket.Id == bucketId || bucket.Name == bucketId {
			return &bucket, nil
		}
	}

	return nil, fmt.Errorf("bucket not found: 404")
}

func (r *StorageBucketResource) updateDataFromBucket(data *StorageBucketResourceModel, bucket *storage.BucketResponse) {
	data.Id = types.StringValue(bucket.Id)
	data.Name = types.StringValue(bucket.Name)
	data.Public = types.BoolValue(bucket.Public)
	data.CreatedAt = types.StringValue(bucket.CreatedAt)
	data.UpdatedAt = types.StringValue(bucket.UpdatedAt)
	data.Owner = types.StringValue(bucket.Owner)

	if bucket.FileSizeLimit != nil {
		data.FileSizeLimit = types.Int64Value(int64(*bucket.FileSizeLimit))
	} else {
		data.FileSizeLimit = types.Int64Null()
	}

	if bucket.AllowedMimeTypes != nil && len(bucket.AllowedMimeTypes) > 0 {
		mimeTypes := make([]types.String, len(bucket.AllowedMimeTypes))
		for i, mt := range bucket.AllowedMimeTypes {
			mimeTypes[i] = types.StringValue(mt)
		}
		data.AllowedMimeTypes, _ = types.ListValueFrom(context.Background(), types.StringType, mimeTypes)
	} else {
		// Only set to empty list if the original config had an empty list
		// Otherwise preserve the null state to maintain consistency
		if !data.AllowedMimeTypes.IsNull() {
			// Config specified allowed_mime_types (even if empty), so use empty list
			data.AllowedMimeTypes, _ = types.ListValueFrom(context.Background(), types.StringType, []types.String{})
		}
		// If data.AllowedMimeTypes.IsNull(), leave it as null to maintain consistency
	}
}
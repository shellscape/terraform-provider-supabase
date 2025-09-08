package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type StorageBucket struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	Public           bool     `json:"public"`
	FileSizeLimit    *int64   `json:"file_size_limit"`
	AllowedMimeTypes []string `json:"allowed_mime_types"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
	Owner            string   `json:"owner"`
}

type CreateBucketRequest struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	Public           bool     `json:"public"`
	FileSizeLimit    *int64   `json:"file_size_limit,omitempty"`
	AllowedMimeTypes []string `json:"allowed_mime_types,omitempty"`
}

type UpdateBucketRequest struct {
	Id               string   `json:"id"`
	Name             string   `json:"name"`
	Public           bool     `json:"public"`
	FileSizeLimit    *int64   `json:"file_size_limit,omitempty"`
	AllowedMimeTypes []string `json:"allowed_mime_types,omitempty"`
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

	// Prepare the create request
	createReq := CreateBucketRequest{
		Id:     data.Name.ValueString(),
		Name:   data.Name.ValueString(),
		Public: data.Public.ValueBool(),
	}

	if !data.FileSizeLimit.IsNull() {
		limit := data.FileSizeLimit.ValueInt64()
		createReq.FileSizeLimit = &limit
	}

	if !data.AllowedMimeTypes.IsNull() {
		var mimeTypes []string
		resp.Diagnostics.Append(data.AllowedMimeTypes.ElementsAs(ctx, &mimeTypes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.AllowedMimeTypes = mimeTypes
	}

	// Create the bucket via Storage API
	bucket, err := r.createBucketViaStorageAPI(ctx, data.ProjectRef.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Storage API Error", fmt.Sprintf("Unable to create storage bucket: %s", err))
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

	// Get the bucket details via Storage API
	bucket, err := r.getBucketViaStorageAPI(ctx, data.ProjectRef.ValueString(), data.Name.ValueString())
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

	// Prepare the update request
	updateReq := UpdateBucketRequest{
		Id:     data.Name.ValueString(),
		Name:   data.Name.ValueString(),
		Public: data.Public.ValueBool(),
	}

	if !data.FileSizeLimit.IsNull() {
		limit := data.FileSizeLimit.ValueInt64()
		updateReq.FileSizeLimit = &limit
	}

	if !data.AllowedMimeTypes.IsNull() {
		var mimeTypes []string
		resp.Diagnostics.Append(data.AllowedMimeTypes.ElementsAs(ctx, &mimeTypes, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		updateReq.AllowedMimeTypes = mimeTypes
	}

	// Update the bucket via Storage API
	bucket, err := r.updateBucketViaStorageAPI(ctx, data.ProjectRef.ValueString(), data.Name.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Storage API Error", fmt.Sprintf("Unable to update storage bucket: %s", err))
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

	// Delete the bucket via Storage API
	err := r.deleteBucketViaStorageAPI(ctx, data.ProjectRef.ValueString(), data.Name.ValueString())
	if err != nil && !strings.Contains(err.Error(), "404") {
		resp.Diagnostics.AddError("Storage API Error", fmt.Sprintf("Unable to delete storage bucket: %s", err))
		return
	}
}

// Helper methods for direct Storage API communication
func (r *StorageBucketResource) createBucketViaStorageAPI(ctx context.Context, projectRef string, bucket CreateBucketRequest) (*StorageBucket, error) {
	return r.makeStorageAPIRequest(ctx, "POST", projectRef, "/bucket", bucket)
}

func (r *StorageBucketResource) getBucketViaStorageAPI(ctx context.Context, projectRef string, bucketId string) (*StorageBucket, error) {
	return r.makeStorageAPIRequest(ctx, "GET", projectRef, fmt.Sprintf("/bucket/%s", bucketId), nil)
}

func (r *StorageBucketResource) updateBucketViaStorageAPI(ctx context.Context, projectRef string, bucketId string, bucket UpdateBucketRequest) (*StorageBucket, error) {
	return r.makeStorageAPIRequest(ctx, "PUT", projectRef, fmt.Sprintf("/bucket/%s", bucketId), bucket)
}

func (r *StorageBucketResource) deleteBucketViaStorageAPI(ctx context.Context, projectRef string, bucketId string) error {
	_, err := r.makeStorageAPIRequest(ctx, "DELETE", projectRef, fmt.Sprintf("/bucket/%s", bucketId), nil)
	return err
}

func (r *StorageBucketResource) makeStorageAPIRequest(ctx context.Context, method, projectRef, path string, body interface{}) (*StorageBucket, error) {
	// Storage API URL
	url := fmt.Sprintf("https://%s.supabase.co/storage/v1%s", projectRef, path)

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Get the appropriate authorization header for storage operations
	authHeader, err := r.getAuthorizationHeader(ctx, projectRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get authorization header: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)

	httpClient := &http.Client{}
	httpResp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		// If authentication failed, invalidate cached tokens and potentially retry
		if httpResp.StatusCode == 401 || httpResp.StatusCode == 403 {
			r.providerData.TokenManager.InvalidateProjectTokens(projectRef)
		}
		return nil, fmt.Errorf("API error (status %d): %s", httpResp.StatusCode, string(respBody))
	}

	// For GET requests, parse the bucket response
	if method == "GET" {
		var bucket StorageBucket
		if err := json.Unmarshal(respBody, &bucket); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		return &bucket, nil
	}

	// For POST (create), the response contains the created bucket
	if method == "POST" {
		// Create might return just {name: "bucket-name"} or the full bucket object
		// Let's fetch the full bucket details after creation
		return r.getBucketViaStorageAPI(ctx, projectRef, body.(CreateBucketRequest).Id)
	}

	// For PUT (update), fetch the updated bucket
	if method == "PUT" {
		return r.getBucketViaStorageAPI(ctx, projectRef, body.(UpdateBucketRequest).Id)
	}

	// For DELETE, return nil (no content expected)
	return nil, nil
}

func (r *StorageBucketResource) getAuthorizationHeader(ctx context.Context, projectRef string) (string, error) {
	// Use token manager to get the appropriate token for storage operations
	serviceRoleToken, err := r.providerData.TokenManager.GetServiceRoleToken(ctx, projectRef)
	if err != nil {
		return "", fmt.Errorf("failed to get service role token: %w", err)
	}
	return "Bearer " + serviceRoleToken, nil
}

func (r *StorageBucketResource) updateDataFromBucket(data *StorageBucketResourceModel, bucket *StorageBucket) {
	data.Id = types.StringValue(bucket.Id)
	data.Name = types.StringValue(bucket.Name)
	data.Public = types.BoolValue(bucket.Public)
	data.CreatedAt = types.StringValue(bucket.CreatedAt)
	data.UpdatedAt = types.StringValue(bucket.UpdatedAt)
	data.Owner = types.StringValue(bucket.Owner)

	if bucket.FileSizeLimit != nil {
		data.FileSizeLimit = types.Int64Value(*bucket.FileSizeLimit)
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
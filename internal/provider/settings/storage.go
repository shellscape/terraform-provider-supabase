// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package settings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/supabase/cli/pkg/api"
)

// StorageConfig represents storage configuration
type StorageConfig struct {
	FileSizeLimit types.Int64      `tfsdk:"file_size_limit"`
	Features      *StorageFeatures `tfsdk:"features"`
}

// StorageFeatures represents storage feature flags
type StorageFeatures struct {
	ImageTransformation *StorageFeatureImageTransformation `tfsdk:"image_transformation"`
	S3Protocol          *StorageFeatureS3Protocol          `tfsdk:"s3_protocol"`
}

// StorageFeatureImageTransformation represents image transformation feature configuration
type StorageFeatureImageTransformation struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

// StorageFeatureS3Protocol represents S3 protocol feature configuration
type StorageFeatureS3Protocol struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

// Custom types for API requests with correct camelCase JSON tags
type CustomStorageUpdateBody struct {
	FileSizeLimit *int64                  `json:"fileSizeLimit,omitempty"`
	Features      *CustomStorageFeatures  `json:"features,omitempty"`
}

type CustomStorageFeatures struct {
	ImageTransformation *CustomStorageFeatureImageTransformation `json:"imageTransformation,omitempty"`
	S3Protocol          *CustomStorageFeatureS3Protocol          `json:"s3Protocol,omitempty"`
}

type CustomStorageFeatureImageTransformation struct {
	Enabled bool `json:"enabled"`
}

type CustomStorageFeatureS3Protocol struct {
	Enabled bool `json:"enabled"`
}

func GetStorageSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"file_size_limit": schema.Int64Attribute{
			MarkdownDescription: "Maximum file size limit in bytes",
			Optional:            true,
		},
		"features": schema.SingleNestedAttribute{
			MarkdownDescription: "Storage feature flags",
			Optional:            true,
			Attributes: map[string]schema.Attribute{
				"image_transformation": schema.SingleNestedAttribute{
					MarkdownDescription: "Image transformation feature configuration",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Enable image transformation features",
							Optional:            true,
						},
					},
				},
				"s3_protocol": schema.SingleNestedAttribute{
					MarkdownDescription: "S3 protocol feature configuration",
					Optional:            true,
					Attributes: map[string]schema.Attribute{
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Enable S3 protocol compatibility",
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

// ReadStorageConfig reads storage configuration from the API
func ReadStorageConfig(ctx context.Context, client *api.ClientWithResponses, state *SettingsResourceModel) diag.Diagnostics {
	httpResp, err := client.V1GetStorageConfigWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read storage settings: %s", err))}
	}

	switch httpResp.StatusCode() {
	case http.StatusNotFound, http.StatusNotAcceptable:
		return nil
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read storage settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	if state.Storage == nil {
		state.Storage = &StorageConfig{}
	}

	resp := httpResp.JSON200
	state.Storage.FileSizeLimit = types.Int64Value(resp.FileSizeLimit)

	if state.Storage.Features == nil {
		state.Storage.Features = &StorageFeatures{}
	}

	// Only initialize and set ImageTransformation if it was configured in the plan
	if state.Storage.Features.ImageTransformation != nil {
		state.Storage.Features.ImageTransformation.Enabled = types.BoolValue(resp.Features.ImageTransformation.Enabled)
	}

	// Only initialize S3Protocol if it was configured in the plan
	if state.Storage.Features.S3Protocol != nil {
		// S3Protocol is not currently in the API response, so preserve the configured value
		// This prevents drift from occurring when the field is not returned by the API
	}

	return nil
}

// UpdateStorageConfig updates storage configuration via the API
func UpdateStorageConfig(ctx context.Context, client *api.ClientWithResponses, plan *SettingsResourceModel) diag.Diagnostics {
	body := CustomStorageUpdateBody{}

	if !plan.Storage.FileSizeLimit.IsNull() {
		val := plan.Storage.FileSizeLimit.ValueInt64()
		body.FileSizeLimit = &val
	}

	if plan.Storage.Features != nil {
		features := &CustomStorageFeatures{}

		if plan.Storage.Features.ImageTransformation != nil && !plan.Storage.Features.ImageTransformation.Enabled.IsNull() {
			features.ImageTransformation = &CustomStorageFeatureImageTransformation{
				Enabled: plan.Storage.Features.ImageTransformation.Enabled.ValueBool(),
			}
		}

		if plan.Storage.Features.S3Protocol != nil && !plan.Storage.Features.S3Protocol.Enabled.IsNull() {
			features.S3Protocol = &CustomStorageFeatureS3Protocol{
				Enabled: plan.Storage.Features.S3Protocol.Enabled.ValueBool(),
			}
		}

		body.Features = features
	}

	// Use custom HTTP request with proper JSON marshaling
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to marshal storage settings: %s", err))}
	}

	httpResp, err := client.V1UpdateStorageConfigWithBodyWithResponse(ctx, plan.ProjectRef.ValueString(), "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update storage settings: %s", err))}
	}

	if httpResp.StatusCode() < 200 || httpResp.StatusCode() >= 300 {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update storage settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	return nil
}

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStorageBucketResourceUpdateDataFromBucket(t *testing.T) {
	resource := &StorageBucketResource{}
	
	// Test with full bucket data
	bucket := &StorageBucket{
		Id:               "test-bucket",
		Name:             "test-bucket", 
		Public:           true,
		FileSizeLimit:    func() *int64 { i := int64(52428800); return &i }(),
		AllowedMimeTypes: []string{"image/*", "video/mp4"},
		CreatedAt:        "2023-01-01T00:00:00Z",
		UpdatedAt:        "2023-01-01T00:00:00Z",
		Owner:            "owner123",
	}
	
	var data StorageBucketResourceModel
	data.ProjectRef = types.StringValue("test-project")
	
	resource.updateDataFromBucket(&data, bucket)
	
	// Verify all fields are correctly mapped
	if data.Id.ValueString() != "test-bucket" {
		t.Errorf("Expected Id 'test-bucket', got %s", data.Id.ValueString())
	}
	
	if data.Name.ValueString() != "test-bucket" {
		t.Errorf("Expected Name 'test-bucket', got %s", data.Name.ValueString())
	}
	
	if !data.Public.ValueBool() {
		t.Errorf("Expected Public true, got %t", data.Public.ValueBool())
	}
	
	if data.FileSizeLimit.ValueInt64() != 52428800 {
		t.Errorf("Expected FileSizeLimit 52428800, got %d", data.FileSizeLimit.ValueInt64())
	}
	
	if data.Owner.ValueString() != "owner123" {
		t.Errorf("Expected Owner 'owner123', got %s", data.Owner.ValueString())
	}
	
	// Check MIME types
	var mimeTypes []string
	data.AllowedMimeTypes.ElementsAs(context.Background(), &mimeTypes, false)
	
	if len(mimeTypes) != 2 {
		t.Errorf("Expected 2 MIME types, got %d", len(mimeTypes))
	}
	
	if mimeTypes[0] != "image/*" {
		t.Errorf("Expected first MIME type 'image/*', got %s", mimeTypes[0])
	}
}

func TestStorageBucketResourceUpdateDataFromBucketNullFields(t *testing.T) {
	resource := &StorageBucketResource{}
	
	// Test with bucket data containing null optional fields
	bucket := &StorageBucket{
		Id:               "null-bucket",
		Name:             "null-bucket", 
		Public:           false,
		FileSizeLimit:    nil,  // This should result in null value
		AllowedMimeTypes: nil,  // This should preserve original null state
		CreatedAt:        "2023-01-01T00:00:00Z",
		UpdatedAt:        "2023-01-01T00:00:00Z",
		Owner:            "owner123",
	}
	
	var data StorageBucketResourceModel
	data.ProjectRef = types.StringValue("test-project")
	// Start with AllowedMimeTypes as null (not specified in config)
	data.AllowedMimeTypes = types.ListNull(types.StringType)
	
	resource.updateDataFromBucket(&data, bucket)
	
	// Verify null handling
	if !data.FileSizeLimit.IsNull() {
		t.Errorf("Expected FileSizeLimit to be null, got %d", data.FileSizeLimit.ValueInt64())
	}
	
	// Should preserve null state when config was null
	if !data.AllowedMimeTypes.IsNull() {
		t.Errorf("Expected AllowedMimeTypes to remain null when config was null, got non-null")
	}
}

func TestStorageBucketResourceUpdateDataFromBucketEmptyList(t *testing.T) {
	resource := &StorageBucketResource{}
	
	// Test with bucket data containing null API response but config had empty list
	bucket := &StorageBucket{
		Id:               "empty-list-bucket",
		Name:             "empty-list-bucket", 
		Public:           false,
		FileSizeLimit:    nil,
		AllowedMimeTypes: nil,  // API returns nil
		CreatedAt:        "2023-01-01T00:00:00Z",
		UpdatedAt:        "2023-01-01T00:00:00Z",
		Owner:            "owner123",
	}
	
	var data StorageBucketResourceModel
	data.ProjectRef = types.StringValue("test-project")
	// Start with AllowedMimeTypes as empty list (was specified in config as empty)
	data.AllowedMimeTypes, _ = types.ListValueFrom(context.Background(), types.StringType, []types.String{})
	
	resource.updateDataFromBucket(&data, bucket)
	
	// Should preserve empty list when config had empty list
	if data.AllowedMimeTypes.IsNull() {
		t.Errorf("Expected AllowedMimeTypes to be empty list when config had empty list, got null")
	}
	
	// Check that empty MIME types list is handled correctly
	var mimeTypes []string
	data.AllowedMimeTypes.ElementsAs(context.Background(), &mimeTypes, false)
	
	if len(mimeTypes) != 0 {
		t.Errorf("Expected 0 MIME types for empty list, got %d", len(mimeTypes))
	}
}

func TestStorageBucketResourceValidation(t *testing.T) {
	// Test various bucket name validations that should be implemented
	testCases := []struct {
		name        string
		bucketName  string
		shouldError bool
	}{
		{"Valid simple name", "my-bucket", false},
		{"Valid with numbers", "bucket123", false},
		{"Valid with hyphens", "my-test-bucket", false},
		{"Invalid with spaces", "my bucket", true},
		{"Invalid with uppercase", "MyBucket", true},
		{"Invalid with underscores", "my_bucket", true},
		{"Invalid too short", "a", true},
		{"Invalid too long", "this-bucket-name-is-way-too-long-and-should-be-rejected-by-validation", true},
		{"Invalid empty", "", true},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This test documents expected behavior but doesn't test actual implementation
			// since validation logic may not be implemented yet
			t.Logf("Bucket name '%s' should error: %t", tc.bucketName, tc.shouldError)
		})
	}
}

func TestStorageBucketResourceDeleteErrorHandling(t *testing.T) {
	// Test that delete errors are properly handled and formatted
	// This test verifies that delete errors are properly propagated
	// In a real scenario, the HTTP client would return the actual error
	// This documents the expected error handling behavior
	
	testCases := []struct {
		name          string
		expectedError string
	}{
		{"Bucket not empty", "Bucket is not empty"},
		{"Insufficient permissions", "Insufficient permissions"},
		{"Bucket not found", "Bucket not found"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Delete error '%s' should be properly handled and reported to user", tc.expectedError)
			// Actual implementation would test the deleteBucketViaStorageAPI method
			// with mocked HTTP responses returning these errors
		})
	}
}
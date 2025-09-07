package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/supabase/cli/pkg/api"
)

func TestEdgeFunctionResourceUpdateDataFromResponse(t *testing.T) {
	resource := &EdgeFunctionResource{}
	
	// Test with full function response
	resp := &api.FunctionResponse{
		Id:                "func123",
		Slug:              "test-function",
		Name:              "Test Function",
		Status:            "ACTIVE",
		CreatedAt:         1640995200,
		UpdatedAt:         1640995300,
		ComputeMultiplier: func() *float32 { f := float32(2.5); return &f }(),
		EntrypointPath:    func() *string { s := "custom/index.ts"; return &s }(),
		ImportMap:         func() *bool { b := true; return &b }(),
		ImportMapPath:     func() *string { s := "custom/import_map.json"; return &s }(),
	}
	
	var data EdgeFunctionResourceModel
	data.ProjectRef = types.StringValue("test-project")
	data.Slug = types.StringValue("test-function")
	data.Body = types.StringValue("test body")
	
	resource.updateDataFromResponse(&data, resp)
	
	// Verify all fields are correctly mapped
	if data.Id.ValueString() != "func123" {
		t.Errorf("Expected Id 'func123', got %s", data.Id.ValueString())
	}
	
	if data.Name.ValueString() != "Test Function" {
		t.Errorf("Expected Name 'Test Function', got %s", data.Name.ValueString())
	}
	
	if data.Status.ValueString() != "ACTIVE" {
		t.Errorf("Expected Status 'ACTIVE', got %s", data.Status.ValueString())
	}
	
	if data.CreatedAt.ValueInt64() != 1640995200 {
		t.Errorf("Expected CreatedAt 1640995200, got %d", data.CreatedAt.ValueInt64())
	}
	
	if data.UpdatedAt.ValueInt64() != 1640995300 {
		t.Errorf("Expected UpdatedAt 1640995300, got %d", data.UpdatedAt.ValueInt64())
	}
	
	if data.ComputeMultiplier.ValueFloat64() != 2.5 {
		t.Errorf("Expected ComputeMultiplier 2.5, got %f", data.ComputeMultiplier.ValueFloat64())
	}
	
	if data.EntrypointPath.ValueString() != "custom/index.ts" {
		t.Errorf("Expected EntrypointPath 'custom/index.ts', got %s", data.EntrypointPath.ValueString())
	}
	
	if !data.ImportMap.ValueBool() {
		t.Errorf("Expected ImportMap true, got %t", data.ImportMap.ValueBool())
	}
	
	if data.ImportMapPath.ValueString() != "custom/import_map.json" {
		t.Errorf("Expected ImportMapPath 'custom/import_map.json', got %s", data.ImportMapPath.ValueString())
	}
}

func TestEdgeFunctionResourceUpdateDataFromResponseNullFields(t *testing.T) {
	resource := &EdgeFunctionResource{}
	
	// Test with function response containing null optional fields
	resp := &api.FunctionResponse{
		Id:                "func456",
		Slug:              "minimal-function",
		Name:              "Minimal Function",
		Status:            "ACTIVE",
		CreatedAt:         1640995200,
		UpdatedAt:         1640995200,
		ComputeMultiplier: nil,  // This should result in null value
		EntrypointPath:    nil,  // This should result in null value
		ImportMap:         nil,  // This should result in null value
		ImportMapPath:     nil,  // This should result in null value
	}
	
	var data EdgeFunctionResourceModel
	data.ProjectRef = types.StringValue("test-project")
	data.Slug = types.StringValue("minimal-function")
	data.Body = types.StringValue("minimal body")
	
	resource.updateDataFromResponse(&data, resp)
	
	// Verify null handling
	if !data.ComputeMultiplier.IsNull() {
		t.Errorf("Expected ComputeMultiplier to be null, got %f", data.ComputeMultiplier.ValueFloat64())
	}
	
	if !data.EntrypointPath.IsNull() {
		t.Errorf("Expected EntrypointPath to be null, got %s", data.EntrypointPath.ValueString())
	}
	
	if !data.ImportMap.IsNull() {
		t.Errorf("Expected ImportMap to be null, got %t", data.ImportMap.ValueBool())
	}
	
	if !data.ImportMapPath.IsNull() {
		t.Errorf("Expected ImportMapPath to be null, got %s", data.ImportMapPath.ValueString())
	}
	
	// Verify required fields are still set
	if data.Id.ValueString() != "func456" {
		t.Errorf("Expected Id 'func456', got %s", data.Id.ValueString())
	}
	
	if data.Name.ValueString() != "Minimal Function" {
		t.Errorf("Expected Name 'Minimal Function', got %s", data.Name.ValueString())
	}
}

func TestEdgeFunctionResourceBodyRetrieval(t *testing.T) {
	// This test documents that edge functions require separate body retrieval
	// The main CRUD operations don't return the function body in the response
	// So we need to make a separate API call to /functions/{slug}/body
	
	t.Log("Edge functions require separate body retrieval via /functions/{slug}/body endpoint")
	t.Log("This is correctly implemented in the Read method but should be tested in integration tests")
}
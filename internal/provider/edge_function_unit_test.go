package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/supabase/cli/pkg/api"
)

func TestEdgeFunctionResourceUpdateDataFromResponse(t *testing.T) {
	resource := &EdgeFunctionResource{}
	
	// Test with basic function data
	response := &api.FunctionResponse{
		Id:        "func123",
		Slug:      "test-function",
		Name:      "Test Function",
		Status:    "ACTIVE",
		CreatedAt: 1640995200,
		UpdatedAt: 1640995200,
	}
	
	var data EdgeFunctionResourceModel
	data.ProjectRef = types.StringValue("test-project")
	
	resource.updateDataFromResponse(&data, response)
	
	// Verify all fields are correctly mapped
	if data.Id.ValueString() != "func123" {
		t.Errorf("Expected Id 'func123', got %s", data.Id.ValueString())
	}
	
	if data.Slug.ValueString() != "test-function" {
		t.Errorf("Expected Slug 'test-function', got %s", data.Slug.ValueString())
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
}

func TestEdgeFunctionResourceUpdateDataFromResponseWithOptionalFields(t *testing.T) {
	resource := &EdgeFunctionResource{}
	
	// Test with optional fields set
	computeMultiplier := float32(2.5)
	entrypointPath := "/path/to/index.ts"
	importMapPath := "/path/to/import_map.json"
	
	response := &api.FunctionResponse{
		Id:                "func456", 
		Slug:              "advanced-function",
		Name:              "Advanced Function",
		Status:            "INACTIVE",
		CreatedAt:         1640995300,
		UpdatedAt:         1640995400,
		ComputeMultiplier: &computeMultiplier,
		EntrypointPath:    &entrypointPath,
		ImportMapPath:     &importMapPath,
	}
	
	var data EdgeFunctionResourceModel
	data.ProjectRef = types.StringValue("test-project")
	
	resource.updateDataFromResponse(&data, response)
	
	// Verify optional fields are correctly mapped
	if data.ComputeMultiplier.ValueFloat64() != float64(computeMultiplier) {
		t.Errorf("Expected ComputeMultiplier %f, got %f", float64(computeMultiplier), data.ComputeMultiplier.ValueFloat64())
	}
	
	if data.EntrypointPath.ValueString() != entrypointPath {
		t.Errorf("Expected EntrypointPath '%s', got %s", entrypointPath, data.EntrypointPath.ValueString())
	}
	
	if data.ImportMapPath.ValueString() != importMapPath {
		t.Errorf("Expected ImportMapPath '%s', got %s", importMapPath, data.ImportMapPath.ValueString())
	}
}

func TestEdgeFunctionResourceUpdateDataFromResponseWithNullOptionalFields(t *testing.T) {
	resource := &EdgeFunctionResource{}
	
	// Test with optional fields as nil
	response := &api.FunctionResponse{
		Id:                "func789",
		Slug:              "simple-function", 
		Name:              "Simple Function",
		Status:            "ACTIVE",
		CreatedAt:         1640995500,
		UpdatedAt:         1640995600,
		ComputeMultiplier: nil,
		EntrypointPath:    nil,
		ImportMapPath:     nil,
	}
	
	var data EdgeFunctionResourceModel
	data.ProjectRef = types.StringValue("test-project")
	
	resource.updateDataFromResponse(&data, response)
	
	// Verify optional fields are null when not provided
	if !data.ComputeMultiplier.IsNull() {
		t.Errorf("Expected ComputeMultiplier to be null, got %f", data.ComputeMultiplier.ValueFloat64())
	}
	
	if !data.EntrypointPath.IsNull() {
		t.Errorf("Expected EntrypointPath to be null, got %s", data.EntrypointPath.ValueString())
	}
	
	if !data.ImportMapPath.IsNull() {
		t.Errorf("Expected ImportMapPath to be null, got %s", data.ImportMapPath.ValueString())
	}
}

func TestBoolPtrHelperFunction(t *testing.T) {
	// Test boolPtr function
	truePtr := boolPtr(true)
	if truePtr == nil || *truePtr != true {
		t.Errorf("Expected boolPtr(true) to return pointer to true")
	}
	
	falsePtr := boolPtr(false)
	if falsePtr == nil || *falsePtr != false {
		t.Errorf("Expected boolPtr(false) to return pointer to false")
	}
}

func TestBoolPtrFromTypesHelperFunction(t *testing.T) {
	// Test boolPtrFromTypes with valid bool
	validBool := types.BoolValue(true)
	ptr := boolPtrFromTypes(validBool)
	if ptr == nil || *ptr != true {
		t.Errorf("Expected boolPtrFromTypes(true) to return pointer to true")
	}
	
	// Test boolPtrFromTypes with null
	nullBool := types.BoolNull()
	ptr = boolPtrFromTypes(nullBool)
	if ptr != nil {
		t.Errorf("Expected boolPtrFromTypes(null) to return nil")
	}
	
	// Test boolPtrFromTypes with unknown
	unknownBool := types.BoolUnknown()
	ptr = boolPtrFromTypes(unknownBool)
	if ptr != nil {
		t.Errorf("Expected boolPtrFromTypes(unknown) to return nil")
	}
}

func TestToFileURLHelperFunction(t *testing.T) {
	// Test with empty path
	result := toFileURL("")
	if result != nil {
		t.Errorf("Expected toFileURL(\"\") to return nil")
	}
	
	// Test with valid path
	path := "/path/to/file.ts"
	result = toFileURL(path)
	if result == nil || *result != path {
		t.Errorf("Expected toFileURL(\"%s\") to return pointer to same string", path)
	}
}
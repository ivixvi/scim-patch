package scimpatch_test

import (
	"testing"

	"github.com/elimity-com/scim"
	"github.com/ivixvi/scimpatch"
)

func TestPatcher_Apply(t *testing.T) {
	// Create a Patcher instance with a dummy schema
	patcher := scimpatch.Patcher{}

	// Define a dummy PatchOperation
	op := scim.PatchOperation{
		Op:    "add",
		Value: map[string]interface{}{"displayName": "Alice Green"},
	}

	// Define the input data
	data := scim.ResourceAttributes{
		"displayName": "Bob Green",
	}

	// Apply the PatchOperation
	result, err := patcher.Apply(op, data)
	if err != nil {
		t.Fatalf("Apply() returned an unexpected error: %v", err)
	}

	// Check if the result matches the input data
	// NOTE: now Apply is not implemented. so no change value.
	if result["displayName"] != "Bob Green" {
		t.Errorf("expected %v, got %v", data, result)
	}
}

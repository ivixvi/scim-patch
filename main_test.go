package scimpatch_test

import (
	"fmt"
	"testing"

	"github.com/elimity-com/scim"
	"github.com/ivixvi/scimpatch"
	filter "github.com/scim2/filter-parser/v2"
)

func path(s string) *filter.Path {
	p, err := filter.ParsePath([]byte(s))
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to parse %s occurred by %s", s, err))
	}
	return &p
}

func TestPatcher_Apply(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name     string
		op       scim.PatchOperation
		data     scim.ResourceAttributes
		expected scim.ResourceAttributes
	}{
		{
			name: "Add operation",
			op: scim.PatchOperation{
				Op:    "add",
				Value: map[string]interface{}{"displayName": "Alice Green"},
			},
			data: scim.ResourceAttributes{
				"displayName": "Bob Green",
			},
			expected: scim.ResourceAttributes{
				"displayName": "Bob Green",
			},
		},
		{
			name: "Replace operation",
			op: scim.PatchOperation{
				Op:    "replace",
				Value: map[string]interface{}{"displayName": "Alice Green"},
			},
			data: scim.ResourceAttributes{
				"displayName": "Bob Green",
			},
			expected: scim.ResourceAttributes{
				"displayName": "Bob Green",
			},
		},
		{
			name: "Replace operation",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("displayName"),
			},
			data: scim.ResourceAttributes{
				"displayName": "Bob Green",
			},
			expected: scim.ResourceAttributes{
				"displayName": "Bob Green",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			// Create a Patcher instance with a dummy schema
			patcher := scimpatch.Patcher{}

			// Apply the PatchOperation
			result, err := patcher.Apply(tc.op, tc.data)
			if err != nil {
				t.Fatalf("Apply() returned an unexpected error: %v", err)
			}

			// Check if the result matches the expected data
			for key, expectedValue := range tc.expected {
				if result[key] != expectedValue {
					t.Errorf("for key %q, expected %v, got %v", key, expectedValue, result[key])
				}
			}
		})
	}
}

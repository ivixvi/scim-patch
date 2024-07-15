package scimpatch_test

import (
	"fmt"
	"testing"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/schema"
	scimpatch "github.com/ivixvi/scim-patch"
)

// TestPathSpecifiedReplace は Pacher.Apply の path指定をしているreplace操作の正常系をテストします
func TestPathSpecifiedReplace(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name            string
		op              scim.PatchOperation
		data            scim.ResourceAttributes
		expected        scim.ResourceAttributes
		expectedChanged bool
	}{
		// Replace Singular Attribute
		{
			name: "Replace operation - Core Singular Attributes",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`displayName`),
				Value: "Alice Green",
			},
			data: scim.ResourceAttributes{
				"displayName": "Bob Green",
			},
			expected: scim.ResourceAttributes{
				"displayName": "Alice Green",
			},
			expectedChanged: true,
		},
		{
			name: "Replace operation - Core MultiValued Complex Attributes - For Item",
			op: scim.PatchOperation{
				Op:   "replace",
				Path: path(`emails[type eq "work"]`),
				Value: map[string]interface{}{
					"type":  "work",
					"value": "ivixvi-updated@example.com",
				},
			},
			data: scim.ResourceAttributes{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi@example.com",
					},
				},
			},
			expected: scim.ResourceAttributes{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi-updated@example.com",
					},
				},
			},
			expectedChanged: true,
		},
		{
			name: "Replace operation - Core MultiValued Complex Attributes - For Attribute",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`emails[type eq "work"].value`),
				Value: "ivixvi-updated@example.com",
			},
			data: scim.ResourceAttributes{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi@example.com",
					},
				},
			},
			expected: scim.ResourceAttributes{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi-updated@example.com",
					},
				},
			},
			expectedChanged: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			// Create a Patcher instance with a dummy schema
			patcher := scimpatch.NewPatcher(schema.CoreUserSchema(), []schema.Schema{schema.ExtensionEnterpriseUser()})

			// Apply the PatchOperation
			result, changed, err := patcher.Apply(tc.op, tc.data)
			if err != nil {
				t.Fatalf("Apply() returned an unexpected error: %v", err)
			}
			// Check if the result matches the expected data
			if changed != tc.expectedChanged {
				t.Errorf("changed:\n    actual  : %v\n    expected: %v", changed, tc.expectedChanged)
			}
			// Check if the result matches the expected data
			if !(fmt.Sprint(result) == fmt.Sprint(tc.expected)) {
				t.Errorf("result:\n    actual  : %v\n    expected: %v", result, tc.expected)
			}
		})
	}
}

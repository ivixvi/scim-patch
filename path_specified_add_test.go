package scimpatch_test

import (
	"fmt"
	"testing"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/schema"
	scimpatch "github.com/ivixvi/scim-patch"
)

// TestPathSpecifiedAdd は Pacher.Apply の path指定をしているadd操作の正常系をテストします
func TestPathSpecifiedAdd(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name            string
		op              scim.PatchOperation
		data            scim.ResourceAttributes
		expected        scim.ResourceAttributes
		expectedChanged bool
	}{
		{
			name: "Add operation",
			op: scim.PatchOperation{
				Op:    "add",
				Path:  path(`displayName`),
				Value: "Alice Green",
			},
			data: scim.ResourceAttributes{},
			expected: scim.ResourceAttributes{
				"displayName": "Alice Green",
			},
			expectedChanged: true,
		},
		{
			name: "Add operation - MultiValued Complex Attributes - Filter & Value addition",
			op: scim.PatchOperation{
				Op:    "add",
				Path:  path(`emails[type eq "work"].value`),
				Value: "ivixvi@example.com",
			},
			data: scim.ResourceAttributes{},
			expected: scim.ResourceAttributes{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi@example.com",
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

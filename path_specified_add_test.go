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
		// Singular Attribute
		{
			name: "Add operation - Core Singular Attributes - add",
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
			name: "Add operation - Core Singular Attributes - replace",
			op: scim.PatchOperation{
				Op:    "add",
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
			name: "Add operation - Core Singular Attributes - no changed",
			op: scim.PatchOperation{
				Op:    "add",
				Path:  path(`displayName`),
				Value: "Alice Green",
			},
			data: scim.ResourceAttributes{
				"displayName": "Alice Green",
			},
			expected: scim.ResourceAttributes{
				"displayName": "Alice Green",
			},
			expectedChanged: false,
		},
		// Complex Attribute
		{
			name: "Add operation - Core Complex Attributes - replace",
			op: scim.PatchOperation{
				Op:    "add",
				Path:  path(`name.familyName`),
				Value: "Green",
			},
			data: scim.ResourceAttributes{
				"name": map[string]interface{}{
					"familyName": "Blue",
					"givenName":  "Alice",
				},
			},
			expected: scim.ResourceAttributes{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expectedChanged: true,
		},
		{
			name: "Add operation - Core Complex Attributes - no value",
			op: scim.PatchOperation{
				Op:    "add",
				Path:  path(`name.familyName`),
				Value: "Green",
			},
			data: scim.ResourceAttributes{
				"name": map[string]interface{}{
					"givenName": "Alice",
				},
			},
			expected: scim.ResourceAttributes{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expectedChanged: true,
		},
		{
			name: "Add operation - Core Complex Attributes - no changed.",
			op: scim.PatchOperation{
				Op:    "add",
				Path:  path(`name.familyName`),
				Value: "Green",
			},
			data: scim.ResourceAttributes{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expected: scim.ResourceAttributes{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expectedChanged: false,
		},
		{
			name: "Replace operation - Core Complex Attributes - map specified merge",
			op: scim.PatchOperation{
				Op:   "add",
				Path: path(`name`),
				Value: map[string]interface{}{
					"givenName": "Alice",
					"formatted": "Alice Green",
				},
			},
			data: scim.ResourceAttributes{
				"name": map[string]interface{}{
					"familyName": "Green",
					"formatted":  "Bob Green",
				},
			},
			expected: scim.ResourceAttributes{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
					"formatted":  "Alice Green",
				},
			},
			expectedChanged: true,
		},
		{
			name: "Replace operation - Core Singular Attributes - map specified no changed",
			op: scim.PatchOperation{
				Op:   "replace",
				Path: path(`name`),
				Value: map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			data: scim.ResourceAttributes{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expected: scim.ResourceAttributes{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expectedChanged: false,
		},
		// MultiValued Complex Attribute
		{
			name: "Add operation - Core MultiValued Complex Attributes - Add All",
			op: scim.PatchOperation{
				Op:   "add",
				Path: path(`emails`),
				Value: []interface{}{
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi-added@example.com",
					},
				},
			},
			data: scim.ResourceAttributes{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "home",
						"value": "ivixvi@example.com",
					},
				},
			},
			expected: scim.ResourceAttributes{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "home",
						"value": "ivixvi@example.com",
					},
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi-added@example.com",
					},
				},
			},
			expectedChanged: true,
		},
		{
			name: "Add operation - Core MultiValued Complex Attributes - Add All no changed",
			op: scim.PatchOperation{
				Op:   "add",
				Path: path(`emails`),
				Value: []interface{}{
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi-added@example.com",
					},
				},
			},
			data: scim.ResourceAttributes{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "home",
						"value": "ivixvi@example.com",
					},
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi-added@example.com",
					},
				},
			},
			expected: scim.ResourceAttributes{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "home",
						"value": "ivixvi@example.com",
					},
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi-added@example.com",
					},
				},
			},
			expectedChanged: false,
		},
		{
			name: "Add operation - Core MultiValued Complex Attributes - Add For Item",
			op: scim.PatchOperation{
				Op:   "add",
				Path: path(`emails[type eq "work"]`),
				Value: map[string]interface{}{
					"type":    "work",
					"value":   "ivixvi-updated@example.com",
					"primary": true,
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
						"type":    "work",
						"value":   "ivixvi-updated@example.com",
						"primary": true,
					},
				},
			},
			expectedChanged: true,
		},
		{
			name: "Add operation - Core MultiValued Complex Attributes - Add For Item no changed",
			op: scim.PatchOperation{
				Op:   "add",
				Path: path(`emails[type eq "work"]`),
				Value: map[string]interface{}{
					"type":  "work",
					"value": "ivixvi@example.com",
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
						"value": "ivixvi@example.com",
					},
				},
			},
			expectedChanged: false,
		},
		{
			name: "Add operation - Core MultiValued Complex Attributes - Replace For Attribute",
			op: scim.PatchOperation{
				Op:    "add",
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
		{
			name: "Add operation - Core MultiValued Complex Attributes - Replace For Attribute no changed",
			op: scim.PatchOperation{
				Op:    "add",
				Path:  path(`emails[type eq "work"].value`),
				Value: "ivixvi@example.com",
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
						"value": "ivixvi@example.com",
					},
				},
			},
			expectedChanged: false,
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

package scimpatch_test

import (
	"fmt"
	"testing"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/schema"
	scimpatch "github.com/ivixvi/scim-patch"
)

// TestPathSpecifiedReplace は Patcher.Apply の path指定をしているreplace操作の正常系をテストします
func TestPathSpecifiedReplace(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name            string
		op              scim.PatchOperation
		data            map[string]interface{}
		expected        map[string]interface{}
		expectedChanged bool
	}{
		// Common Attribute
		{
			name: "Replace operation - common Attributes - externalId",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`externalId`),
				Value: "galice",
			},
			data: map[string]interface{}{
				"externalId": "gbob",
			},
			expected: map[string]interface{}{
				"externalId": "galice",
			},
			expectedChanged: true,
		},
		// Singular Attribute
		{
			name: "Replace operation - Core Singular Attributes - replace",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`displayName`),
				Value: "Alice Green",
			},
			data: map[string]interface{}{
				"displayName": "Bob Green",
			},
			expected: map[string]interface{}{
				"displayName": "Alice Green",
			},
			expectedChanged: true,
		},
		{
			name: "Replace operation - Core Singular Attributes - no value",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`displayName`),
				Value: "Alice Green",
			},
			data: map[string]interface{}{},
			expected: map[string]interface{}{
				"displayName": "Alice Green",
			},
			expectedChanged: true,
		},
		{
			name: "Replace operation - Core Singular Attributes - no changed",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`displayName`),
				Value: "Alice Green",
			},
			data: map[string]interface{}{
				"displayName": "Alice Green",
			},
			expected: map[string]interface{}{
				"displayName": "Alice Green",
			},
			expectedChanged: false,
		},
		// Complex Attribute
		{
			name: "Replace operation - Core Complex Attributes - replace",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`name.familyName`),
				Value: "Green",
			},
			data: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Blue",
					"givenName":  "Alice",
				},
			},
			expected: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expectedChanged: true,
		},
		{
			name: "Replace operation - Core Complex Attributes - no value",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`name.familyName`),
				Value: "Green",
			},
			data: map[string]interface{}{
				"name": map[string]interface{}{
					"givenName": "Alice",
				},
			},
			expected: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expectedChanged: true,
		},
		{
			name: "Replace operation - Core Complex Attributes - no changed.",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`name.familyName`),
				Value: "Green",
			},
			data: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expected: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expectedChanged: false,
		},
		{
			name: "Replace operation - Core Complex Attributes - map specified replace",
			op: scim.PatchOperation{
				Op:   "replace",
				Path: path(`name`),
				Value: map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			data: map[string]interface{}{
				"name": map[string]interface{}{
					"givenName": "Bob",
				},
			},
			expected: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
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
			data: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expected: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expectedChanged: false,
		},
		// MultiValued Complex Attribute
		{
			name: "Replace operation - Core MultiValued Complex Attributes - Replace All",
			op: scim.PatchOperation{
				Op:   "replace",
				Path: path(`emails`),
				Value: []interface{}{
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi-updated@example.com",
					},
				},
			},
			data: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "home",
						"value": "ivixvi@example.com",
					},
				},
			},
			expected: map[string]interface{}{
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
			name: "Replace operation - Core MultiValued Complex Attributes - Replace All no changed",
			op: scim.PatchOperation{
				Op:   "replace",
				Path: path(`emails`),
				Value: []interface{}{
					map[string]interface{}{
						"type":  "home",
						"value": "ivixvi@example.com",
					},
				},
			},
			data: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "home",
						"value": "ivixvi@example.com",
					},
				},
			},
			expected: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "home",
						"value": "ivixvi@example.com",
					},
				},
			},
			expectedChanged: false,
		},
		{
			name: "Replace operation - Core MultiValued Complex Attributes - Replace For Item",
			op: scim.PatchOperation{
				Op:   "replace",
				Path: path(`emails[type eq "work"]`),
				Value: map[string]interface{}{
					"type":  "work",
					"value": "ivixvi-updated@example.com",
				},
			},
			data: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"type":    "work",
						"value":   "ivixvi@example.com",
						"primary": true,
					},
				},
			},
			expected: map[string]interface{}{
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
			name: "Replace operation - Core MultiValued Complex Attributes - Replace For Item no changed",
			op: scim.PatchOperation{
				Op:   "replace",
				Path: path(`emails[type eq "work"]`),
				Value: map[string]interface{}{
					"type":  "work",
					"value": "ivixvi@example.com",
				},
			},
			data: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi@example.com",
					},
				},
			},
			expected: map[string]interface{}{
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
			name: "Replace operation - Core MultiValued Complex Attributes - Replace For Attribute",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`emails[type eq "work"].value`),
				Value: "ivixvi-updated@example.com",
			},
			data: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi@example.com",
					},
				},
			},
			expected: map[string]interface{}{
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
			name: "Replace operation - Core MultiValued Complex Attributes - Replace For Attribute no changed",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`emails[type eq "work"].value`),
				Value: "ivixvi@example.com",
			},
			data: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"type":  "work",
						"value": "ivixvi@example.com",
					},
				},
			},
			expected: map[string]interface{}{
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
			name: "Replace operation - MultiValued Complex Attributes - Filter & Value addition",
			op: scim.PatchOperation{
				Op:    "replace",
				Path:  path(`emails[type eq "work"].value`),
				Value: "ivixvi@example.com",
			},
			data:            map[string]interface{}{},
			expected:        map[string]interface{}{},
			expectedChanged: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			// Create a Patcher instance with a dummy schema
			patcher := scimpatch.NewPatcher(schema.CoreUserSchema(), []schema.Schema{schema.ExtensionEnterpriseUser()}, nil)

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

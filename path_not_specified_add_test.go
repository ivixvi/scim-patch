package scimpatch_test

import (
	"fmt"
	"testing"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/schema"
	scimpatch "github.com/ivixvi/scim-patch"
)

// TestPathNotspecifiedAdd は Pacher.Apply の path指定をしていいないadd操作の正常系をテストします
func TestPathNotspecifiedAdd(t *testing.T) {
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
				Op: "add",
				Value: map[string]interface{}{
					"displayName": "Alice Green",
				},
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
				Op: "add",
				Value: map[string]interface{}{
					"displayName": "Alice Green",
				},
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
				Op: "add",
				Value: map[string]interface{}{
					"displayName": "Alice Green",
				},
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
				Op: "add",
				Value: map[string]interface{}{
					"name": map[string]interface{}{
						"familyName": "Green",
					},
				},
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
				Op: "add",
				Value: map[string]interface{}{
					"name": map[string]interface{}{
						"familyName": "Green",
					},
				},
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
				Op: "add",
				Value: map[string]interface{}{
					"name": map[string]interface{}{
						"familyName": "Green",
					},
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
		{
			name: "Add operation - Core Complex Attributes - map specified merge",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"name": map[string]interface{}{
						"givenName": "Alice",
						"formatted": "Alice Green",
					},
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
			name: "Add operation - Core Singular Attributes - map specified no changed",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"name": map[string]interface{}{
						"familyName": "Green",
						"givenName":  "Alice",
					},
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
				Op: "add",
				Value: map[string]interface{}{
					"emails": []interface{}{
						map[string]interface{}{
							"type":  "work",
							"value": "ivixvi-added@example.com",
						},
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
				Op: "add",
				Value: map[string]interface{}{
					"emails": []interface{}{
						map[string]interface{}{
							"type":  "work",
							"value": "ivixvi-added@example.com",
						},
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
			name: "Add operation - Core MultiValued Complex Attributes - Add For Item no changed",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"emails": []interface{}{
						map[string]interface{}{
							"type":  "work",
							"value": "ivixvi@example.com",
						},
					},
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
			name: "Add operation - Core MultiValued Complex Attributes - Replace For Attribute no changed",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"emails": []interface{}{
						map[string]interface{}{
							"type":  "work",
							"value": "ivixvi@example.com",
						},
					},
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
			name: "Add operation - MultiValued Complex Attributes - Filter & Value addition",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"emails": []interface{}{
						map[string]interface{}{
							"type":  "work",
							"value": "ivixvi@example.com",
						},
					},
				},
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
		{
			name: "Add operation - Extention Singular Attribute - URI Prefix not exists.",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
						"department": "department1",
					},
				},
			},
			data: scim.ResourceAttributes{},
			expected: scim.ResourceAttributes{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"department": "department1",
				},
			},
			expectedChanged: true,
		},
		{
			name: "Add operation - Extention Singular Attribute - URI Prefix exists.",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
						"department": "updated-department",
					},
				},
			},
			data: scim.ResourceAttributes{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"department": "department",
				},
			},
			expected: scim.ResourceAttributes{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"department": "updated-department",
				},
			},
			expectedChanged: true,
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

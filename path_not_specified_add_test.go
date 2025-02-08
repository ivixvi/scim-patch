package scimpatch_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/schema"
	scimpatch "github.com/ivixvi/scim-patch"
)

// TestPathNotSpecifiedAdd は Patcher.Apply の path指定をしていない add 操作の正常系をテストします
func TestPathNotSpecifiedAdd(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name            string
		op              scim.PatchOperation
		data            map[string]interface{}
		expected        map[string]interface{}
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
			data: map[string]interface{}{},
			expected: map[string]interface{}{
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
			data: map[string]interface{}{
				"displayName": "Bob Green",
			},
			expected: map[string]interface{}{
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
			name: "Add operation - Core Complex Attributes - replace",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"name": map[string]interface{}{
						"familyName": "Green",
					},
				},
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
			name: "Add operation - Core Complex Attributes - no value",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"name": map[string]interface{}{
						"familyName": "Green",
					},
				},
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
			name: "Add operation - Core Complex Attributes - no changed.",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"name": map[string]interface{}{
						"familyName": "Green",
					},
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
			data: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Green",
					"formatted":  "Bob Green",
				},
			},
			expected: map[string]interface{}{
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
			data: map[string]interface{}{
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
			expected: map[string]interface{}{
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
			data: map[string]interface{}{},
			expected: map[string]interface{}{
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
			name: "Add operation - Extension Singular Attribute - URI Prefix not exists.",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
						"department": "department1",
					},
				},
			},
			data: map[string]interface{}{},
			expected: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"department": "department1",
				},
			},
			expectedChanged: true,
		},
		{
			name: "Add operation - Extension Singular Attribute - URI Prefix exists.",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
						"department": "updated-department",
					},
				},
			},
			data: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"department": "department",
				},
			},
			expected: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"department": "updated-department",
				},
			},
			expectedChanged: true,
		},
		// MultiValued Attribute
		{
			name: "Add operation - Extension MultiValued Attributes - add",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"urn:ivixvi:testSchema": map[string]interface{}{
						"testString": []interface{}{"value"},
					},
				},
			},
			data: map[string]interface{}{},
			expected: map[string]interface{}{
				"urn:ivixvi:testSchema": map[string]interface{}{
					"testString": []interface{}{"value"},
				},
			},
			expectedChanged: true,
		},
		{
			name: "Add operation - Extension MultiValued Attributes - merge slice",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"urn:ivixvi:testSchema": map[string]interface{}{
						"testString": []interface{}{"newValue"},
					},
				},
			},
			data: map[string]interface{}{
				"urn:ivixvi:testSchema": map[string]interface{}{
					"testString": []interface{}{"oldValue"},
				},
			},
			expected: map[string]interface{}{
				"urn:ivixvi:testSchema": map[string]interface{}{
					"testString": []interface{}{"oldValue", "newValue"},
				},
			},
			expectedChanged: true,
		},
		{
			name: "Add operation - Extension MultiValued Attributes - no changed",
			op: scim.PatchOperation{
				Op: "add",
				Value: map[string]interface{}{
					"urn:ivixvi:testSchema": map[string]interface{}{
						"testString": []interface{}{"value"},
					},
				},
			},
			data: map[string]interface{}{
				"urn:ivixvi:testSchema": map[string]interface{}{
					"testString": []interface{}{"value"},
				},
			},
			expected: map[string]interface{}{
				"urn:ivixvi:testSchema": map[string]interface{}{
					"testString": []interface{}{"value"},
				},
			},
			expectedChanged: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			// Create a Patcher instance with a dummy schema
			patcher := scimpatch.NewPatcher(schema.CoreUserSchema(), []schema.Schema{
				schema.ExtensionEnterpriseUser(),
				TestExtensionSchema,
			}, nil)

			// Apply the PatchOperation
			result, changed, err := patcher.Apply(context.TODO(),tc.op, tc.data)
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

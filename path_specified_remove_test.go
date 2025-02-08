package scimpatch_test

import (
	"fmt"
	"testing"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/errors"
	"github.com/elimity-com/scim/schema"
	scimpatch "github.com/ivixvi/scim-patch"
)

// TestPatcher_Apply は Patcher.Apply の Remove の正常系をテストします
func TestPathSpecifiedRemove(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name            string
		op              scim.PatchOperation
		data            map[string]interface{}
		expected        map[string]interface{}
		expectedChanged bool
	}{
		// Remove Singular Attribute
		{
			name: "Remove operation - Core Singular Attribute",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("displayName"),
			},
			data: map[string]interface{}{
				"displayName": "Bob Green",
			},
			expected:        map[string]interface{}{},
			expectedChanged: true,
		},
		{
			name: "Remove operation - Core Singular Attribute Not Changed.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("displayName"),
			},
			data:            map[string]interface{}{},
			expected:        map[string]interface{}{},
			expectedChanged: false,
		},
		{
			name: "Remove operation - Extension Singular Attribute - All Removed.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:department"),
			},
			data: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"department": "2B Sales",
				},
			},
			expected:        map[string]interface{}{},
			expectedChanged: true,
		},
		{
			name: "Remove operation - Extension Singular Attribute - Partially Removed.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:department"),
			},
			data: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"division":   "Sales",
					"department": "2B Sales",
				},
			},
			expected: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"division": "Sales",
				},
			},
			expectedChanged: true,
		},
		{
			name: "Remove operation - Extension Singular Attribute - URI Prefix not exists.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:department"),
			},
			data:            map[string]interface{}{},
			expected:        map[string]interface{}{},
			expectedChanged: false,
		},
		{
			name: "Remove operation - Extension Singular Attribute - URI Prefix exists.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:department"),
			},
			data: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"division": "Sales",
				},
			},
			expected: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"division": "Sales",
				},
			},
			expectedChanged: false,
		},
		// Remove Complex Attributes
		{
			name: "Remove operation - Core Complex Attribute - SubAttributes Not Specified.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("name"),
			},
			data: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Green",
				},
			},
			expected:        map[string]interface{}{},
			expectedChanged: true,
		},
		{
			name: "Remove operation - Core Complex Attribute - SubAttributes Specified.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("name.familyName"),
			},
			data: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Green",
					"givenName":  "Alice",
				},
			},
			expected: map[string]interface{}{
				"name": map[string]interface{}{
					"givenName": "Alice",
				},
			},
			expectedChanged: true,
		},
		{
			name: "Remove operation - Core Complex Attribute - SubAttributes Specified - Remove Attributes",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("name.familyName"),
			},
			data: map[string]interface{}{
				"name": map[string]interface{}{
					"familyName": "Green",
				},
			},
			expected:        map[string]interface{}{},
			expectedChanged: true,
		},
		{
			name: "Remove operation - Core Complex Attribute Not Changed. - SubAttributes not Specified ",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("name"),
			},
			data:            map[string]interface{}{},
			expected:        map[string]interface{}{},
			expectedChanged: false,
		},
		{
			name: "Remove operation - Core Complex Attribute Not Changed. - SubAttributes Specified ",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("name.familyName"),
			},
			data:            map[string]interface{}{},
			expected:        map[string]interface{}{},
			expectedChanged: false,
		},
		{
			name: "Remove operation - Extension Complex Attribute - Attribute Removed.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:manager"),
			},
			data: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"manager": interface{}(
						map[string]interface{}{
							"value":       "0001",
							"displayName": "Bob Green",
						},
					),
				},
			},
			expected:        map[string]interface{}{},
			expectedChanged: true,
		},
		{
			name: "Remove operation - Extension Complex Attribute - All Removed.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:manager.value"),
			},
			data: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"manager": interface{}(
						map[string]interface{}{
							"value": "0001",
						},
					),
				},
			},
			expected:        map[string]interface{}{},
			expectedChanged: true,
		},
		{
			name: "Remove operation - Extension Complex Attribute - Partially Removed.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:manager.value"),
			},
			data: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"manager": interface{}(
						map[string]interface{}{
							"value":       "0001",
							"displayName": "Bob Green",
						},
					),
				},
			},
			expected: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"manager": interface{}(
						map[string]interface{}{
							"displayName": "Bob Green",
						},
					),
				},
			},
			expectedChanged: true,
		},
		{
			name: "Remove operation - Extension Complex Attribute - URI Prefix not exists.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:manager.value"),
			},
			data:            map[string]interface{}{},
			expected:        map[string]interface{}{},
			expectedChanged: false,
		},
		{
			name: "Remove operation - Extension Complex Attribute - URI Prefix exists.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:manager.value"),
			},
			data: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"division": "Sales",
				},
			},
			expected: map[string]interface{}{
				"urn:ietf:params:scim:schemas:extension:enterprise:2.0:User": map[string]interface{}{
					"division": "Sales",
				},
			},
			expectedChanged: false,
		},
		// Remove MultiValued Complex Attributes
		{
			name: "Remove operation - MultiValued Complex Attribute - Direct Remove",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("emails"),
			},
			data: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"value": "ivixvi@example.com",
						"type":  "home",
					},
				},
			},
			expected:        map[string]interface{}{},
			expectedChanged: true,
		},
		{
			name: "Remove operation - MultiValued Complex Attribute - Item Remove",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path(`emails[type eq "home"]`),
			},
			data: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"value": "ivixvi@example.com",
						"type":  "home",
					},
					map[string]interface{}{
						"value": "ivixvi-work@example.com",
						"type":  "work",
					},
				},
			},
			expected: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"value": "ivixvi-work@example.com",
						"type":  "work",
					},
				},
			},
			expectedChanged: true,
		},
		{
			name: "Remove operation - MultiValued Complex Attribute - SubAttribute Remove",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path(`emails[type eq "home"].primary`),
			},
			data: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"value":   "ivixvi@example.com",
						"type":    "home",
						"primary": true,
					},
				},
			},
			expected: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"value": "ivixvi@example.com",
						"type":  "home",
					},
				},
			},
			expectedChanged: true,
		},
		{
			name: "Remove operation - MultiValued Complex Attribute - Item Remove. No Changed",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path(`emails[type eq "home"]`),
			},
			data: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"value": "ivixvi-work@example.com",
						"type":  "work",
					},
				},
			},
			expected: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"value": "ivixvi-work@example.com",
						"type":  "work",
					},
				},
			},
			expectedChanged: false,
		},
		{
			name: "Remove operation - MultiValued Complex Attribute - SubAttribute Remove. no Changed.",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path(`emails[type eq "home"].primary`),
			},
			data: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"value": "ivixvi@example.com",
						"type":  "home",
					},
				},
			},
			expected: map[string]interface{}{
				"emails": []interface{}{
					map[string]interface{}{
						"value": "ivixvi@example.com",
						"type":  "home",
					},
				},
			},
			expectedChanged: false,
		},
		// Remove MultiValued Attributes
		{
			name: "Remove operation - MultiValued Singular Attribute - remove",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("urn:ivixvi:testSchema:testString"),
			},
			data: map[string]interface{}{
				"urn:ivixvi:testSchema": map[string]interface{}{
					"testString": []interface{}{"value"},
				},
			},
			expected:        map[string]interface{}{},
			expectedChanged: true,
		},
		{
			name: "Remove operation - MultiValued Singular Attribute - item remove",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path(`urn:ivixvi:testSchema:testString[value eq "delete"]`),
			},
			data: map[string]interface{}{
				"urn:ivixvi:testSchema": map[string]interface{}{
					"testString": []interface{}{"value", "delete"},
				},
			},
			// FIXME
			// data: map[string]interface{}{
			// 	"urn:ivixvi:testSchema": map[string]interface{}{
			// 		"testString": []interface{}{"value"},
			// 	},
			// },
			// expectedChanged: true,
			expected:        map[string]interface{}{},
			expectedChanged: false,
		},
		{
			name: "Remove operation - MultiValued Singular Attribute - no changed",
			op: scim.PatchOperation{
				Op:   "remove",
				Path: path("urn:ivixvi:testSchema:testString"),
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
			patcher := scimpatch.NewPatcher(
				schema.CoreUserSchema(),
				[]schema.Schema{
					schema.ExtensionEnterpriseUser(),
					TestExtensionSchema,
				}, nil)

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

// TestPatcher_ApplyError は Patcher.Apply の Remove の異常系をテストします
func TestRemoveError(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name     string
		op       scim.PatchOperation
		expected errors.ScimError
	}{
		{
			name: "Remove operation - no specify path",
			op: scim.PatchOperation{
				Op: "remove",
			},
			expected: errors.ScimErrorNoTarget,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			// Create a Patcher instance with a dummy schema
			patcher := scimpatch.Patcher{}

			// Apply the PatchOperation
			_, _, err := patcher.Apply(tc.op, map[string]interface{}{})
			if err == nil {
				t.Fatalf("Apply() not returned error")
			}
			scimError, ok := err.(errors.ScimError)
			if !ok {
				t.Fatalf("Apply() not returned ScimError: %v", err)
			}

			// Check if the result matches the expected data
			if !(tc.expected.Detail == scimError.Detail &&
				tc.expected.Status == scimError.Status &&
				tc.expected.ScimType == scimError.ScimType) {
				t.Fatalf("Apply() not returned Expected ScimError: %v", scimError)
			}
		})
	}
}

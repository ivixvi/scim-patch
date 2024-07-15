package scimpatch_test

import (
	"fmt"
	"testing"

	scimpatch "github.com/ivixvi/scim-patch"
)

// TestAreEveryItemsMap は areEveryItemsMap をテストします
func TestAreEveryItemsMap(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name          string
		target        []interface{}
		expectedOk    bool
		expectedSlice []map[string]interface{}
	}{
		{
			name:          "empty",
			target:        []interface{}{},
			expectedOk:    true,
			expectedSlice: []map[string]interface{}{},
		},
		{
			name: "success",
			target: []interface{}{
				map[string]interface{}{},
				map[string]interface{}{},
			},
			expectedOk: true,
			expectedSlice: []map[string]interface{}{
				map[string]interface{}{},
				map[string]interface{}{},
			},
		},
		{
			name: "failed",
			target: []interface{}{
				"hoge", 31,
			},
			expectedOk:    false,
			expectedSlice: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			// call AreEveryItemsMap
			actual, ok := scimpatch.AreEveryItemsMap(tc.target)

			// Check if the result matches the expected data
			if !(tc.expectedOk == ok &&
				fmt.Sprint(tc.expectedSlice) == fmt.Sprint(actual)) {
				t.Fatalf("AreEveryItemsMap() not returned Expected: %v, %v", actual, ok)
			}
		})
	}
}

// TestEqMap は eqMap をテストします
func TestEqMap(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name     string
		m1       map[string]interface{}
		m2       map[string]interface{}
		expected bool
	}{
		{
			name:     "empty",
			m1:       map[string]interface{}{},
			m2:       map[string]interface{}{},
			expected: true,
		},
		{
			name: "match length",
			m1: map[string]interface{}{
				"a2": "string",
			},
			m2: map[string]interface{}{
				"a1": "string",
			},
			expected: false,
		},
		{
			name: "contains m1",
			m1: map[string]interface{}{
				"a1": "string",
			},
			m2: map[string]interface{}{
				"a1": "string",
				"a2": "string",
			},
			expected: false,
		},
		{
			name: "contains m2",
			m1: map[string]interface{}{
				"a1": "string",
				"a2": "string",
			},
			m2: map[string]interface{}{
				"a1": "string",
			},
			expected: false,
		},
		{
			name: "match",
			m1: map[string]interface{}{
				"a1": "string",
				"a2": "string",
			},
			m2: map[string]interface{}{
				"a1": "string",
				"a2": "string",
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			// call AreEveryItemsMap
			ok := scimpatch.EqMap(tc.m1, tc.m2)

			// Check if the result matches the expected data
			if !(tc.expected == ok) {
				t.Fatalf("EqMap() not returned Expected: %v, %v", tc.expected, ok)
			}
		})
	}
}

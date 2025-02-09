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
				{},
				{},
			},
		},
		{
			name: "failed",
			target: []interface{}{
				"value", 31,
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

// TestMergeMap は mergeMap をテストします
func TestMergeMap(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name        string
		m1          map[string]interface{}
		m2          map[string]interface{}
		expectedMap map[string]interface{}
		expectedOk  bool
	}{
		{
			name:        "empty",
			m1:          map[string]interface{}{},
			m2:          map[string]interface{}{},
			expectedMap: map[string]interface{}{},
			expectedOk:  false,
		},
		{
			name: "merger contains mergee",
			m1: map[string]interface{}{
				"a1": "string",
			},
			m2: map[string]interface{}{
				"a1": "string",
				"a2": "string",
			},
			expectedMap: map[string]interface{}{
				"a1": "string",
				"a2": "string",
			},
			expectedOk: true,
		},
		{
			name: "mergee contains merger",
			m1: map[string]interface{}{
				"a1": "string",
				"a2": "string",
			},
			m2: map[string]interface{}{
				"a1": "string",
			},
			expectedMap: map[string]interface{}{
				"a1": "string",
				"a2": "string",
			},
			expectedOk: false,
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
			expectedMap: map[string]interface{}{
				"a1": "string",
				"a2": "string",
			},
			expectedOk: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			// call AreEveryItemsMap
			actual, ok := scimpatch.MergeMap(tc.m1, tc.m2)

			// Check if the result matches the expected data
			if !(tc.expectedOk == ok && fmt.Sprint(tc.expectedMap) == fmt.Sprint(actual)) {
				t.Fatalf("EqMap() not returned Expected ok: %v, %v", tc.expectedOk, ok)
				t.Fatalf("EqMap() not returned Expected map: %v, %v", tc.expectedMap, actual)
			}
		})
	}
}

// TestContainsMap は containsMap をテストします
func TestContainsMap(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name     string
		slice    []map[string]interface{}
		item     map[string]interface{}
		expected bool
	}{
		{
			name:     "empty slice",
			slice:    []map[string]interface{}{},
			item:     map[string]interface{}{"key": "value"},
			expected: false,
		},
		{
			name: "item not in slice",
			slice: []map[string]interface{}{
				{"key1": "value1"},
				{"key2": "value2"},
			},
			item:     map[string]interface{}{"key": "value"},
			expected: false,
		},
		{
			name: "item in slice",
			slice: []map[string]interface{}{
				{"key1": "value1"},
				{"key": "value"},
			},
			item:     map[string]interface{}{"key": "value"},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			// call ContainsMap
			ok := scimpatch.ContainsMap(tc.slice, tc.item)

			// Check if the result matches the expected data
			if !(tc.expected == ok) {
				t.Fatalf("ContainsMap() not returned Expected: %v, %v", tc.expected, ok)
			}
		})
	}
}

// TestContainsItem は containsItem をテストします
func TestContainsItem(t *testing.T) {
	// Define the test cases
	testCases := []struct {
		name     string
		slice    []interface{}
		item     interface{}
		expected bool
	}{
		{
			name:     "empty slice",
			slice:    []interface{}{},
			item:     "value",
			expected: false,
		},
		{
			name: "item not in slice",
			slice: []interface{}{
				"value1",
				"value2",
			},
			item:     "value",
			expected: false,
		},
		{
			name: "item in slice",
			slice: []interface{}{
				"value1",
				"value",
			},
			item:     "value",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Log(tc.name)
			// call ContainsItem
			ok := scimpatch.ContainsItem(tc.slice, tc.item)

			// Check if the result matches the expected data
			if !(tc.expected == ok) {
				t.Fatalf("ContainsItem() not returned Expected: %v, %v", tc.expected, ok)
			}
		})
	}
}

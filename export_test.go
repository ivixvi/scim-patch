// export_test is export function, variable, etc... for testing
package scimpatch

// AreEveryItemsMap is export areEveryItemsMap for testing
func AreEveryItemsMap(s interface{}) ([]map[string]interface{}, bool) {
	return areEveryItemsMap(s)
}

// EqMap is export eqMap for testing
func EqMap(m1 map[string]interface{}, m2 map[string]interface{}) bool {
	return eqMap(m1, m2)
}

// MergeMap is export mergeMap for testing
func MergeMap(m1 map[string]interface{}, m2 map[string]interface{}) (map[string]interface{}, bool) {
	return mergeMap(m1, m2)
}

// ContainsMap is export containsMap for testing
func ContainsMap(slice []map[string]interface{}, item map[string]interface{}) bool {
	return containsMap(slice, item)
}

// ContainsItem is export containsItem for testing
func ContainsItem(slice []interface{}, item interface{}) bool {
	return containsItem(slice, item)
}

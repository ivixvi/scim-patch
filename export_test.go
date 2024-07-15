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

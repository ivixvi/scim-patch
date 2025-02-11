package scimpatch

import (
	"github.com/scim2/filter-parser/v2"
)

type adder struct{}

var adderInstance *adder

func (r *adder) Direct(scopedMap map[string]interface{}, scopedAttr string, value interface{}) bool {
	switch newValue := value.(type) {
	case []map[string]interface{}:
		return r.addMapSlice(scopedMap, scopedAttr, newValue)
	case map[string]interface{}:
		return r.addMap(scopedMap, scopedAttr, newValue)
	case []interface{}:
		return r.addSlice(scopedMap, scopedAttr, newValue)
	case interface{}:
		return r.addValue(scopedMap, scopedAttr, newValue)
	}
	return false
}

func (r *adder) addMapSlice(scopedMap map[string]interface{}, scopedAttr string, newValue []map[string]interface{}) bool {
	oldSlice, ok := scopedMap[scopedAttr]
	if !ok {
		scopedMap[scopedAttr] = newValue
		return true
	}
	oldMaps, ok := areEveryItemsMap(oldSlice)
	if !ok {
		scopedMap[scopedAttr] = newValue
		return true
	}
	changed := false
	for _, newMap := range newValue {
		if !containsMap(oldMaps, newMap) {
			oldMaps = append(oldMaps, newMap)
			changed = true
		}
	}
	if changed {
		scopedMap[scopedAttr] = oldMaps
	}
	return changed
}

func (r *adder) addMap(scopedMap map[string]interface{}, scopedAttr string, newValue map[string]interface{}) bool {
	oldMap, ok := scopedMap[scopedAttr].(map[string]interface{})
	if ok {
		changed := false
		scopedMap[scopedAttr], changed = mergeMap(oldMap, newValue)
		return changed
	}
	scopedMap[scopedAttr] = newValue
	return true
}

func (r *adder) addSlice(scopedMap map[string]interface{}, scopedAttr string, newValue []interface{}) bool {
	oldSlice, ok := scopedMap[scopedAttr].([]interface{})
	// oldSlice is nil
	if !ok {
		scopedMap[scopedAttr] = newValue
		return true
	}

	// Complex MultiValued
	if newMaps, ok := areEveryItemsMap(newValue); ok {
		return r.addMapSlice(scopedMap, scopedAttr, newMaps)
	}

	// Singular MultiValued
	changed := false
	for _, newItem := range newValue {
		if !containsItem(oldSlice, newItem) {
			oldSlice = append(oldSlice, newItem)
			changed = true
		}
	}
	if changed {
		scopedMap[scopedAttr] = oldSlice
	}
	return changed
}

func (r *adder) addValue(scopedMap map[string]interface{}, scopedAttr string, newValue interface{}) bool {
	if oldValue, ok := scopedMap[scopedAttr]; !ok || oldValue != newValue {
		scopedMap[scopedAttr] = newValue
		return true
	}
	return false
}

func (r *adder) ByValueForItem(scopedSlice []interface{}, value interface{}) ([]interface{}, bool) {
	changed := false
	if !containsItem(scopedSlice, value) {
		changed = true
		scopedSlice = append(scopedSlice, value)
	}
	return scopedSlice, changed
}

func (r *adder) ByValueExpressionForItem(scopedMaps []map[string]interface{}, expr filter.Expression, value interface{}) ([]map[string]interface{}, bool) {
	newValue, ok := value.(map[string]interface{})

	// unexpected input
	if !ok {
		return scopedMaps, false
	}

	changed := false
	for i, oldValue := range scopedMaps {
		if isMatchExpression(oldValue, expr) && !eqMap(oldValue, newValue) {
			var merger map[string]interface{}
			merger, changed = mergeMap(oldValue, newValue)
			scopedMaps[i] = merger
		}
	}
	return scopedMaps, changed
}

func (r *adder) ByValueExpressionForAttribute(scopedMaps []map[string]interface{}, expr filter.Expression, subAttr string, value interface{}) ([]map[string]interface{}, bool) {
	scopedMaps, changed, found := replaceByValueExpressionForAttribute(scopedMaps, expr, subAttr, value)
	if !found {
		changed = true
		newMap := toMap(expr)
		newMap[subAttr] = value
		scopedMaps = append(scopedMaps, newMap)
	}
	return scopedMaps, changed
}

package scimpatch

import (
	"github.com/scim2/filter-parser/v2"
)

type replacer struct{}

var replacerInstance *replacer

func (r *replacer) Direct(scopedMap map[string]interface{}, scopedAttr string, value interface{}) bool {
	switch newValue := value.(type) {
	case []map[string]interface{}:
		return r.replaceMapSlice(scopedMap, scopedAttr, newValue)
	case map[string]interface{}:
		return r.replaceMap(scopedMap, scopedAttr, newValue)
	case []interface{}:
		return r.replaceSlice(scopedMap, scopedAttr, newValue)
	case interface{}:
		return r.replaceValue(scopedMap, scopedAttr, newValue)
	}
	return false
}

func (r *replacer) replaceMapSlice(scopedMap map[string]interface{}, scopedAttr string, newValue []map[string]interface{}) bool {
	oldSlice, ok := scopedMap[scopedAttr]
	if !ok {
		scopedMap[scopedAttr] = newValue
		return true
	}
	oldMaps, ok := areEveryItemsMap(oldSlice)
	if !ok || len(oldMaps) != len(newValue) {
		scopedMap[scopedAttr] = newValue
		return true
	}
	for _, newMap := range newValue {
		if !containsMap(oldMaps, newMap) {
			scopedMap[scopedAttr] = newValue
			return true
		}
	}
	return false
}

func (r *replacer) replaceMap(scopedMap map[string]interface{}, scopedAttr string, newValue map[string]interface{}) bool {
	oldMap, ok := scopedMap[scopedAttr].(map[string]interface{})
	if ok && eqMap(newValue, oldMap) {
		return false
	}
	scopedMap[scopedAttr] = newValue
	return true
}

func (r *replacer) replaceSlice(scopedMap map[string]interface{}, scopedAttr string, newValue []interface{}) bool {
	oldSlice, ok := scopedMap[scopedAttr].([]interface{})
	// oldSlice is nil
	if !ok || len(oldSlice) != len(newValue) {
		scopedMap[scopedAttr] = newValue
		return true
	}

	// Complex MultiValued
	if newMaps, ok := areEveryItemsMap(newValue); ok {
		return r.replaceMapSlice(scopedMap, scopedAttr, newMaps)
	}

	// Singular MultiValued
	for _, newItem := range newValue {
		if !containsItem(oldSlice, newItem) {
			scopedMap[scopedAttr] = newValue
			return true
		}
	}

	return false
}

func (r *replacer) replaceValue(scopedMap map[string]interface{}, scopedAttr string, newValue interface{}) bool {
	if oldValue, ok := scopedMap[scopedAttr]; !ok || oldValue != newValue {
		scopedMap[scopedAttr] = newValue
		return true
	}
	return false
}

func (r *replacer) ByValueForItem(scopedSlice []interface{}, value interface{}) ([]interface{}, bool) {
	changed := false
	if !containsItem(scopedSlice, value) {
		changed = true
		scopedSlice = append(scopedSlice, value)
	}
	return scopedSlice, changed
}

func (r *replacer) ByValueExpressionForItem(scopedMaps []map[string]interface{}, expr filter.Expression, value interface{}) ([]map[string]interface{}, bool) {
	newValue, ok := value.(map[string]interface{})

	// unexpected input
	if !ok {
		return scopedMaps, false
	}

	changed := false
	for i, oldValue := range scopedMaps {
		if isMatchExpression(oldValue, expr) && !eqMap(oldValue, newValue) {
			changed = true
			scopedMaps[i] = newValue
		}
	}
	return scopedMaps, changed
}

func (r *replacer) ByValueExpressionForAttribute(scopedMaps []map[string]interface{}, expr filter.Expression, subAttr string, value interface{}) ([]map[string]interface{}, bool) {
	scopedMaps, changed, _ := replaceByValueExpressionForAttribute(scopedMaps, expr, subAttr, value)
	return scopedMaps, changed
}

func replaceByValueExpressionForAttribute(
	scopedMaps []map[string]interface{},
	expr filter.Expression,
	subAttr string,
	value interface{},
) ([]map[string]interface{}, bool, bool) {
	changed := false
	found := false
	for _, oldValue := range scopedMaps {
		if isMatchExpression(oldValue, expr) {
			found = true
			oldAttrValue, ok := oldValue[subAttr]
			if !ok || oldAttrValue != value {
				changed = true
				oldValue[subAttr] = value
			}
		}
	}
	return scopedMaps, changed, found
}

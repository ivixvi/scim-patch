package scimpatch

import (
	"fmt"

	"github.com/scim2/filter-parser/v2"
)

type Replacer struct{}

var replacer *Replacer

func (r *Replacer) Direct(scopedMap map[string]interface{}, scopedAttr string, value interface{}) (map[string]interface{}, bool) {
	switch newValue := value.(type) {
	case []map[string]interface{}:
	case map[string]interface{}:
		oldMap, ok := scopedMap[scopedAttr].(map[string]interface{})
		if ok && eqMap(newValue, oldMap) {
			return scopedMap, false
		}
		scopedMap[scopedAttr] = value
		return scopedMap, true
	case []interface{}:
	case interface{}:
		if oldValue, ok := scopedMap[scopedAttr]; !ok || oldValue != newValue {
			scopedMap[scopedAttr] = value
			return scopedMap, true
		}
	}

	return scopedMap, false
}

func (r *Replacer) ByValueForItem(scopedSlice []interface{}, value interface{}) ([]interface{}, bool) {
	changed := false
	found := false
	for _, oldValue := range scopedSlice {
		if oldValue == value {
			found = true
			break
		}
	}
	if !found {
		changed = true
		scopedSlice = append(scopedSlice, value)
	}
	return scopedSlice, changed
}

func (r *Replacer) ByValueExpressionForItem(scopedMaps []map[string]interface{}, expr filter.Expression, value interface{}) ([]map[string]interface{}, bool) {
	changed := false
	newValues := []map[string]interface{}{}
	for _, oldValue := range scopedMaps {
		if !isMatchExpression(oldValue, expr) {
			newValues = append(newValues, oldValue)
		} else {
			newMap, ok := value.(map[string]interface{})
			fmt.Printf("\nnewMap, ok = %v, %v\n", newMap, ok)
			if ok && !eqMap(oldValue, newMap) {
				changed = true
				newValues = append(newValues, newMap)
			}
		}
	}
	return newValues, changed
}

func (r *Replacer) ByValueExpressionForAttribute(scopedMaps []map[string]interface{}, expr filter.Expression, subAttr string, value interface{}) ([]map[string]interface{}, bool) {
	changed := false
	newValues := []map[string]interface{}{}
	for _, oldValue := range scopedMaps {
		if !isMatchExpression(oldValue, expr) {
			newValues = append(newValues, oldValue)
		} else {
			oldAttrValue, ok := oldValue[subAttr]
			if !ok || oldAttrValue != value {
				changed = true
				oldValue[subAttr] = value
				newValues = append(newValues, oldValue)
			} else {
				newValues = append(newValues, oldValue)
			}
		}
	}
	return newValues, changed
}

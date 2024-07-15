package scimpatch

import "github.com/scim2/filter-parser/v2"

type Remover struct{}

var remover *Remover

func (r *Remover) Direct(scopedMap map[string]interface{}, scopedAttr string, value interface{}) (map[string]interface{}, bool) {
	if _, ok := scopedMap[scopedAttr]; ok {
		delete(scopedMap, scopedAttr)
		return scopedMap, true
	}
	return scopedMap, false
}

func (r *Remover) ByValueForItem(scopedSlice []interface{}, value interface{}) ([]interface{}, bool) {
	changed := false
	newValues := []interface{}{}
	for _, oldValue := range scopedSlice {
		if oldValue != value {
			newValues = append(newValues, oldValue)
		} else {
			changed = true
		}
	}
	return newValues, changed
}

func (r *Remover) ByValueExpressionForItem(scopedMaps []map[string]interface{}, expr filter.Expression, value interface{}) ([]map[string]interface{}, bool) {
	changed := false
	newValues := []map[string]interface{}{}
	for _, oldValue := range scopedMaps {
		if !isMatchExpression(oldValue, expr) {
			newValues = append(newValues, oldValue)
		} else {
			changed = true
		}
	}
	return newValues, changed
}

func (r *Remover) ByValueExpressionForAttribute(scopedMaps []map[string]interface{}, expr filter.Expression, subAttr string, value interface{}) ([]map[string]interface{}, bool) {
	changed := false
	newValues := []map[string]interface{}{}
	for _, oldValue := range scopedMaps {
		if !isMatchExpression(oldValue, expr) {
			newValues = append(newValues, oldValue)
		} else {
			_, ok := oldValue[subAttr]
			if ok {
				changed = true
				delete(oldValue, subAttr)
				newValues = append(newValues, oldValue)
			} else {
				newValues = append(newValues, oldValue)
			}
		}
	}
	return newValues, changed
}

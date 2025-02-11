package scimpatch

import "github.com/scim2/filter-parser/v2"

type remover struct{}

var removerInstance *remover

func (r *remover) Direct(scopedMap map[string]interface{}, scopedAttr string, value interface{}) bool {
	if _, ok := scopedMap[scopedAttr]; ok {
		delete(scopedMap, scopedAttr)
		return true
	}
	return false
}

func (r *remover) ByValueForItem(scopedSlice []interface{}, value interface{}) ([]interface{}, bool) {
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

func (r *remover) ByValueExpressionForItem(scopedMaps []map[string]interface{}, expr filter.Expression, value interface{}) ([]map[string]interface{}, bool) {
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

func (r *remover) ByValueExpressionForAttribute(scopedMaps []map[string]interface{}, expr filter.Expression, subAttr string, value interface{}) ([]map[string]interface{}, bool) {
	changed := false
	newValues := []map[string]interface{}{}
	for _, oldValue := range scopedMaps {
		if !isMatchExpression(oldValue, expr) {
			newValues = append(newValues, oldValue)
		} else {
			if _, ok := oldValue[subAttr]; ok {
				changed = true
				delete(oldValue, subAttr)
			}
			newValues = append(newValues, oldValue)
		}
	}
	return newValues, changed
}

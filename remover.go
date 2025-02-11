package scimpatch

import (
	"context"

	"github.com/scim2/filter-parser/v2"
)

type remover struct{}

var removerInstance *remover

func (r *remover) Direct(ctx context.Context, scopedMap map[string]interface{}, scopedAttr string, value interface{}) bool {
	if _, ok := scopedMap[scopedAttr]; ok {
		delete(scopedMap, scopedAttr)
		return true
	}
	return false
}

func (r *remover) ByValueExpressionForItem(ctx context.Context, scopedMaps []map[string]interface{}, expr filter.Expression, value interface{}) ([]map[string]interface{}, bool) {
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

func (r *remover) ByValueExpressionForAttribute(ctx context.Context, scopedMaps []map[string]interface{}, expr filter.Expression, subAttr string, value interface{}) ([]map[string]interface{}, bool) {
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

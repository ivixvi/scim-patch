package scimpatch

import "github.com/scim2/filter-parser/v2"

type Operator interface {
	Direct(scopedMap map[string]interface{}, scopedAttr string, value interface{}) (map[string]interface{}, bool)
	ByValueForItem(scopedSlice []interface{}, value interface{}) ([]interface{}, bool)
	ByValueExpressionForItem(scopedMaps []map[string]interface{}, expr filter.Expression, value interface{}) ([]map[string]interface{}, bool)
	ByValueExpressionForAttribute(scopedMaps []map[string]interface{}, expr filter.Expression, subAttr string, value interface{}) ([]map[string]interface{}, bool)
}

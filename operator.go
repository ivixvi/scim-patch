package scimpatch

import "github.com/scim2/filter-parser/v2"

// Operator は Patch Operation の各操作のドメインとなるインターフェースです。
// Direct は map とその map 内で更新対象となる属性と更新後の値を受け取って更新後のmapと変更有無を返却します
// この関数のみ、pathが未指定の場合でも利用されます
// ByValueForItem は path が指定されているときに、MultiValued な属性名が指定されたとき対象のスライスと指定された値を受け取って、更新後のスライスと変更有無を返却します。
// ByValueExpressionForItem は 対象の属性が、MultiValuedComplexAttribute で path にて valFilter が指定されているときにそれを受けとって更新後のスライスと変更有無を返却します。
// ByValueExpressionForAttribute は 対象の属性が、MultiValuedComplexAttribute で path にて valFilter と subAttr が指定されているときにそれを受けとって更新後のスライスと変更有無を返却します。
type Operator interface {
	Direct(scopedMap map[string]interface{}, scopedAttr string, value interface{}) bool
	ByValueForItem(scopedSlice []interface{}, value interface{}) ([]interface{}, bool)
	ByValueExpressionForItem(scopedMaps []map[string]interface{}, expr filter.Expression, value interface{}) ([]map[string]interface{}, bool)
	ByValueExpressionForAttribute(scopedMaps []map[string]interface{}, expr filter.Expression, subAttr string, value interface{}) ([]map[string]interface{}, bool)
}

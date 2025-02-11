package scimpatch

import (
	"context"

	"github.com/scim2/filter-parser/v2"
)

// Operator は Patch Operation の各操作のドメインとなるインターフェースです。
// Direct は map とその map 内で更新対象となる属性と更新後の値を受け取って更新後のmapと変更有無を返却します
// この関数のみ、pathが未指定の場合でも利用されます
// ByValueExpressionForItem は 対象の属性が、MultiValuedComplexAttribute で path にて valFilter が指定されているときにそれを受けとって更新後のスライスと変更有無を返却します。
// ByValueExpressionForAttribute は 対象の属性が、MultiValuedComplexAttribute で path にて valFilter と subAttr が指定されているときにそれを受けとって更新後のスライスと変更有無を返却します。
type Operator interface {
	Direct(ctx context.Context, scopedMap map[string]interface{}, scopedAttr string, value interface{}) bool
	ByValueExpressionForItem(ctx context.Context, scopedMaps []map[string]interface{}, expr filter.Expression, value interface{}) ([]map[string]interface{}, bool)
	ByValueExpressionForAttribute(ctx context.Context, scopedMaps []map[string]interface{}, expr filter.Expression, subAttr string, value interface{}) ([]map[string]interface{}, bool)
}

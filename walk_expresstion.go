package scimpatch

import (
	"github.com/scim2/filter-parser/v2"
)

// isMatchExpression は、 value が expr の条件に一致するかどうかを確認します
func isMatchExpression(value map[string]interface{}, expr filter.Expression) bool {
	switch typedExpr := expr.(type) {
	case *filter.AttributeExpression:
		attrValue, ok := value[typedExpr.AttributePath.AttributeName]
		if !ok {
			return false
		}
		switch typedExpr.Operator {
		case "eq":
			return typedExpr.CompareValue == attrValue
		}
	}
	return false
}

// toMap は、 expr を解釈して map を作成します
func toMap(expr filter.Expression) map[string]interface{} {
	switch typedExpr := expr.(type) {
	case *filter.AttributeExpression:
		switch typedExpr.Operator {
		case "eq":
			return map[string]interface{}{
				typedExpr.AttributePath.AttributeName: typedExpr.CompareValue,
			}
		}
	}
	return map[string]interface{}{}
}

package scimpatch

import (
	"github.com/scim2/filter-parser/v2"
)

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
	return true
}

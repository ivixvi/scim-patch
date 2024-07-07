package scimpatch

import (
	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/schema"
	filter "github.com/scim2/filter-parser/v2"
)

var (
	attributeMutabilityImmutable = "immutable"
	attributeMutabilityReadOnly  = "readOnly"
)

func cannotBePatched(op string, attr schema.CoreAttribute) bool {
	return isImmutable(op, attr) || isReadOnly(attr)
}

func isImmutable(op string, attr schema.CoreAttribute) bool {
	return attr.Mutability() == attributeMutabilityImmutable && (op == scim.PatchOperationReplace || op == scim.PatchOperationRemove)
}

func isReadOnly(attr schema.CoreAttribute) bool {
	return attr.Mutability() == attributeMutabilityReadOnly
}

func containsURIPrefix(path *filter.Path) (string, bool) {
	ok := false
	uriPrefix := ""
	if path != nil && path.AttributePath.URIPrefix != nil {
		ok = true
		uriPrefix = *path.AttributePath.URIPrefix
	}
	return uriPrefix, ok
}

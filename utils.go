package scimpatch

import (
	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/schema"
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

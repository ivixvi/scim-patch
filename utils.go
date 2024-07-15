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

func areEveryItemsMap(s interface{}) ([]map[string]interface{}, bool) {
	switch typed := s.(type) {
	case []map[string]interface{}:
		return typed, true
	case []interface{}:
		maps := []map[string]interface{}{}
		for _, item := range typed {
			if map_, ok := item.(map[string]interface{}); ok {
				maps = append(maps, map_)
			} else {
				return nil, false
			}
		}
		return maps, true
	default:
		return nil, false
	}
}

func eqMap(m1 map[string]interface{}, m2 map[string]interface{}) bool {
	if len(m1) != len(m2) {
		return false
	}
	for m1k, m1v := range m1 {
		if m2v, ok := m2[m1k]; !ok || m2v != m1v {
			return false
		}
	}
	return true
}

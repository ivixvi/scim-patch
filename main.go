package scimpatch

import (
	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/schema"
)

type Patcher struct {
	Schema           schema.Schema
	SchemaExtensions []scim.SchemaExtension
}

// Apply is returns a argument.
func (p Patcher) Apply(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	return data, nil
}

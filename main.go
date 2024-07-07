package scimpatch

import (
	"strings"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/errors"
	"github.com/elimity-com/scim/schema"
)

type Patcher struct {
	Schema           schema.Schema
	SchemaExtensions []scim.SchemaExtension
}

// Apply is returns a argument.
func (p Patcher) Apply(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	switch strings.ToLower(op.Op) {
	case scim.PatchOperationAdd:
		return p.Add(op, data)
	case scim.PatchOperationReplace:
		return p.Replace(op, data)
	case scim.PatchOperationRemove:
		return p.Remove(op, data)
	}
	return data, nil
}

// Add is returns a argument.
func (p Patcher) Add(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	return data, nil
}

// Replace is returns a argument.
func (p Patcher) Replace(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	return data, nil
}

// Remove is returns a argument.
func (p Patcher) Remove(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	if op.Path == nil {
		// 3.5.2.2. Remove Operation
		// If "path" is unspecified, the operation fails with HTTP status code 400 and a "scimType" error code of "noTarget".
		return scim.ResourceAttributes{}, errors.ScimErrorNoTarget
	}
	return data, nil
}

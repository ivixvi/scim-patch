package scimpatch

import (
	"github.com/elimity-com/scim"
)

// Apply is returns a argument.
func Apply(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	return data, nil
}

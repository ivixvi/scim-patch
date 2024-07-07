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

// Apply は RFC7644 3.5.2.  Modifying with PATCH の実装です。
// data に op が適用された ResourceAttributes を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2
func (p Patcher) Apply(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	switch strings.ToLower(op.Op) {
	case scim.PatchOperationAdd:
		return p.add(op, data)
	case scim.PatchOperationReplace:
		return p.replace(op, data)
	case scim.PatchOperationRemove:
		return p.remove(op, data)
	}
	return data, nil
}

// add は RFC7644 3.5.2.1. Add Operation の実装です。
// data に op が適用された ResourceAttributes を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.1
// 基本は Validated な op を想定しているため、エラーハンドリングは属性を確認するうえで対応することになる最小限のチェックとなっています。
func (p Patcher) add(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	return data, nil
}

// replace は RFC7644 3.5.2.3. Replace Operation の実装です。
// data に op が適用された ResourceAttributes を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.3
// 基本は Validated な op を想定しているため、エラーハンドリングは属性を確認するうえで対応することになる最小限のチェックとなっています。
func (p Patcher) replace(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	return data, nil
}

// remove は RFC7644 3.5.2.2. Remove Operation の実装です。
// data に op が適用された ResourceAttributes を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.2
// 基本は Validated な op を想定しているため、エラーハンドリングは属性を確認するうえで対応することになる最小限のチェックとなっています。
func (p Patcher) remove(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	if op.Path == nil {
		// If "path" is unspecified, the operation fails with HTTP status code 400 and a "scimType" error code of "noTarget".
		return scim.ResourceAttributes{}, errors.ScimErrorNoTarget
	}
	return data, nil
}

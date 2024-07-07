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
	retData := data
	if op.Path == nil {
		// If "path" is unspecified, the operation fails with HTTP status code 400 and a "scimType" error code of "noTarget".
		return scim.ResourceAttributes{}, errors.ScimErrorNoTarget
	}
	// Resolve Attribute
	attrName := op.Path.AttributePath.AttributeName
	attr, ok := p.containsAttribute(attrName)
	if !ok {
		return scim.ResourceAttributes{}, errors.ScimErrorInvalidPath
	}
	if cannotBePatched(op.Op, attr) {
		return scim.ResourceAttributes{}, errors.ScimErrorMutability
	}
	uriPrefix, containsURI := containsURIPrefix(op.Path)
	if containsURI {
		data, ok = data[uriPrefix].(map[string]interface{})
		if !ok {
			data = scim.ResourceAttributes{}
		}
	}
	switch {
	case attr.HasSubAttributes() && attr.MultiValued():
	case !attr.HasSubAttributes() && attr.MultiValued():
	case attr.HasSubAttributes() && !attr.MultiValued():
	case !attr.HasSubAttributes() && !attr.MultiValued():
		delete(data, attrName)
	}
	if containsURI {
		if len(data) == 0 {
			delete(retData, uriPrefix)
		} else {
			retData[uriPrefix] = data
		}
	} else {
		retData = data
	}
	return retData, nil
}

func (p Patcher) containsAttribute(attrName string) (schema.CoreAttribute, bool) {
	attr, ok := p.Schema.Attributes.ContainsAttribute(attrName)
	if ok {
		return attr, ok
	}
	for _, schema := range p.SchemaExtensions {
		attr, ok := schema.Schema.Attributes.ContainsAttribute(attrName)
		if ok {
			return attr, ok
		}
	}
	return schema.CoreAttribute{}, false
}

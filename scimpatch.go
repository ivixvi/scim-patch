package scimpatch

import (
	"strings"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/errors"
	"github.com/elimity-com/scim/schema"
)

type Patcher struct {
	schema  schema.Schema
	schemas map[string]schema.Schema
}

// NewPatcher は Pathcher の実態を取得します。
func NewPatcher(s schema.Schema, extentions []schema.Schema) *Patcher {
	schemas := map[string]schema.Schema{s.ID: s}
	for _, s := range extentions {
		schemas[s.ID] = s
	}
	return &Patcher{
		schema:  s,
		schemas: schemas,
	}
}

// Apply は RFC7644 3.5.2.  Modifying with PATCH の実装です。
// data に op が適用された ResourceAttributes と実際に適用されたかどうかの真偽値を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2
func (p *Patcher) Apply(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, bool, error) {
	switch strings.ToLower(op.Op) {
	case scim.PatchOperationAdd:
		return p.add(op, data)
	case scim.PatchOperationReplace:
		return p.replace(op, data)
	case scim.PatchOperationRemove:
		return p.remove(op, data)
	}
	return data, false, nil
}

// add は RFC7644 3.5.2.1. Add Operation の実装です。
// data に op が適用された ResourceAttributes と実際に適用されたかどうかの真偽値を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.1
// 基本は Validated な op を想定しているため、エラーハンドリングは属性を確認するうえで対応することになる最小限のチェックとなっています。
func (p *Patcher) add(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, bool, error) {
	if op.Path == nil {
		return data, false, nil
	}
	return p.pathSpecifiedOperate(op, data, additionner)
}

// replace は RFC7644 3.5.2.3. Replace Operation の実装です。
// data に op が適用された ResourceAttributes と実際に適用されたかどうかの真偽値を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.3
// 基本は Validated な op を想定しているため、エラーハンドリングは属性を確認するうえで対応することになる最小限のチェックとなっています。
func (p *Patcher) replace(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, bool, error) {
	if op.Path == nil {
		return data, false, nil
	}
	return p.pathSpecifiedOperate(op, data, replacer)
}

// remove は RFC7644 3.5.2.2. Remove Operation の実装です。
// data に op が適用された ResourceAttributes と実際に適用されたかどうかの真偽値を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.2
// 基本は Validated な op を想定しているため、エラーハンドリングは属性を確認するうえで対応することになる最小限のチェックとなっています。
func (p *Patcher) remove(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, bool, error) {
	if op.Path == nil {
		// If "path" is unspecified, the operation fails with HTTP status code 400 and a "scimType" error code of "noTarget".
		return scim.ResourceAttributes{}, false, errors.ScimErrorNoTarget
	}
	return p.pathSpecifiedOperate(op, data, remover)
}

// containsAttribute は attrName がサーバーで利用されているスキーマの属性名として適切かを確認し、取得します。
func (p *Patcher) containsAttribute(attrName string) (schema.CoreAttribute, bool) {
	for _, schema := range p.schemas {
		attr, ok := schema.Attributes.ContainsAttribute(attrName)
		if ok {
			return attr, ok
		}
	}
	return schema.CoreAttribute{}, false
}

func (p *Patcher) pathSpecifiedOperate(
	op scim.PatchOperation,
	data scim.ResourceAttributes,
	operator Operator,
) (scim.ResourceAttributes, bool, error) {
	var changed = false
	// Resolve Attribute
	attrName := op.Path.AttributePath.AttributeName
	attr, ok := p.containsAttribute(attrName)
	if !ok {
		return scim.ResourceAttributes{}, false, errors.ScimErrorInvalidPath
	}
	if cannotBePatched(op.Op, attr) {
		return scim.ResourceAttributes{}, false, errors.ScimErrorMutability
	}
	n := NewScopeNavigator(op, data, attr)
	switch {
	case attr.MultiValued() && op.Path.ValueExpression != nil && op.Path.SubAttribute != nil:
		var newValues []map[string]interface{}
		oldValues := n.GetScopedMapSlice()
		newValues, changed = operator.ByValueExpressionForAttribute(oldValues, op.Path.ValueExpression, *op.Path.SubAttribute, op.Value)
		n.ApplyScopedMapSlice(newValues)
	case attr.MultiValued() && op.Path.ValueExpression != nil:
		var newValues []map[string]interface{}
		oldValues := n.GetScopedMapSlice()
		newValues, changed = operator.ByValueExpressionForItem(oldValues, op.Path.ValueExpression, op.Value)
		n.ApplyScopedMapSlice(newValues)
	case !attr.MultiValued() || op.Path.ValueExpression == nil:
		scopedMap, scopedAttr := n.GetScopedMap()
		scopedMap, changed = operator.Direct(scopedMap, scopedAttr, op.Value)
		n.ApplyScopedMap(scopedMap)
		data = n.GetMap()
	}

	return data, changed, nil
}

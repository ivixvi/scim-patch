package scimpatch

import (
	"context"
	"strings"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/errors"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
)

type Patcher struct {
	schema   schema.Schema
	schemas  map[string]schema.Schema
	adder    Operator
	replacer Operator
	remover  Operator
}

// PatcherOpts を利用することで Patcherが利用する各操作の Operator を上書きすることができます。
// 指定しない場合はパッケージデフォルトで実装されている Operator が利用されます。
type PatcherOpts struct {
	Adder    *Operator
	Replacer *Operator
	Remover  *Operator
}

var externalIdAttr = schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
	Description: optional.NewString("A String that is an identifier for the resource as defined by the provisioning client."),
	Name:        "externalId",
}))

// NewPatcher は Patcher の実態を取得します。
func NewPatcher(
	s schema.Schema,
	extensions []schema.Schema,
	opts *PatcherOpts,
) *Patcher {
	schemas := map[string]schema.Schema{s.ID: s}
	for _, s := range extensions {
		schemas[s.ID] = s
	}
	patcher := &Patcher{
		schema:   s,
		schemas:  schemas,
		adder:    adderInstance,
		replacer: replacerInstance,
		remover:  removerInstance,
	}
	if opts != nil {
		if opts.Adder != nil {
			patcher.adder = *opts.Adder
		}
		if opts.Replacer != nil {
			patcher.replacer = *opts.Replacer
		}
		if opts.Remover != nil {
			patcher.remover = *opts.Remover
		}
	}
	return patcher
}

// Apply は RFC7644 3.5.2.  Modifying with PATCH の実装です。
// data に op が適用された ResourceAttributes と実際に適用されたかどうかの真偽値を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2
func (p *Patcher) Apply(ctx context.Context, op scim.PatchOperation, data map[string]interface{}) (map[string]interface{}, bool, error) {
	switch strings.ToLower(op.Op) {
	case scim.PatchOperationAdd:
		return p.add(ctx, op, data)
	case scim.PatchOperationReplace:
		return p.replace(ctx, op, data)
	case scim.PatchOperationRemove:
		return p.remove(ctx, op, data)
	}
	return data, false, nil
}

// add は RFC7644 3.5.2.1. Add Operation の実装です。
// data に op が適用された ResourceAttributes と実際に適用されたかどうかの真偽値を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.1
// 基本は Validated な op を想定しているため、エラーハンドリングは属性を確認するうえで対応することになる最小限のチェックとなっています。
func (p *Patcher) add(ctx context.Context, op scim.PatchOperation, data map[string]interface{}) (map[string]interface{}, bool, error) {
	if op.Path == nil {
		return p.pathUnspecifiedOperate(ctx, op, data, p.adder)
	}
	return p.pathSpecifiedOperate(ctx, op, data, p.adder)
}

// replace は RFC7644 3.5.2.3. Replace Operation の実装です。
// data に op が適用された ResourceAttributes と実際に適用されたかどうかの真偽値を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.3
// 基本は Validated な op を想定しているため、エラーハンドリングは属性を確認するうえで対応することになる最小限のチェックとなっています。
func (p *Patcher) replace(ctx context.Context, op scim.PatchOperation, data map[string]interface{}) (map[string]interface{}, bool, error) {
	if op.Path == nil {
		return p.pathUnspecifiedOperate(ctx, op, data, p.replacer)
	}
	return p.pathSpecifiedOperate(ctx, op, data, p.replacer)
}

// remove は RFC7644 3.5.2.2. Remove Operation の実装です。
// data に op が適用された ResourceAttributes と実際に適用されたかどうかの真偽値を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.2
// 基本は Validated な op を想定しているため、エラーハンドリングは属性を確認するうえで対応することになる最小限のチェックとなっています。
func (p *Patcher) remove(ctx context.Context, op scim.PatchOperation, data map[string]interface{}) (map[string]interface{}, bool, error) {
	if op.Path == nil {
		// If "path" is unspecified, the operation fails with HTTP status code 400 and a "scimType" error code of "noTarget".
		return map[string]interface{}{}, false, errors.ScimErrorNoTarget
	}
	return p.pathSpecifiedOperate(ctx, op, data, p.remover)
}

// containsAttribute は attrName がサーバーで利用されているスキーマの属性名として適切かを確認し、取得します。
func (p *Patcher) containsAttribute(attrName string) (schema.CoreAttribute, bool) {
	for _, schema := range p.schemas {
		attr, ok := schema.Attributes.ContainsAttribute(attrName)
		if ok {
			return attr, ok
		}
	}
	if attrName == externalIdAttr.Name() {
		return externalIdAttr, true
	}
	return schema.CoreAttribute{}, false
}

func (p *Patcher) pathSpecifiedOperate(
	ctx context.Context,
	op scim.PatchOperation,
	data map[string]interface{},
	operator Operator,
) (map[string]interface{}, bool, error) {
	var changed = false
	// Resolve Attribute
	attrName := op.Path.AttributePath.AttributeName
	attr, ok := p.containsAttribute(attrName)
	if !ok {
		return map[string]interface{}{}, false, errors.ScimErrorInvalidPath
	}
	if cannotBePatched(op.Op, attr) {
		return map[string]interface{}{}, false, errors.ScimErrorMutability
	}
	n := newScopeNavigator(op, data, attr)
	switch {
	// request path is `attr[expr].subAttr`
	case attr.MultiValued() && op.Path.ValueExpression != nil && op.Path.SubAttribute != nil:
		var newValues []map[string]interface{}
		oldValues := n.GetScopedMapSlice()
		newValues, changed = operator.ByValueExpressionForAttribute(ctx, oldValues, op.Path.ValueExpression, *op.Path.SubAttribute, op.Value)
		n.ApplyScopedMapSlice(newValues)
	// request path is `attr[expr]`
	case attr.MultiValued() && op.Path.ValueExpression != nil:
		var newValues []map[string]interface{}
		oldValues := n.GetScopedMapSlice()
		newValues, changed = operator.ByValueExpressionForItem(ctx, oldValues, op.Path.ValueExpression, op.Value)
		n.ApplyScopedMapSlice(newValues)
	// request path is `attr`, `attr.subAttr`
	case !attr.MultiValued() || op.Path.ValueExpression == nil:
		scopedMap, scopedAttr := n.GetScopedMap()
		changed = operator.Direct(ctx, scopedMap, scopedAttr, op.Value)
		n.ApplyScopedMap(scopedMap)
	}

	return data, changed, nil
}

func (p *Patcher) pathUnspecifiedOperate(
	ctx context.Context,
	op scim.PatchOperation,
	data map[string]interface{},
	operator Operator,
) (map[string]interface{}, bool, error) {
	switch newMap := op.Value.(type) {
	case map[string]interface{}:
		changed := false
		for attr, value := range newMap {
			uriPrefix, ok := p.schemas[attr]
			// Core Attributes
			if !ok {
				changed = changed || operator.Direct(ctx, data, attr, value)
				continue
			}

			// Schema Extension Attributes
			oldMap, ok := data[uriPrefix.ID].(map[string]interface{})

			// if not exists, write all attributes
			if !ok {
				changed = true
				data[uriPrefix.ID] = value
				continue
			}

			// if exists, write by every attributes
			if newUriMap, ok := value.(map[string]interface{}); ok {
				for scopedAttr, scopedValue := range newUriMap {
					changed = changed || operator.Direct(ctx, oldMap, scopedAttr, scopedValue)
				}
			}
		}
		return data, changed, nil
	default:
		// unexpected input
		return data, false, nil
	}
}

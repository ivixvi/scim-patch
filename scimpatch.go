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
// data に op が適用された ResourceAttributes を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2
func (p *Patcher) Apply(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
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
func (p *Patcher) add(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	return data, nil
}

// replace は RFC7644 3.5.2.3. Replace Operation の実装です。
// data に op が適用された ResourceAttributes を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.3
// 基本は Validated な op を想定しているため、エラーハンドリングは属性を確認するうえで対応することになる最小限のチェックとなっています。
func (p *Patcher) replace(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
	return data, nil
}

// remove は RFC7644 3.5.2.2. Remove Operation の実装です。
// data に op が適用された ResourceAttributes を返却します。
// see. https://datatracker.ietf.org/doc/html/rfc7644#section-3.5.2.2
// 基本は Validated な op を想定しているため、エラーハンドリングは属性を確認するうえで対応することになる最小限のチェックとなっています。
func (p *Patcher) remove(op scim.PatchOperation, data scim.ResourceAttributes) (scim.ResourceAttributes, error) {
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
	scopedMap := p.getScopedMap(op, data)
	switch {
	case !attr.HasSubAttributes() && !attr.MultiValued():
		delete(scopedMap, attrName)
	}
	data = p.setScopedMap(op, data, scopedMap)
	return data, nil
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

// getScopedMap は 処理対象であるmapまでのスコープをたどり該当のmapを返却します
func (p *Patcher) getScopedMap(op scim.PatchOperation, data scim.ResourceAttributes) scim.ResourceAttributes {
	uriPrefix, containsURI := containsURIPrefix(op.Path)
	if containsURI {
		data, ok := data[uriPrefix].(map[string]interface{})
		if !ok {
			data = scim.ResourceAttributes{}
		}
		return data
	}

	return data
}

// setScopedMap は 処理対象であるmapまでのスコープをたどりscopedMapに置換したdataを返却します
func (p *Patcher) setScopedMap(op scim.PatchOperation, data scim.ResourceAttributes, scopedMap scim.ResourceAttributes) scim.ResourceAttributes {
	uriPrefix, containsURI := containsURIPrefix(op.Path)
	if containsURI {
		if len(scopedMap) == 0 {
			delete(data, uriPrefix)
		} else {
			data[uriPrefix] = scopedMap
		}
	}
	return data
}

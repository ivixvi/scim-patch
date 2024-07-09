package scimpatch

import (
	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/schema"
)

// getScopedMap は 処理対象であるmapまでのスコープをたどり該当のmapを返却します
func (p *Patcher) getScopedMap(op scim.PatchOperation, data scim.ResourceAttributes, attr schema.CoreAttribute) (scim.ResourceAttributes, string) {
	// initialize returns
	attrName := attr.Name()
	data = p.getURIScopedMap(op, data, attr)
	data, attrName = p.getAttributeScopedMap(op, data, attr)

	return data, attrName
}

// setScopedMap は 処理対象であるmapまでのスコープをたどりscopedMapに置換したdataを返却します
func (p *Patcher) setScopedMap(op scim.PatchOperation, data scim.ResourceAttributes, scopedMap scim.ResourceAttributes, attr schema.CoreAttribute) scim.ResourceAttributes {
	uriScoped := p.getURIScopedMap(op, data, attr)

	if attr.HasSubAttributes() && op.Path != nil && op.Path.AttributePath.SubAttribute != nil {
		if len(scopedMap) == 0 {
			delete(uriScoped, op.Path.AttributePath.AttributeName)
		} else {
			uriScoped[op.Path.AttributePath.AttributeName] = scopedMap
		}
	}

	uriPrefix, containsURI := containsURIPrefix(op.Path)
	if containsURI {
		if len(uriScoped) == 0 {
			delete(data, uriPrefix)
		} else {
			data[uriPrefix] = uriScoped
		}
	}
	return data
}

// getURIScopedMap は URIに応じて、処理対象のMapを返却します
func (p *Patcher) getURIScopedMap(op scim.PatchOperation, data scim.ResourceAttributes, attr schema.CoreAttribute) scim.ResourceAttributes {
	uriScoped := data
	uriPrefix, containsURI := containsURIPrefix(op.Path)
	if containsURI {
		data_, ok := data[uriPrefix].(map[string]interface{})
		switch ok {
		case true:
			uriScoped = data_
		case false:
			uriScoped = scim.ResourceAttributes{}
		}
	}
	return uriScoped
}

// getAttributeScopedMap は 属性に応じて、処理対象のMapを返却します
func (p *Patcher) getAttributeScopedMap(op scim.PatchOperation, data scim.ResourceAttributes, attr schema.CoreAttribute) (scim.ResourceAttributes, string) {
	// initialize returns
	attrName := attr.Name()
	if attr.HasSubAttributes() && op.Path != nil && op.Path.AttributePath.SubAttribute != nil {
		data_, ok := data[op.Path.AttributePath.AttributeName].(map[string]interface{})
		switch ok {
		case true:
			data = data_
			attrName = *op.Path.AttributePath.SubAttribute
		case false:
			data = scim.ResourceAttributes{}
		}
	}
	return data, attrName
}

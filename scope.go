package scimpatch

import (
	"fmt"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/schema"
)

type ScopeNavigator struct {
	op   scim.PatchOperation
	data scim.ResourceAttributes
	attr schema.CoreAttribute
}

func NewScopeNavigator(op scim.PatchOperation, data scim.ResourceAttributes, attr schema.CoreAttribute) *ScopeNavigator {
	return &ScopeNavigator{
		op:   op,
		data: data,
		attr: attr,
	}
}

// GetMap は 処理対象であるmapまでのスコープをたどり該当のmapを返却します
func (n *ScopeNavigator) GetMap() scim.ResourceAttributes {
	return n.data
}

// GetScopedMap は 処理対象であるmapまでのスコープをたどり該当のmapを返却します
func (n *ScopeNavigator) GetScopedMap() (scim.ResourceAttributes, string) {
	return n.getAttributeScopedMap()
}

// ApplyScopedMap は 処理対象であるmapまでのスコープをたどりscopedMapに置換します
func (n *ScopeNavigator) ApplyScopedMap(scopedMap scim.ResourceAttributes) {
	uriScoped := n.getURIScopedMap()
	if n.attr.HasSubAttributes() && n.op.Path != nil && n.op.Path.AttributePath.SubAttribute != nil {
		if len(scopedMap) == 0 {
			delete(uriScoped, n.op.Path.AttributePath.AttributeName)
		} else {
			uriScoped[n.op.Path.AttributePath.AttributeName] = scopedMap
		}
	}

	data := n.data
	uriPrefix, containsURI := n.containsURIPrefix()
	if containsURI {
		if len(uriScoped) == 0 {
			delete(data, uriPrefix)
		} else {
			data[uriPrefix] = uriScoped
		}
	}
	fmt.Printf("%v\n", data)
	n.data = data
}

// getURIScopedMap は URIに応じて、処理対象のMapを返却します
func (n *ScopeNavigator) getURIScopedMap() scim.ResourceAttributes {
	uriScoped := n.data
	uriPrefix, ok := n.containsURIPrefix()
	uriScoped = n.navigateToMap(uriScoped, uriPrefix, ok)
	return uriScoped
}

// getAttributeScopedMap は 属性に応じて、処理対象のMapを返却します
func (n *ScopeNavigator) getAttributeScopedMap() (scim.ResourceAttributes, string) {
	// initialize returns
	data := n.getURIScopedMap()
	subAttrName, ok := n.requiredSubAttributes()
	data = n.navigateToMap(data, n.attr.Name(), ok)
	return data, subAttrName
}

func (n *ScopeNavigator) navigateToMap(data map[string]interface{}, attr string, ok bool) scim.ResourceAttributes {
	if ok {
		data_, ok := data[attr].(map[string]interface{})
		switch ok {
		case true:
			data = data_
		case false:
			data = scim.ResourceAttributes{}
		}
	}
	return data
}

func (n *ScopeNavigator) containsURIPrefix() (string, bool) {
	ok := false
	uriPrefix := ""
	if n.op.Path != nil && n.op.Path.AttributePath.URIPrefix != nil {
		ok = true
		uriPrefix = *n.op.Path.AttributePath.URIPrefix
	}
	return uriPrefix, ok
}

func (n *ScopeNavigator) requiredSubAttributes() (string, bool) {
	ok := false
	subAttr := n.attr.Name()
	if n.attr.HasSubAttributes() && n.op.Path != nil && n.op.Path.AttributePath.SubAttribute != nil {
		ok = true
		subAttr = *n.op.Path.AttributePath.SubAttribute
	}
	return subAttr, ok
}

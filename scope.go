package scimpatch

import (
	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/schema"
)

type scopeNavigator struct {
	op   scim.PatchOperation
	data map[string]interface{}
	attr schema.CoreAttribute
}

func newScopeNavigator(op scim.PatchOperation, data map[string]interface{}, attr schema.CoreAttribute) *scopeNavigator {
	return &scopeNavigator{
		op:   op,
		data: data,
		attr: attr,
	}
}

// GetMap は 処理対象であるmapまでのスコープをたどり該当のmapを返却します
func (n *scopeNavigator) GetMap() map[string]interface{} {
	return n.data
}

// ApplyScopedMap は 処理対象であるmapまでのスコープをたどりscopedMapに置換します
func (n *scopeNavigator) ApplyScopedMap(scopedMap map[string]interface{}) {
	uriScoped := n.GetURIScopedMap()
	if _, required := n.requiredSubAttributes(); required {
		attachToMap(uriScoped, scopedMap, n.attr.Name(), required)
	} else {
		uriScoped = scopedMap
	}

	uriPrefix, containsURI := n.containsURIPrefix()
	attachToMap(n.data, uriScoped, uriPrefix, containsURI)
}

// ApplyScopedMapSlice は 処理対象であるmapまでのスコープをたどりscopedMapに置換します
func (n *scopeNavigator) ApplyScopedMapSlice(scopedMapSlice []map[string]interface{}) {
	uriScoped := n.GetURIScopedMap()
	attachToMapSlice(uriScoped, scopedMapSlice, n.attr.Name(), true)
	uriPrefix, containsURI := n.containsURIPrefix()
	attachToMap(n.data, uriScoped, uriPrefix, containsURI)
}

// GetURIScopedMap は URIに応じて、処理対象のMapを返却します
func (n *scopeNavigator) GetURIScopedMap() map[string]interface{} {
	uriScoped := n.data
	uriPrefix, ok := n.containsURIPrefix()
	uriScoped = navigateToMap(uriScoped, uriPrefix, ok)
	return uriScoped
}

// GetScopedMap は 属性に応じて、処理対象のMapを返却します
func (n *scopeNavigator) GetScopedMap() (map[string]interface{}, string) {
	// initialize returns
	data := n.GetURIScopedMap()
	subAttrName, ok := n.requiredSubAttributes()
	data = navigateToMap(data, n.attr.Name(), ok)
	return data, subAttrName
}

// GetScopedMapSlice は 属性に応じて、処理対象のMapを返却します
func (n *scopeNavigator) GetScopedMapSlice() []map[string]interface{} {
	// initialize returns
	scoped := n.GetURIScopedMap()
	scopedSlice := navigateToMapSlice(scoped, n.attr.Name(), true)
	return scopedSlice
}

// containsURIPrefix は対象の属性がURIPrefixを持ったmapの中に格納されているかどうかを判断します
func (n *scopeNavigator) containsURIPrefix() (string, bool) {
	ok := false
	uriPrefix := ""
	if n.op.Path != nil && n.op.Path.AttributePath.URIPrefix != nil {
		ok = true
		uriPrefix = *n.op.Path.AttributePath.URIPrefix
	}
	return uriPrefix, ok
}

// requiredSubAttributes は対象の属性がサブ属性を保持したマップであるかどうかと、サブ属性が対象となったPatchOperationかどうかを判断します
func (n *scopeNavigator) requiredSubAttributes() (string, bool) {
	ok := false
	subAttr := n.attr.Name()
	if n.attr.HasSubAttributes() && n.op.Path != nil && n.op.Path.AttributePath.SubAttribute != nil {
		ok = true
		subAttr = *n.op.Path.AttributePath.SubAttribute
	}
	return subAttr, ok
}

// navigateToMap は必要に応じて、パスをたどる処理です
func navigateToMap(data map[string]interface{}, attr string, ok bool) map[string]interface{} {
	if ok {
		data_, ok := data[attr].(map[string]interface{})
		switch ok {
		case true:
			data = data_
		case false:
			data = map[string]interface{}{}
		}
	}
	return data
}

// attachToMap は必要に応じて、パスを戻す処理です
func attachToMap(data map[string]interface{}, scoped map[string]interface{}, attr string, ok bool) {
	if ok {
		if len(scoped) == 0 {
			delete(data, attr)
		} else {
			data[attr] = scoped
		}
	}
}

// navigateToMapSlice は必要に応じて、パスをたどる処理です
func navigateToMapSlice(data map[string]interface{}, attr string, ok bool) []map[string]interface{} {
	ret := []map[string]interface{}{}
	if ok {
		value, ok := data[attr]
		switch ok {
		case true:
			switch typedSlice := value.(type) {
			case []map[string]interface{}:
				ret = typedSlice
			// 各々の item が map として変換できる可能性があるため、一つ一つ確認する必要がある
			case []interface{}:
				for _, item := range typedSlice {
					if mappedItem, ok := item.(map[string]interface{}); ok {
						ret = append(ret, mappedItem)
					}
				}
			}
		}
	}
	return ret
}

// attachToMapSlice は必要に応じて、パスを戻す処理です
func attachToMapSlice(data map[string]interface{}, scoped []map[string]interface{}, attr string, ok bool) {
	if ok {
		if len(scoped) == 0 {
			delete(data, attr)
		} else {
			data[attr] = scoped
		}
	}
}

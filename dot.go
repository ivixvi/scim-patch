package scimpatch

import (
	"strings"
)

// Resolve an attribute name with dot notation ("name.givenName") to a new scopedMap ("name") and scopedAttr ("givenName")
// This is used prominently by MS Entra. See https://learn.microsoft.com/en-us/entra/identity/app-provisioning/application-provisioning-config-problem-scim-compatibility#flags-to-alter-the-scim-behavior
func resolveDotNotationAttribute(scopedMap map[string]interface{}, scopedAttr string) (map[string]interface{}, string) {
	attrParts := strings.SplitN(scopedAttr, ".", 2)
	if len(attrParts) == 1 {
		return scopedMap, scopedAttr
	}

	if subMap, exists := scopedMap[attrParts[0]]; exists {
		scopedMap = subMap.(map[string]interface{})
	} else {
		subMap := map[string]interface{}{}
		scopedMap[attrParts[0]] = subMap
		scopedMap = subMap
	}
	scopedAttr = attrParts[1]

	return scopedMap, scopedAttr
}

package scimpatch_test

import (
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
)

var TestExtensionSchema = schema.Schema{
	Description: optional.NewString("test"),
	ID:          "urn:ivixvi:testSchema",
	Name:        optional.NewString("test"),
	Attributes: []schema.CoreAttribute{
		schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
			Name:        "testString",
			MultiValued: true,
		})),
	},
}

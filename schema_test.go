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

// userSchemaWithExternalId returns CoreUserSchema with externalId added,
// since ApplyPatch resolves attributes from the main schema and externalId
// is a common attribute (RFC 7643 Section 3.1) not included in CoreUserSchema.
func userSchemaWithExternalId() schema.Schema {
	s := schema.CoreUserSchema()
	s.Attributes = append(s.Attributes,
		schema.SimpleCoreAttribute(schema.SimpleStringParams(schema.StringParams{
			Name: "externalId",
		})),
	)
	return s
}

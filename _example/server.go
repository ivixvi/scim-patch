package main

import (
	"log"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
	scimpatch "github.com/ivixvi/scim-patch"
)

func newTestServer() scim.Server {
	userSchema := schema.CoreUserSchema()
	enterprizeUserSchema := schema.ExtensionEnterpriseUser()
	patcher := scimpatch.NewPatcher(userSchema, []schema.Schema{enterprizeUserSchema}, nil)
	s, err := scim.NewServer(
		&scim.ServerArgs{
			ServiceProviderConfig: &scim.ServiceProviderConfig{},
			ResourceTypes: []scim.ResourceType{
				{
					ID:          optional.NewString("User"),
					Name:        "User",
					Endpoint:    "/Users",
					Description: optional.NewString("User Account"),
					Schema:      schema.CoreUserSchema(),
					SchemaExtensions: []scim.SchemaExtension{
						{Schema: schema.ExtensionEnterpriseUser()},
					},
					Handler: &testResourceHandler{
						data:   map[string]testData{},
						schema: schema.CoreUserSchema(),
						schemaExtensions: []scim.SchemaExtension{
							{Schema: schema.ExtensionEnterpriseUser()},
						},
						patcher: patcher,
					},
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

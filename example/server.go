package main

import (
	"log"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
)

func newTestServer() scim.Server {
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
					Handler: &testResourceHandler{
						data:   map[string]testData{},
						schema: schema.CoreUserSchema(),
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

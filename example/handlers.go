package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/errors"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
	"github.com/ivixvi/scimpatch"
)

type testData struct {
	resourceAttributes scim.ResourceAttributes
	meta               scim.Meta
}

// simple in-memory resource database.
type testResourceHandler struct {
	data             map[string]testData
	schema           schema.Schema
	schemaExtensions []scim.SchemaExtension
	patcher          *scimpatch.Patcher
}

// notImplemented is error returns 501 Not Implemented.
var notImplemented = errors.ScimError{Detail: "Not implemented.", Status: 501}

func (h testResourceHandler) Create(r *http.Request, attributes scim.ResourceAttributes) (scim.Resource, error) {
	// create unique identifier
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	id := fmt.Sprintf("%04d", rng.Intn(9999))

	// store resource
	now := time.Now()
	meta := scim.Meta{
		Created:      &now,
		LastModified: &now,
		Version:      fmt.Sprintf("v%s", id),
	}
	h.data[id] = testData{
		resourceAttributes: attributes,
		meta:               meta,
	}

	// return stored resource
	return scim.Resource{
		ID:         id,
		ExternalID: h.externalID(attributes),
		Attributes: attributes,
		Meta:       meta,
	}, nil
}

func (h testResourceHandler) Delete(r *http.Request, id string) error {
	// check if resource exists
	_, ok := h.data[id]
	if !ok {
		return errors.ScimErrorResourceNotFound(id)
	}

	// delete resource
	delete(h.data, id)

	return nil
}

func (h testResourceHandler) Get(r *http.Request, id string) (scim.Resource, error) {
	// check if resource exists and get resource
	data, ok := h.data[id]
	if !ok {
		return scim.Resource{}, errors.ScimErrorResourceNotFound(id)
	}
	// return resource
	return scim.Resource{
		ID:         id,
		ExternalID: h.externalID(data.resourceAttributes),
		Attributes: data.resourceAttributes,
		Meta:       data.meta,
	}, nil
}

func (h testResourceHandler) GetAll(r *http.Request, params scim.ListRequestParams) (scim.Page, error) {
	return scim.Page{}, notImplemented
}

func (h testResourceHandler) Patch(r *http.Request, id string, operations []scim.PatchOperation) (scim.Resource, error) {
	// check if resource exists and get resource
	data, ok := h.data[id]
	if !ok {
		return scim.Resource{}, errors.ScimErrorResourceNotFound(id)
	}

	// Apply PATCH operations
	var err error
	var changed bool
	for _, op := range operations {
		data.resourceAttributes, changed, err = h.patcher.Apply(op, data.resourceAttributes)
		if err != nil {
			return scim.Resource{}, err
		}
	}

	// store resource
	if changed {
		now := time.Now()
		data.meta.LastModified = &now
		h.data[id] = data
	}

	return scim.Resource{
		ID:         id,
		ExternalID: h.externalID(data.resourceAttributes),
		Attributes: data.resourceAttributes,
		Meta:       data.meta,
	}, nil
}

func (h testResourceHandler) Replace(r *http.Request, id string, attributes scim.ResourceAttributes) (scim.Resource, error) {
	return scim.Resource{}, notImplemented

}

func (h testResourceHandler) externalID(attributes scim.ResourceAttributes) optional.String {
	if eID, ok := attributes["externalId"]; ok {
		externalID, ok := eID.(string)
		if !ok {
			return optional.String{}
		}
		return optional.NewString(externalID)
	}

	return optional.String{}
}

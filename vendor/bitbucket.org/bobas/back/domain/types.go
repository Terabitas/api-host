package domain

import (
	"github.com/juju/errors"
	"github.com/nu7hatch/gouuid"
)

const (
	HANDLER_TYPE_EMAIL HandlerType = "email"
)

type (
	// Endpoint is base model
	Endpoint struct {
		Id             string    `bson:"_id" json:"id"`
		Email          string    `bson:"email" json:"email"`
		Name           string    `bson:"name" json:"name"`
		SampleResponse string    `bson:"sampleResponse" json:"sampleResponse"`
		Handlers       []Handler `bson:"handlers" json:"handlers"`
	}

	// Handler type
	Handler struct {
		Metadata HandlerMetadata `bson:"metadata" json:"metadata"`
		Value    []Email         `bson:"value" json:"value"`
	}

	// HandlerMetadata type
	HandlerMetadata struct {
		Type HandlerType `bson:"type" json:"type"`
	}

	// HandlerType type
	HandlerType string

	// Email type
	Email string

	// RequestResponse type
	RequestResponse struct {
		SampleResponse string `json:"sampleResponse"`
	}

	// Request type
	Request struct {
		Id         string      `bson:"_id" json:"id"`
		EndpointId string      `bson:"eid" json:"eid"`
		Data       interface{} `bson:"data" json:"data"`
	}
)

// MakeEndpoint constructor
func MakeEndpoint(name, email, sampleResponse string, handlers []Handler) (*Endpoint, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &Endpoint{
		Id:             u.String(),
		Name:           name,
		Email:          email,
		SampleResponse: sampleResponse,
		Handlers:       handlers,
	}, nil
}

// MakeRequest constructor
func MakeRequest(eid string, data interface{}) *Request {
	u, err := uuid.NewV4()
	if err != nil {
		return nil
	}
	return &Request{
		Id:         u.String(),
		EndpointId: eid,
		Data:       data,
	}
}

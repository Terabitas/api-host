package storage

import (
	"os"

	"github.com/juju/errors"

	"bitbucket.org/bobas/back/domain"
)

const (
	STORAGE_MONGODB = "mongodb"
)

type (
	// EndpointRepository type
	EndpointRepository interface {
		Save(*domain.Endpoint) error
		GetById(id string) (*domain.Endpoint, error)
	}

	// RequestRepository type
	RequestRepository interface {
		Save(*domain.Request) error
		GetById(id string) (*domain.Request, error)
		GetByEndpointId(eid string) ([]domain.Request, error)
	}
)

// MakeRepository factory
func MakeRepositoryFromEnv() (EndpointRepository, error) {
	storage := os.Getenv("ND_STORAGE")
	switch storage {
	case STORAGE_MONGODB:
		return makeMongoDB(), nil
	}

	return nil, errors.Trace(errors.Errorf("Storage [%s] is not supported", storage))
}

// MakeRequestRepositoryFromEnv factory
func MakeRequestRepositoryFromEnv() (RequestRepository, error) {
	storage := os.Getenv("ND_STORAGE")
	switch storage {
	case STORAGE_MONGODB:
		return makeRequestMongoDB(), nil
	}

	return nil, errors.Trace(errors.Errorf("Storage [%s] is not supported", storage))
}

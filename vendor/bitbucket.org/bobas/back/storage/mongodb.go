package storage

import (
	"bitbucket.org/bobas/back/domain"
	"bitbucket.org/bobas/back/provision"
	"github.com/juju/errors"
	"github.com/nildev/lib/registry"
	"gopkg.in/mgo.v2/bson"
)

type (
	mongodb        struct{}
	mongodbRequest struct{}
)

func makeMongoDB() *mongodb {
	return &mongodb{}
}

func makeRequestMongoDB() *mongodbRequest {
	return &mongodbRequest{}
}

// Save method
func (ds *mongodb) Save(endpoint *domain.Endpoint) error {
	if endpoint == nil {
		return errors.Trace(errors.Errorf("endpoint is nil!"))
	}

	session, err := registry.CreateMongoDBClient()
	if err != nil {
		return err
	}

	collection := session.DB(registry.GetDatabaseName()).C(provision.TABLE_NAME_ENDPOINTS)
	_, err = collection.UpsertId(endpoint.Id, endpoint)

	if err != nil {
		return err
	}

	return nil
}

// GetById method
func (ds *mongodb) GetById(id string) (*domain.Endpoint, error) {
	session, err := registry.CreateMongoDBClient()
	if err != nil {
		return nil, err
	}
	collection := session.DB(registry.GetDatabaseName()).C(provision.TABLE_NAME_ENDPOINTS)
	endpoint := &domain.Endpoint{}
	collection.FindId(id).One(&endpoint)

	return endpoint, nil
}

// Save method
func (ds *mongodbRequest) Save(request *domain.Request) error {
	if request == nil {
		return errors.Trace(errors.Errorf("request is nil!"))
	}

	session, err := registry.CreateMongoDBClient()
	if err != nil {
		return err
	}

	collection := session.DB(registry.GetDatabaseName()).C(provision.TABLE_NAME_REQUESTS)
	_, err = collection.UpsertId(request.Id, request)

	if err != nil {
		return err
	}

	return nil
}

// GetById method
func (ds *mongodbRequest) GetById(id string) (*domain.Request, error) {
	session, err := registry.CreateMongoDBClient()
	if err != nil {
		return nil, err
	}
	collection := session.DB(registry.GetDatabaseName()).C(provision.TABLE_NAME_REQUESTS)
	req := &domain.Request{}
	collection.FindId(id).One(&req)

	return req, nil
}

// GetByEndpointId method
func (ds *mongodbRequest) GetByEndpointId(eid string) ([]domain.Request, error) {
	session, err := registry.CreateMongoDBClient()
	if err != nil {
		return nil, err
	}
	collection := session.DB(registry.GetDatabaseName()).C(provision.TABLE_NAME_REQUESTS)
	var req []domain.Request
	collection.Find(bson.M{"eid": eid}).All(&req)

	return req, nil
}

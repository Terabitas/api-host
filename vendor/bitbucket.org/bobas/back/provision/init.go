package provision

import (
	"fmt"

	"github.com/nildev/lib/registry"
	"gopkg.in/mgo.v2"
)

var (
	TABLE_NAME_ENDPOINTS = "Endpoints"
	TABLE_NAME_REQUESTS  = "Requests"
)

// NildevInitMongoDB init
func NildevInitMongoDB() {
	// In which environment we are
	session, err := registry.CreateMongoDBClient()
	if err != nil {
		fmt.Printf("%s", err)
	}

	err = session.DB(registry.GetDatabaseName()).C(TABLE_NAME_ENDPOINTS).Create(&mgo.CollectionInfo{})

	if err != nil {
		fmt.Printf("%s", err)
	}

	err = session.DB(registry.GetDatabaseName()).C(TABLE_NAME_REQUESTS).Create(&mgo.CollectionInfo{})

	if err != nil {
		fmt.Printf("%s", err)
	}

}

func DestroyMongoDB() {
	session, err := registry.CreateMongoDBClient()
	if err != nil {
		fmt.Printf("%s", err)
	}

	err = session.DB(registry.GetDatabaseName()).C(TABLE_NAME_ENDPOINTS).DropCollection()

	if err != nil {
		fmt.Printf("%s", err)
	}

	err = session.DB(registry.GetDatabaseName()).C(TABLE_NAME_REQUESTS).DropCollection()

	if err != nil {
		fmt.Printf("%s", err)
	}
}

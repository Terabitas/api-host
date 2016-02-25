package registry

import (
	"os"

	"net/http"

	"github.com/iron-io/iron_go3/mq"
	"google.golang.org/cloud"
	"google.golang.org/cloud/datastore"
	"google.golang.org/cloud/pubsub"
	"gopkg.in/mgo.v2"
)

var (
	mongoDBClientURL   = "mongodb://localhost:27017/test"
	googleDatastoreURL = "http://localhost:8080/"
	googlePubSubURL    = "http://localhost:8080/"
	googleProjectID    = ""
	databaseName       = "default"
)

// These constants is public API, being taken by lib from environment
const (
	NDMongoDBURL         = "ND_MONGODB_URL"
	NDGoogleDatastoreURL = "ND_GOOGLE_DATASTORE_URL"
	NDGooglePubSubURL    = "ND_GOOGLE_PUBSUB_URL"
	NDGoogleProjectID    = "ND_GOOGLE_PROJECT_ID"
	NDDatabaseName       = "ND_DATABASE_NAME"
	NDEnv                = "ND_ENV"

	EnvDev = "dev"
)

// GetEnv return environment
func GetEnv() string {
	env := os.Getenv(NDEnv)
	if env == "" {
		env = EnvDev
	}

	return env
}

// CreateMongoDBClient create MongoDB client
func CreateMongoDBClient() (*mgo.Session, error) {
	envValue := os.Getenv(NDMongoDBURL)
	if envValue != "" {
		mongoDBClientURL = envValue
	}

	session, err := mgo.Dial(mongoDBClientURL)
	//	session.SetSafe(&mgo.Safe{})
	return session, err
}

// GetDatabaseName returns DB name
func GetDatabaseName() string {
	envValue := os.Getenv(NDDatabaseName)
	if envValue != "" {
		databaseName = envValue
	}
	return databaseName
}

// GetGoogleProjectID returns google project ID
func GetGoogleProjectID() string {
	envValue := os.Getenv(NDGoogleProjectID)
	if envValue != "" {
		googleProjectID = envValue
	}
	return googleProjectID
}

// CreateIronQueue creates and returns Iron.io queue
func CreateIronQueue(name string) (*mq.Queue, error) {
	subscribers := []mq.QueueSubscriber{}
	subscription := mq.PushInfo{
		Retries:      3,
		RetriesDelay: 60,
		ErrorQueue:   "error_queue",
		Subscribers:  subscribers,
	}
	queue_type := "push"
	queueInfo := mq.QueueInfo{Type: queue_type, MessageExpiration: 60, MessageTimeout: 56, Push: &subscription}
	_, err := mq.CreateQueue(name, queueInfo)
	if err != nil {
		return nil, err
	}

	queue := mq.New(name)
	return &queue, nil
}

// CreateGoogleDatastoreClient create client to communicate to google Datastore service
func CreateGoogleDatastoreClient() (*datastore.Client, error) {
	envValue := os.Getenv(NDGoogleDatastoreURL)
	if envValue != "" {
		googleDatastoreURL = envValue
	}

	envValue = os.Getenv(NDGoogleProjectID)
	if envValue != "" {
		googleProjectID = envValue
	}

	ctx := cloud.NewContext(googleProjectID, http.DefaultClient)
	o := []cloud.ClientOption{
		cloud.WithBaseHTTP(http.DefaultClient),
		cloud.WithEndpoint(googleDatastoreURL),
	}
	client, err := datastore.NewClient(ctx, googleProjectID, o...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// CreateGooglePubSubClient returns client communicate with google Pub/Sub service
func CreateGooglePubSubClient() (*pubsub.Client, error) {
	envValue := os.Getenv(NDGooglePubSubURL)
	if envValue != "" {
		googlePubSubURL = envValue
	}

	envValue = os.Getenv(NDGoogleProjectID)
	if envValue != "" {
		googleProjectID = envValue
	}

	ctx := cloud.NewContext(googleProjectID, http.DefaultClient)
	o := []cloud.ClientOption{
		cloud.WithBaseHTTP(http.DefaultClient),
		cloud.WithEndpoint(googleDatastoreURL),
	}
	client, _ := pubsub.NewClient(ctx, googleProjectID, o...)
	return client, nil
}

package common

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config stores the configuration for the database
type Config struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
	Collection   string
}

// DatabaseHelper Helper functions for the Database
type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

// CollectionHelper Helper functions for the Collection
type CollectionHelper interface {
	InsertOne(context.Context, interface{}) (interface{}, error)
}

// ClientHelper Helper functions for the Client
type ClientHelper interface {
	Database(string) DatabaseHelper
	Connect() error
	StartSession() (mongo.Session, error)
}

type mongoClient struct {
	cl *mongo.Client
}
type mongoDatabase struct {
	db *mongo.Database
}
type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSession struct {
	mongo.Session
}

// NewClient creates a new ClientHelper
func NewClient(cnf *Config) (ClientHelper, error) {
	var uri string
	if len(cnf.Username) > 0 || len(cnf.Password) > 0 {
		uri = fmt.Sprintf(`mongodb://%s:%s@%s:%s`,
			cnf.Username,
			cnf.Password,
			cnf.Host,
			cnf.Port,
		)
	} else {
		uri = fmt.Sprintf(`mongodb://%s:%s`,
			cnf.Host,
			cnf.Port,
		)
	}
	c, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	return &mongoClient{cl: c}, err

}

func (mc *mongoClient) Database(dbName string) DatabaseHelper {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.cl.StartSession()
	return &mongoSession{session}, err
}

func (mc *mongoClient) Connect() error {
	return mc.cl.Connect(context.Background())
}

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	return id.InsertedID, err
}

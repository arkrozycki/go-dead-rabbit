package main

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// DatastoreClientHelper
type DatastoreClientHelper interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
	Database(string) DatabaseHelper
	Insert([]byte) error
}

// DatabaseHelper
type DatabaseHelper interface {
	Collection(string) CollectionHelper
}

// CollectionHelper
type CollectionHelper interface {
	InsertOne(context.Context, interface{}) (interface{}, error)
}

// MongoClientHelper
type MongoClientHelper struct {
	cl      *mongo.Client
	dbname  string
	colname string
}

// Connect
func (m *MongoClientHelper) Connect(ctx context.Context) error {
	return m.cl.Connect(context.Background())
}

// Disconnect
func (m *MongoClientHelper) Disconnect(ctx context.Context) error {
	return m.cl.Disconnect(context.Background())
}

// Database
func (m *MongoClientHelper) Database(name string) DatabaseHelper {
	db := m.cl.Database(name)
	return &mongoDatabase{db: db}
}

// Insert
func (m *MongoClientHelper) Insert(doc []byte) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := m.Connect(ctx)
	if err != nil {
		return err
	}
	db := m.Database(m.dbname)
	var bdoc interface{}
	err = bson.UnmarshalJSON(doc, &bdoc)
	id, err := db.Collection(m.colname).InsertOne(context.Background(), bdoc)
	if err == nil {
		log.Debug().Msgf("DATASTORE: inserted record %v", id)
	}
	defer m.Disconnect(ctx)
	return err
}

// mongoDatabase
type mongoDatabase struct {
	db *mongo.Database
}

// Collection
func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

// mongoCollection
type mongoCollection struct {
	coll *mongo.Collection
}

// InsertOne
func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	return id, err
}

// GetDatastoreClient
func GetDatastoreClient(uri string) (DatastoreClientHelper, error) {
	var datastoreClient DatastoreClientHelper
	mcl, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	datastoreClient = &MongoClientHelper{
		cl:      mcl,
		dbname:  Conf.Datastore.Mongodb.Database,
		colname: Conf.Datastore.Mongodb.Collection,
	}
	return datastoreClient, nil
}

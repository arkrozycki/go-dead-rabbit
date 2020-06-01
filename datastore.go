package main

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// DatastoreClientHelper
type DatastoreClientHelper interface {
	Connect() error
	Disconnect(context.Context) error
	Database(string) DatabaseHelper
	Insert([]byte) (string, error)
	Count() (int64, error)
}

// DatabaseHelper
type DatabaseHelper interface {
	Collection(string) CollectionHelper
}

// CollectionHelper
type CollectionHelper interface {
	InsertOne(context.Context, interface{}) (string, error)
	CountDocuments(context.Context) (int64, error)
}

// MongoClientHelper
type MongoClientHelper struct {
	cl      *mongo.Client
	dbname  string
	colname string
}

// Connect
func (m *MongoClientHelper) Connect() error {
	var err error
	// Check the connection
	ctxPing, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = m.cl.Ping(ctxPing, nil)
	// if not connected go ahead and connect
	if err != nil {
		log.Debug().Msg("Connecting to datastore...")
		err = m.cl.Connect(context.Background())
		if err != nil {
			log.Info().Err(err).Msg("error: connect")
		}
	}

	return err
}

// Disconnect
func (m *MongoClientHelper) Disconnect(ctx context.Context) error {
	log.Debug().Msg("Disconnecting ...")
	err := m.cl.Disconnect(ctx)
	if err != nil {
		log.Info().Err(err).Msg("error")
	}
	return err
}

// Database
func (m *MongoClientHelper) Database(name string) DatabaseHelper {
	db := m.cl.Database(name)
	return &mongoDatabase{db: db}
}

// count
func (m *MongoClientHelper) Count() (int64, error) {
	err := m.Connect()
	if err != nil {
		return 0, err
	}

	db := m.Database(m.dbname)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cnt, err := db.Collection(m.colname).CountDocuments(ctx)
	return cnt, err
}

// Insert
func (m *MongoClientHelper) Insert(doc []byte) (string, error) {
	err := m.Connect()
	if err != nil {
		return "", err
	}
	db := m.Database(m.dbname)
	var bdoc interface{}
	err = bson.UnmarshalJSON(doc, &bdoc)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := db.Collection(m.colname).InsertOne(ctx, bdoc)

	if err == nil {
		log.Debug().Msgf("DATASTORE: inserted record %v", result)
	}

	return result, err
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
func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (string, error) {
	result, err := mc.coll.InsertOne(ctx, document)
	objectID := result.InsertedID.(primitive.ObjectID)
	return objectID.String(), err
}

// CountDocuments
func (mc *mongoCollection) CountDocuments(ctx context.Context) (int64, error) {
	return mc.coll.CountDocuments(ctx, bson.M{}, nil)
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

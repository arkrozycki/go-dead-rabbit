package main

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

type DatastoreClientHelper interface {
	Connect(context.Context) error
	Disconnect(context.Context) error
	Database(string) DatabaseHelper
	Insert([]byte) error
}

type DatabaseHelper interface {
	Collection(string) CollectionHelper
}

type CollectionHelper interface {
	InsertOne(context.Context, interface{}) (interface{}, error)
}

type MongoClient struct {
	cl      *mongo.Client
	dbname  string
	colname string
}

func (m *MongoClient) Connect(ctx context.Context) error {
	return m.cl.Connect(context.Background())
}

func (m *MongoClient) Disconnect(ctx context.Context) error {
	return m.cl.Disconnect(context.Background())
}

func (m *MongoClient) Database(name string) DatabaseHelper {
	db := m.cl.Database(name)
	return &mongoDatabase{db: db}
}

func (m *MongoClient) Insert(doc []byte) error {
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
		log.Debug().Msgf("DATASTORE: inserted %v", id)
	}
	defer m.Disconnect(ctx)
	return err
}

type mongoDatabase struct {
	db *mongo.Database
}

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

type mongoCollection struct {
	coll *mongo.Collection
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	return id, err
}

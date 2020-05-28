package main

import (
	"context"
)

type MockDatastoreClient struct{}

func (m *MockDatastoreClient) Connect(ctx context.Context) error {
	return nil
}

func (m *MockDatastoreClient) Disconnect(ctx context.Context) error {
	return nil
}

func (m *MockDatastoreClient) Database(name string) DatabaseHelper {
	return nil
}

func (m *MockDatastoreClient) Insert(doc []byte) error {
	return nil
}

type MockMongoClientHelper struct {
	cl      *MockDatastoreClient
	dbname  string
	colname string
}

func (m *MockMongoClientHelper) Connect(ctx context.Context) error {
	return nil
}

func (m *MockMongoClientHelper) Disconnect(ctx context.Context) error {
	return nil
}

func (m *MockMongoClientHelper) Database(name string) DatabaseHelper {
	return nil
}

func (m *MockMongoClientHelper) Insert(doc []byte) error {
	return nil
}

func (m *MockMongoClientHelper) Count() (int64, error) {
	return int64(10), nil
}

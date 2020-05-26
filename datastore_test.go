package main

import "context"

type MockDSClient struct{}

type MockMongoClient struct {
	cl      *MockDSClient
	dbname  string
	colname string
}

func (m *MockMongoClient) Connect(ctx context.Context) error {
	return nil
}

func (m *MockMongoClient) Disconnect(ctx context.Context) error {
	return nil
}

func (m *MockMongoClient) Database(name string) DatabaseHelper {
	return nil
}

func (m *MockMongoClient) Insert(doc []byte) error {
	return nil
}

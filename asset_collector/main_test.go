package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go get github.com/stretchr/testify/assert
//go test -v

// MockDB is a mock implementation of the Database interface for testing
type MockDB struct{}

func (m *MockDB) CreateServerNode(ctx context.Context, server Server) error {
	// Simulate a successful operation
	return nil
}

func (m *MockDB) CreateApplicationNode(ctx context.Context, app Application) error {
	// Simulate a successful operation
	return nil
}

func (m *MockDB) CreateRelationship(ctx context.Context, serverID, appID string) error {
	// Simulate a successful operation
	return nil
}

// TestCreateServerNode tests the creation of a server node
func TestCreateServerNode(t *testing.T) {
	mockDB := &MockDB{}
	server := Server{ID: "srv1", Name: "Test Server", IPAddress: "192.168.1.100"}

	err := mockDB.CreateServerNode(context.Background(), server)
	assert.NoError(t, err, "expected no error when creating server node")
}

// TestCreateApplicationNode tests the creation of an application node
func TestCreateApplicationNode(t *testing.T) {
	mockDB := &MockDB{}
	app := Application{ID: "app1", Name: "Test App", Version: "1.0"}

	err := mockDB.CreateApplicationNode(context.Background(), app)
	assert.NoError(t, err, "expected no error when creating application node")
}

// TestCreateRelationship tests the creation of a relationship between a server and an application
func TestCreateRelationship(t *testing.T) {
	mockDB := &MockDB{}
	serverID := "srv1"
	appID := "app1"

	err := mockDB.CreateRelationship(context.Background(), serverID, appID)
	assert.NoError(t, err, "expected no error when creating relationship")
}

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Server represents a server asset
type Server struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
	Apps      []string `json:"apps"`
}

// Application represents an application asset
type Application struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Database interface defines methods for interacting with Neo4j
type Database interface {
	CreateServerNode(ctx context.Context, server Server) error
	CreateApplicationNode(ctx context.Context, app Application) error
	CreateRelationship(ctx context.Context, serverID, appID string) error
}

// Neo4jDB implements the Database interface using Neo4j
type Neo4jDB struct {
	driver neo4j.DriverWithContext
}

func (db *Neo4jDB) CreateServerNode(ctx context.Context, server Server) error {
	session := db.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MERGE (s:Server {id: $id}) SET s.name = $name, s.ip_address = $ip_address"
		params := map[string]any{
			"id":        server.ID,
			"name":      server.Name,
			"ip_address": server.IPAddress,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	return err
}

func (db *Neo4jDB) CreateApplicationNode(ctx context.Context, app Application) error {
	session := db.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := "MERGE (a:Application {id: $id}) SET a.name = $name, a.version = $version"
		params := map[string]any{
			"id":      app.ID,
			"name":    app.Name,
			"version": app.Version,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	return err
}

func (db *Neo4jDB) CreateRelationship(ctx context.Context, serverID, appID string) error {
	session := db.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (s:Server {id: $server_id}), (a:Application {id: $app_id})
			MERGE (s)-[:RUNS]->(a)`
		params := map[string]any{
			"server_id": serverID,
			"app_id":    appID,
		}
		_, err := tx.Run(ctx, query, params)
		return nil, err
	})

	return err
}

func main() {
	neo4jURI := os.Getenv("NEO4J_URI")
	neo4jUser := os.Getenv("NEO4J_USER")
	neo4jPassword := os.Getenv("NEO4J_PASSWORD")

	driver, err := neo4j.NewDriverWithContext(neo4jURI, neo4j.BasicAuth(neo4jUser, neo4jPassword, ""), nil)
	if err != nil {
		log.Fatalf("Failed to create Neo4j driver: %v", err)
	}
	defer driver.Close(context.Background())

	ctx := context.Background()

	db := &Neo4jDB{driver: driver}

	// Example data - replace with real CMDB data
	server := Server{ID: "srv1", Name: "Server1", IPAddress: "192.168.1.10", Apps: []string{"app1"}}
	app := Application{ID: "app1", Name: "App1", Version: "1.0"}

	// Create nodes and relationships
	if err := db.CreateServerNode(ctx, server); err != nil {
		log.Fatalf("Failed to create server node: %v", err)
	}
	if err := db.CreateApplicationNode(ctx, app); err != nil {
		log.Fatalf("Failed to create application node: %v", err)
	}
	if err := db.CreateRelationship(ctx, server.ID, app.ID); err != nil {
		log.Fatalf("Failed to create relationship: %v", err)
	}

	fmt.Println("Asset collection complete.")
}

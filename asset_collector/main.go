package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Server struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
	Apps      []string `json:"apps"`
}

type Application struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

func createServerNode(ctx context.Context, driver neo4j.DriverWithContext, server Server) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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

func createApplicationNode(ctx context.Context, driver neo4j.DriverWithContext, app Application) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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

func createRelationship(ctx context.Context, driver neo4j.DriverWithContext, serverID, appID string) error {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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

	// Example data - replace with real CMDB data
	server := Server{ID: "srv1", Name: "Server1", IPAddress: "192.168.1.10", Apps: []string{"app1"}}
	app := Application{ID: "app1", Name: "App1", Version: "1.0"}

	// Create nodes and relationships
	if err := createServerNode(ctx, driver, server); err != nil {
		log.Fatalf("Failed to create server node: %v", err)
	}
	if err := createApplicationNode(ctx, driver, app); err != nil {
		log.Fatalf("Failed to create application node: %v", err)
	}
	if err := createRelationship(ctx, driver, server.ID, app.ID); err != nil {
		log.Fatalf("Failed to create relationship: %v", err)
	}

	fmt.Println("Asset collection complete.")
}

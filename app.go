package main

import (
	"context"
	"fmt"
	"os"
)

// App struct
type App struct {
	ctx             context.Context
	typesenseClient *TypesenseClient
}

// NewApp creates a new App application struct
func NewApp() *App {
	apiKey := os.Getenv("TYPESENSE_API_KEY")
	if apiKey == "" {
		fmt.Println("Warning: TYPESENSE_API_KEY environment variable is not set")
	}

	host := "http://typesense.documentresearch.dev:8080"

	return &App{
		typesenseClient: NewTypesenseClient(apiKey, host),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetCollections retrieves all collections from Typesense
func (a *App) GetCollections() ([]Collection, error) {
	collections, err := a.typesenseClient.GetCollections()
	if err != nil {
		return nil, fmt.Errorf("failed to get collections: %w", err)
	}
	return collections, nil
}

// GetCollectionData retrieves all data from a specified collection
func (a *App) GetCollectionData(collectionName string) (string, error) {
	data, err := a.typesenseClient.GetCollectionData(collectionName)
	if err != nil {
		return "", fmt.Errorf("failed to get collection data: %w", err)
	}
	return data, nil
}

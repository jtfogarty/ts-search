package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

// TypesenseClient represents the Typesense client
type TypesenseClient struct {
	client *typesense.Client
}

// Collection represents a Typesense collection
type Collection struct {
	Name         string  `json:"name"`
	Fields       []Field `json:"fields"`
	NumDocuments int64   `json:"num_documents"`
}

// Field represents a field in a Typesense collection
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// NewTypesenseClient creates a new Typesense client
func NewTypesenseClient(apiKey, host string) *TypesenseClient {
	client := typesense.NewClient(
		typesense.WithServer(host),
		typesense.WithAPIKey(apiKey),
	)
	return &TypesenseClient{client: client}
}

// GetCollections retrieves all collections from the Typesense database
func (tc *TypesenseClient) GetCollections() ([]Collection, error) {
	collections, err := tc.client.Collections().Retrieve(context.Background())
	if err != nil {
		return nil, err
	}

	var result []Collection
	for _, c := range collections {
		fields := make([]Field, len(c.Fields))
		for i, f := range c.Fields {
			fields[i] = Field{
				Name: f.Name,
				Type: f.Type,
			}
		}

		// Get the number of documents for this collection
		collectionInfo, err := tc.client.Collection(c.Name).Retrieve(context.Background())
		if err != nil {
			return nil, err
		}

		var numDocuments int64
		if collectionInfo.NumDocuments != nil {
			numDocuments = *collectionInfo.NumDocuments
			// Log the number of documents for this collection
			//fmt.Printf("Collection '%s' has %d documents\n", c.Name, numDocuments)
		} else {
			fmt.Printf("Collection '%s' has an unknown number of documents\n", c.Name)
		}

		result = append(result, Collection{
			Name:         c.Name,
			Fields:       fields,
			NumDocuments: numDocuments,
		})
	}

	return result, nil
}

// GetCollectionData retrieves all documents from a specified collection and returns them as JSON
func (tc *TypesenseClient) GetCollectionData(collectionName string) (string, error) {
	perPage := int32(1000)
	page := int32(1)

	searchParameters := &api.SearchCollectionParams{
		Q:       "*",
		QueryBy: "",
		PerPage: func() *int { p := int(perPage); return &p }(),
		Page:    func() *int { p := int(page); return &p }(),
	}

	searchResult, err := tc.client.Collection(collectionName).Documents().Search(context.Background(), searchParameters)
	if err != nil {
		return "", fmt.Errorf("failed to search collection: %w", err)
	}

	// Extract the hits from the search result
	hits := searchResult.Hits

	// Convert the hits to a JSON string
	jsonData, err := json.MarshalIndent(hits, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal search results to JSON: %w", err)
	}

	return string(jsonData), nil
}

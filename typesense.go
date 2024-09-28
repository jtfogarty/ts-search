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

const defaultPerPage = 250

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

// ShakespeareWorks represents a document in the shakespeare_works collection
type ShakespeareWorks struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// GetShakespeareWorks retrieves all documents from the shakespeare_works collection
func (tc *TypesenseClient) GetShakespeareWorks() ([]ShakespeareWorks, error) {
	searchParameters := &api.SearchCollectionParams{
		Q:       "*",
		QueryBy: "title",
		PerPage: func() *int { p := int(defaultPerPage); return &p }(),
	}

	searchResult, err := tc.client.Collection("shakespeare_works").Documents().Search(context.Background(), searchParameters)
	if err != nil {
		return nil, fmt.Errorf("error searching shakespeare_works collection: %w", err)
	}

	var works []ShakespeareWorks
	if searchResult.Hits != nil {
		for _, hit := range *searchResult.Hits {
			var work ShakespeareWorks
			docBytes, err := json.Marshal(hit.Document) // Convert to JSON bytes
			if err != nil {
				return nil, fmt.Errorf("error marshaling document: %w", err)
			}
			err = json.Unmarshal(docBytes, &work) // Unmarshal the JSON bytes
			if err != nil {
				return nil, fmt.Errorf("error unmarshaling document: %w", err)
			}
			works = append(works, work)
		}
	}

	return works, nil
}

// GetCollections retrieves all collections from the Typesense database
func (tc *TypesenseClient) GetCollections() ([]Collection, error) {
	collections, err := tc.client.Collections().Retrieve(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve collections: %w", err)
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

		collectionInfo, err := tc.client.Collection(c.Name).Retrieve(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve collection info for '%s': %w", c.Name, err)
		}

		var numDocuments int64
		if collectionInfo.NumDocuments != nil {
			numDocuments = *collectionInfo.NumDocuments
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
	page := int32(1)
	searchParameters := &api.SearchCollectionParams{
		Q:       "*",
		QueryBy: "",
		PerPage: func() *int { p := int(defaultPerPage); return &p }(),
		Page:    func() *int { p := int(page); return &p }(),
	}

	searchResult, err := tc.client.Collection(collectionName).Documents().Search(context.Background(), searchParameters)
	if err != nil {
		return "", fmt.Errorf("failed to search collection: %w", err)
	}

	if searchResult.Hits == nil {
		return "[]", nil
	}

	jsonData, err := json.MarshalIndent(searchResult.Hits, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal search results to JSON: %w", err)
	}

	return string(jsonData), nil
}

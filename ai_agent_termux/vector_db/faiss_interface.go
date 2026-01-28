package vector_db

import (
	"encoding/json"
	"fmt"

	"ai_agent_termux/config"
	"ai_agent_termux/preprocessor"
)

type FaissInterface struct {
	config *config.Config
}

func NewFaissInterface(cfg *config.Config) *FaissInterface {
	return &FaissInterface{
		config: cfg,
	}
}

func (fi *FaissInterface) AddEmbedding(embedding []float64, metadata map[string]interface{}) error {
	// Convert embedding to JSON string
	embeddingJSON, err := json.Marshal(embedding)
	if err != nil {
		return fmt.Errorf("error marshaling embedding: %v", err)
	}

	// Convert metadata to JSON string
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("error marshaling metadata: %v", err)
	}

	// Call Python script to add embedding
	_, err = preprocessor.ManageFaissIndex(
		"add",
		fmt.Sprintf("%d", len(embedding)), // dimension
		string(embeddingJSON),
		string(metadataJSON),
	)

	return err
}

func (fi *FaissInterface) SaveIndex() error {
	_, err := preprocessor.ManageFaissIndex("save")
	return err
}

func (fi *FaissInterface) LoadIndex() error {
	_, err := preprocessor.ManageFaissIndex("load")
	return err
}

func (fi *FaissInterface) Search(queryEmbedding []float64, k int) (map[string]interface{}, error) {
	// Convert query embedding to JSON string
	queryJSON, err := json.Marshal(queryEmbedding)
	if err != nil {
		return nil, fmt.Errorf("error marshaling query embedding: %v", err)
	}

	// Call Python script to search
	resultStr, err := preprocessor.ManageFaissIndex(
		"search",
		fmt.Sprintf("%d", len(queryEmbedding)), // dimension
		string(queryJSON),
	)
	if err != nil {
		return nil, err
	}

	// Parse result
	var result map[string]interface{}
	err = json.Unmarshal([]byte(resultStr), &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling search result: %v", err)
	}

	return result, nil
}

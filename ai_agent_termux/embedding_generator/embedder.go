package embedding_generator

import (
	"ai_agent_termux/config"
	"ai_agent_termux/preprocessor"
)

// GenerateEmbeddingsForFile generates embeddings for preprocessed file content
func GenerateEmbeddingsForFile(content string, cfg *config.Config) ([]float64, error) {
	return preprocessor.GenerateEmbeddings(content, cfg)
}

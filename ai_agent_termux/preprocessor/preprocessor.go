package preprocessor

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"ai_agent_termux/config"
	"ai_agent_termux/file_processor"
	"ai_agent_termux/utils"
)

type PreprocessedContent struct {
	Text     string            `json:"text"`
	Metadata map[string]string `json:"metadata"`
}

func PreprocessFile(file file_processor.FileMetadata, cfg *config.Config) (*PreprocessedContent, error) {
	var content PreprocessedContent
	var err error

	switch file.Type {
	case "image":
		content, err = preprocessImage(file.Path, cfg)
	case "pdf":
		content, err = preprocessPDF(file.Path, cfg)
	case "text":
		content, err = preprocessText(file.Path, cfg)
	default:
		content.Text = "Unsupported file type for preprocessing"
		content.Metadata = make(map[string]string)
		content.Metadata["source"] = "preprocessor"
		content.Metadata["method"] = "none"
	}

	if err != nil {
		return nil, fmt.Errorf("error preprocessing file %s: %v", file.Path, err)
	}

	return &content, nil
}

func preprocessImage(imagePath string, cfg *config.Config) (PreprocessedContent, error) {
	content := PreprocessedContent{
		Metadata: make(map[string]string),
	}
	content.Metadata["source"] = "ocr"
	content.Metadata["language"] = cfg.TesseractLang

	scriptPath := filepath.Join(cfg.PythonScriptsDir, "ocr_processor.py")
	cmd := "python3"
	args := []string{scriptPath, imagePath, cfg.TesseractLang}

	output, err := utils.ExecuteCommand(cmd, args...)
	if err != nil {
		return content, err
	}

	content.Text = output
	return content, nil
}

func preprocessPDF(pdfPath string, cfg *config.Config) (PreprocessedContent, error) {
	content := PreprocessedContent{
		Metadata: make(map[string]string),
	}
	content.Metadata["source"] = "pdf_parser"

	scriptPath := filepath.Join(cfg.PythonScriptsDir, "pdf_processor.py")
	cmd := "python3"
	args := []string{scriptPath, pdfPath}

	output, err := utils.ExecuteCommand(cmd, args...)
	if err != nil {
		return content, err
	}

	content.Text = output
	return content, nil
}

func preprocessText(textPath string, cfg *config.Config) (PreprocessedContent, error) {
	content := PreprocessedContent{
		Metadata: make(map[string]string),
	}
	content.Metadata["source"] = "file_reader"

	// For simplicity, we'll read the file directly in Go
	// In a more complex implementation, we might use a Python script
	data, err := utils.ExecuteCommand("cat", textPath)
	if err != nil {
		return content, err
	}

	content.Text = data
	return content, nil
}

func GenerateEmbeddings(text string, cfg *config.Config) ([]float64, error) {
	scriptPath := filepath.Join(cfg.PythonScriptsDir, "embedding_generator.py")
	cmd := "python3"
	args := []string{scriptPath, text}

	output, err := utils.ExecuteCommand(cmd, args...)
	if err != nil {
		return nil, err
	}

	var embeddings []float64
	err = json.Unmarshal([]byte(output), &embeddings)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling embeddings: %v", err)
	}

	return embeddings, nil
}

func ManageFaissIndex(operation string, args ...string) (string, error) {
	scriptPath := filepath.Join(config.LoadConfig().PythonScriptsDir, "faiss_manager.py")
	cmdArgs := []string{scriptPath}
	cmdArgs = append(cmdArgs, args...)

	return utils.ExecuteCommand("python3", cmdArgs...)
}

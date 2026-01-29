package preprocessor

import (
	"path/filepath"

	"ai_agent_termux/config"
	"ai_agent_termux/utils"
)

// preprocessMedia handles audio/video file metadata extraction
func preprocessMedia(mediaPath string, cfg *config.Config) (PreprocessedContent, error) {
	content := PreprocessedContent{
		Metadata: make(map[string]string),
	}
	content.Metadata["source"] = "media_parser"

	scriptPath := filepath.Join(cfg.PythonScriptsDir, "media_processor.py")
	cmd := "python3"
	args := []string{scriptPath, mediaPath}

	output, err := utils.ExecuteCommand(cmd, args...)
	if err != nil {
		return content, err
	}

	content.Text = output
	return content, nil
}

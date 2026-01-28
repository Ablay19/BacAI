package interactive

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ai_agent_termux/automation"
	"ai_agent_termux/goroutines"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/exp/slog"
)

// CLIConfig contains configuration for interactive CLI mode
type CLIConfig struct {
	GoogleLens       *automation.GoogleLensProcessor
	TaskManager      *goroutines.TaskManager
	HistoryFile      string
	MaxHistory       int
	DefaultOperation string
	ColorEnabled     bool
}

// InteractiveCLI represents the interactive CLI interface
type InteractiveCLI struct {
	config      *CLIConfig
	running     bool
	history     []string
	historyPos  int
	currentCmd  string
	suggestions []string
}

// CLIMode represents different CLI modes
type CLIMode int

const (
	ModeNormal CLIMode = iota
	ModeBatchProcess
	ModeFileWatch
	ModeSettings
)

// TeaModel represents the Bubble Tea model for TUI
type TeaModel struct {
	cli          *InteractiveCLI
	mode         CLIMode
	input        string
	cursor       int
	messages     []string
	selectedFile int
	files        []string
	progress     float64
	processing   bool
	quitting     bool
	width        int
	height       int
}

// DefaultCLIConfig returns default configuration for interactive CLI
func DefaultCLIConfig() *CLIConfig {
	homeDir, _ := os.UserHomeDir()
	return &CLIConfig{
		HistoryFile:      filepath.Join(homeDir, ".google_lens_history"),
		MaxHistory:       1000,
		DefaultOperation: "extract_text",
		ColorEnabled:     true,
	}
}

// NewInteractiveCLI creates a new interactive CLI instance
func NewInteractiveCLI(config *CLIConfig) *InteractiveCLI {
	cli := &InteractiveCLI{
		config:     config,
		running:    true,
		history:    []string{},
		historyPos: -1,
		currentCmd: "",
	}

	// Load history
	cli.loadHistory()

	return cli
}

// Start starts the interactive CLI
func (cli *InteractiveCLI) Start() error {
	slog.Info("Starting interactive CLI mode")

	// Initialize Bubble Tea program
	p := tea.NewProgram(&TeaModel{
		cli:    cli,
		mode:   ModeNormal,
		input:  "",
		cursor: 0,
		messages: []string{
			"Google Lens Interactive CLI",
			"Type 'help' for commands or 'exit' to quit",
		},
		progress:   0,
		processing: false,
	})

	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run interactive CLI: %v", err)
	}

	// Save history
	cli.saveHistory()

	return nil
}

// Bubble Tea initialization
func (m *TeaModel) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
}

// Bubble Tea update function
func (m *TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit

		case tea.KeyEnter:
			if m.processing {
				return m, nil
			}
			return m, cli.processCommand(m.input)

		case tea.KeyUp:
			if m.mode == ModeNormal {
				m.cli.historyPos++
				if m.cli.historyPos >= len(m.cli.history) {
					m.cli.historyPos = len(m.cli.history) - 1
				}
				if m.cli.historyPos >= 0 && m.cli.historyPos < len(m.cli.history) {
					m.input = m.cli.history[len(m.cli.history)-1-m.cli.historyPos]
				}
			}

		case tea.KeyDown:
			if m.mode == ModeNormal {
				m.cli.historyPos--
				if m.cli.historyPos < -1 {
					m.cli.historyPos = -1
				}
				if m.cli.historyPos == -1 {
					m.input = ""
				} else if m.cli.historyPos >= 0 && m.cli.historyPos < len(m.cli.history) {
					m.input = m.cli.history[len(m.cli.history)-1-m.cli.historyPos]
				}
			}

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		case tea.KeyRunes:
			if m.mode == ModeNormal {
				m.input += string(msg.Runes)
			}

		default:
			// Handle mode-specific keys
			switch m.mode {
			case ModeBatchProcess:
				return m, m.handleBatchProcessKeys(msg)
			case ModeFileWatch:
				return m, m.handleFileWatchKeys(msg)
			}
		}
	}

	return m, nil
}

// Bubble Tea view function
func (m *TeaModel) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	var content strings.Builder

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Render("Google Lens Interactive CLI")

	content.WriteString(header + "\n\n")

	// Mode indicator
	modeText := ""
	switch m.mode {
	case ModeNormal:
		modeText = "Normal Mode"
	case ModeBatchProcess:
		modeText = "Batch Processing Mode"
	case ModeFileWatch:
		modeText = "File Watching Mode"
	}

	modeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#04B575")).
		Bold(true)

	content.WriteString(modeStyle.Render("Mode: "+modeText) + "\n\n")

	// Input prompt
	if m.mode == ModeNormal {
		prompt := ">>> "
		if m.processing {
			prompt = "... "
		}
		content.WriteString(prompt + m.input)
	} else {
		// Mode-specific views
		switch m.mode {
		case ModeBatchProcess:
			content.WriteString(m.viewBatchProcess())
		case ModeFileWatch:
			content.WriteString(m.viewFileWatch())
		}
	}

	// Progress bar if processing
	if m.processing {
		bar := progressbar.Default(100)
		bar.Set64(int64(m.progress))
		content.WriteString("\n\nProcessing: ")
		content.WriteString(bar.String())
	}

	// Messages/help
	if len(m.messages) > 0 {
		content.WriteString("\n\n")
		for _, msg := range m.messages {
			content.WriteString(msg + "\n")
		}
	}

	// Help text
	if m.mode == ModeNormal {
		helpText := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#969B86")).
			Render("Commands: help | process <file> | batch <files...> | watch <dir> | history | settings | exit")

		content.WriteString("\n\n" + helpText)
	}

	return content.String()
}

// processCommand processes CLI commands
func (m *TeaModel) processCommand(cmd string) tea.Cmd {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return nil
	}

	// Add to history
	m.cli.addToHistory(cmd)

	// Parse command
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "help", "h":
		m.messages = []string{
			"Available commands:",
			"  process <file> [operation] - Process a single image",
			"  batch <files...> [operation] - Process multiple images",
			"  watch <directory> [operation] - Watch directory for new images",
			"  history - Show command history",
			"  settings - Change settings",
			"  clear - Clear screen",
			"  exit, quit - Exit interactive mode",
			"",
			"Operations: extract_text, identify_object, solve_math, translate_text, barcode, document",
		}

	case "process", "p":
		if len(parts) < 2 {
			m.messages = []string{"Error: Please specify a file to process"}
			return nil
		}
		return m.processSingleFile(parts[1], m.getOperation(parts))

	case "batch", "b":
		if len(parts) < 2 {
			m.messages = []string{"Error: Please specify files to process"}
			return nil
		}
		return m.processBatchFiles(parts[1:], m.getOperation(parts))

	case "watch", "w":
		if len(parts) < 2 {
			m.messages = []string{"Error: Please specify a directory to watch"}
			return nil
		}
		return m.watchDirectory(parts[1], m.getOperation(parts))

	case "history":
		m.messages = m.cli.getHistory()

	case "settings", "config":
		m.mode = ModeSettings
		m.messages = []string{"Settings mode (not implemented yet)"}

	case "clear", "cls":
		return tea.ClearScreen

	case "exit", "quit", "q":
		m.quitting = true
		return tea.Quit

	default:
		m.messages = []string{fmt.Sprintf("Unknown command: %s", parts[0])}
	}

	// Clear input
	m.input = ""
	return nil
}

// getOperation extracts operation from command parts
func (m *TeaModel) getOperation(parts []string) string {
	if len(parts) > 2 {
		return parts[2]
	}
	return m.cli.config.DefaultOperation
}

// processSingleFile processes a single file
func (m *TeaModel) processSingleFile(filename string, operation string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		m.processing = true
		m.progress = 0

		// Progress callback
		progressCB := func(taskID string, progress float64, message string) {
			m.progress = progress
		}

		// Process file
		result, err := m.cli.config.GoogleLens.ProcessImageWithProgress(filename, operation, func(progress float64, message string) {
			m.progress = progress * 100
		})

		m.processing = false

		if err != nil {
			return msgResult{success: false, message: fmt.Sprintf("Error processing %s: %v", filename, err)}
		}

		return msgResult{success: true, message: fmt.Sprintf("Successfully processed %s\nResult: %s", filename, result.ResultText)}
	})
}

// processBatchFiles processes multiple files
func (m *TeaModel) processBatchFiles(filenames []string, operation string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		m.processing = true
		m.progress = 0

		// Progress callback
		results, err := m.cli.config.GoogleLens.BatchProcessImagesWithProgress(filenames, operation, func(progress float64, message string) {
			m.progress = progress * 100
		})

		m.processing = false

		if err != nil {
			return msgResult{success: false, message: fmt.Sprintf("Error in batch processing: %v", err)}
		}

		return msgResult{success: true, message: fmt.Sprintf("Successfully processed %d files", len(results))}
	})
}

// watchDirectory starts watching a directory
func (m *TeaModel) watchDirectory(dirname string, operation string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		m.mode = ModeFileWatch
		return msgResult{success: true, message: fmt.Sprintf("Started watching directory: %s", dirname)}
	})
}

// addToHistory adds command to history
func (cli *InteractiveCLI) addToHistory(cmd string) {
	cli.history = append(cli.history, cmd)
	if len(cli.history) > cli.config.MaxHistory {
		cli.history = cli.history[1:]
	}
	cli.historyPos = -1
}

// getHistory returns command history
func (cli *InteractiveCLI) getHistory() []string {
	if len(cli.history) == 0 {
		return []string{"No command history yet"}
	}

	history := []string{"Command History:"}
	for i, cmd := range cli.history {
		history = append(history, fmt.Sprintf("% 3d: %s", i+1, cmd))
	}
	return history
}

// loadHistory loads command history from file
func (cli *InteractiveCLI) loadHistory() {
	if cli.config.HistoryFile == "" {
		return
	}

	file, err := os.Open(cli.config.HistoryFile)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cmd := scanner.Text()
		if cmd != "" {
			cli.history = append(cli.history, cmd)
		}
	}
}

// saveHistory saves command history to file
func (cli *InteractiveCLI) saveHistory() {
	if cli.config.HistoryFile == "" || len(cli.history) == 0 {
		return
	}

	file, err := os.Create(cli.config.HistoryFile)
	if err != nil {
		slog.Warn("Failed to save history", "error", err)
		return
	}
	defer file.Close()

	for _, cmd := range cli.history {
		fmt.Fprintln(file, cmd)
	}
}

// viewBatchProcess renders batch processing view
func (m *TeaModel) viewBatchProcess() string {
	return "Batch processing view (not implemented yet)\n"
}

// viewFileWatch renders file watching view
func (m *TeaModel) viewFileWatch() string {
	return "File watching view (not implemented yet)\n"
}

// handleBatchProcessKeys handles keys in batch processing mode
func (m *TeaModel) handleBatchProcessKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.mode = ModeNormal
		m.input = ""
	}
	return m, nil
}

// handleFileWatchKeys handles keys in file watching mode
func (m *TeaModel) handleFileWatchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.mode = ModeNormal
		m.input = ""
	}
	return m, nil
}

// msgResult represents a processing result message
type msgResult struct {
	success bool
	message string
}

// BatchProcessWithProgress provides a simple batch processing function with progress bar
func BatchProcessWithProgress(googleLens *automation.GoogleLensProcessor, files []string, operation string) error {
	bar := progressbar.Default(int64(len(files)), "Processing files")
	bar.Describe("Initializing...")

	for i, file := range files {
		bar.ChangeMax(len(files))
		bar.Set(i)
		bar.Describe(fmt.Sprintf("Processing %s", filepath.Base(file)))

		_, err := googleLens.ProcessImage(file, operation)
		if err != nil {
			bar.Describe(fmt.Sprintf("Error processing %s", filepath.Base(file)))
			return err
		}
	}

	bar.Set(len(files))
	bar.Describe("Batch processing completed!")
	return nil
}

// InteractiveProgress creates an interactive progress indicator
type InteractiveProgress struct {
	bar    *progressbar.ProgressBar
	cancel context.CancelFunc
}

// NewInteractiveProgress creates a new interactive progress indicator
func NewInteractiveProgress(total int64, description string) *InteractiveProgress {
	bar := progressbar.Default(total)
	bar.Describe(description)

	return &InteractiveProgress{
		bar: bar,
	}
}

// Update updates the progress
func (ip *InteractiveProgress) Update(current int64, description string) {
	ip.bar.Set64(current)
	if description != "" {
		ip.bar.Describe(description)
	}
}

// Finish completes the progress bar
func (ip *InteractiveProgress) Finish(description string) {
	ip.bar.Describe(description)
}

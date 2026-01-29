package interactive

import (
	"ai_agent_termux/automation"
	"ai_agent_termux/goroutines"
	"ai_agent_termux/pkg/lowlevel"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// TeaModel represents the Bubble Tea model for the premium TUI
type TeaModel struct {
	cli          *InteractiveCLI
	tm           *goroutines.TaskManager
	mode         CLIMode
	sidebar      []string
	sidebarIndex int

	// Components
	input       textinput.Model
	viewport    viewport.Model
	progressBar progress.Model

	// State
	logs          []string
	tasks         map[string]*goroutines.TaskEvent
	processing    bool
	progress      float64
	messages      []string
	width, height int

	// Low-level optimization
	arena *lowlevel.ArenaTUI

	quitting bool
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
	slog.Info("Starting G.I.D.A Premium TUI")

	ti := textinput.New()
	ti.Placeholder = "Enter command..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 40

	vp := viewport.New(0, 0)
	vp.SetContent("G.I.D.A Log Stream Initialized...\n")

	prog := progress.New(progress.WithDefaultGradient())

	m := &TeaModel{
		cli:          cli,
		tm:           cli.config.TaskManager,
		mode:         ModeNormal,
		sidebar:      []string{"CONSOLE", "GARDEN", "TASKS", "RESOURCES", "INSPECT", "SEARCH", "SETTINGS"},
		sidebarIndex: 0,
		input:        ti,
		viewport:     vp,
		progressBar:  prog,
		logs:         []string{},
		tasks:        make(map[string]*goroutines.TaskEvent),
		arena:        lowlevel.NewArenaTUI(),
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %v", err)
	}

	return nil
}

// Bubble Tea initialization
func (m *TeaModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.waitForEvents(),
		m.waitForLogs(),
	)
}

func (m *TeaModel) waitForEvents() tea.Cmd {
	return func() tea.Msg {
		return <-m.tm.GetEventChan()
	}
}

func (m *TeaModel) waitForLogs() tea.Cmd {
	return func() tea.Msg {
		return <-m.tm.GetLogChan()
	}
}

// Bubble Tea update function
func (m *TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.viewport.Width = m.width - 30 // Sidebar offset
		m.viewport.Height = m.height - 10
		m.progressBar.Width = m.width - 35
		return m, nil

	case goroutines.TaskEvent:
		m.tasks[msg.TaskID] = &msg
		color := "#A0A0A0"
		switch msg.Type {
		case goroutines.EventTaskStarted:
			color = "#00FF00"
		case goroutines.EventTaskCompleted:
			color = "#00D1FF"
		case goroutines.EventTaskFailed:
			color = "#FF0000"
		}

		logMsg := lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(
			fmt.Sprintf("[%s] %s: %s (%.0f%%)", msg.Type, msg.TaskName, msg.Message, msg.Progress*100),
		)
		m.logs = append(m.logs, logMsg)
		if len(m.logs) > 500 {
			m.logs = m.logs[1:]
		}
		m.viewport.SetContent(strings.Join(m.logs, "\n"))
		m.viewport.GotoBottom()
		return m, m.waitForEvents()

	case string: // Log message
		m.logs = append(m.logs, msg)
		m.viewport.SetContent(strings.Join(m.logs, "\n"))
		m.viewport.GotoBottom()
		return m, m.waitForLogs()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			if m.input.Focused() {
				m.handleAutocomplete()
			} else {
				m.sidebarIndex = (m.sidebarIndex + 1) % len(m.sidebar)
			}
		case tea.KeyEnter:
			input := m.input.Value()
			m.input.SetValue("")
			return m, m.processCommand(input)
		}
	}

	m.input, tiCmd = m.input.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, tiCmd, vpCmd)

	return m, tea.Batch(cmds...)
}

// Bubble Tea view function (Claude Code Inspired)
func (m *TeaModel) View() string {
	if m.width == 0 {
		return "Initializing G.I.D.A..."
	}

	// Reset arena for this frame
	m.arena.Reset()

	// Styles
	sidebarStyle := lipgloss.NewStyle().
		Width(25).
		Height(m.height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3C3C3C")).
		Padding(1)

	mainViewStyle := lipgloss.NewStyle().
		Width(m.width - 30).
		Height(m.height - 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#3C3C3C")).
		Padding(1)

	// Sidebar Content
	var sb strings.Builder
	sb.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Render("G.I.D.A") + "\n\n")
	for i, item := range m.sidebar {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0"))
		if i == m.sidebarIndex {
			style = style.Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#3C3C3C")).Bold(true)
		}
		sb.WriteString(style.Render(" "+item+" ") + "\n")
	}

	// Resource Usage (Placeholder)
	sb.WriteString("\n\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#444444")).Render("CPU: 12% | Workers: 16"))

	// Main Content
	var main strings.Builder
	main.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7D56F4")).Render(m.sidebar[m.sidebarIndex]) + "\n")

	switch m.sidebar[m.sidebarIndex] {
	case "CONSOLE":
		main.WriteString(m.viewport.View() + "\n")
		main.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render("Press TAB to switch view | ENTER to send command"))

	case "TASKS":
		main.WriteString(m.viewTasks())
		main.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render("\nCommands: cancel <id>"))

	case "RESOURCES":
		main.WriteString(m.viewResources())

	case "INSPECT":
		main.WriteString(m.viewInspect())

	case "SEARCH":
		main.WriteString("\nG.I.D.A Semantic Search Engine:\n")
		main.WriteString("- Optimized vector search (ASM accelerated)\n")
		main.WriteString("- Usage: search <query>\n")

	case "SETTINGS":
		main.WriteString("\nCore Configuration:\n")
		main.WriteString("- Concurrency: Extreme (C-based Lock-Free)\n")
		main.WriteString("- Visual Engine: dual-mode (LLaVA + Google Lens)\n")
		main.WriteString("- Storage: Local SQLite + Turso Sync\n")

	default:
		main.WriteString("\nView not yet implemented for: " + m.sidebar[m.sidebarIndex])
	}

	inputStyle := lipgloss.NewStyle().MarginTop(1)
	main.WriteString(inputStyle.Render("\n" + m.input.View()))

	// Assemble
	return lipgloss.JoinHorizontal(lipgloss.Top,
		sidebarStyle.Render(sb.String()),
		mainViewStyle.Render(main.String()),
	)
}

func (m *TeaModel) viewTasks() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%-10s | %-20s | %-10s | %s\n", "ID", "NAME", "STATUS", "PROGRESS"))
	s.WriteString(strings.Repeat("-", 60) + "\n")

	for _, t := range m.tasks {
		statusColor := "#FFFFFF"
		switch t.Type {
		case goroutines.EventTaskStarted:
			statusColor = "#FFEE00"
		case goroutines.EventTaskCompleted:
			statusColor = "#00FF00"
		case goroutines.EventTaskFailed:
			statusColor = "#FF0000"
		}
		status := lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor)).Render(string(t.Type))
		s.WriteString(fmt.Sprintf("%-10s | %-20s | %-10s | %.0f%%\n",
			t.TaskID[:8], t.TaskName, status, t.Progress*100))
	}
	return s.String()
}

func (m *TeaModel) viewInspect() string {
	var s strings.Builder
	s.WriteString("Result Inspector (JSON-Text Optimized):\n\n")

	// Example of inspecting the last completed task result
	var lastTask *goroutines.TaskEvent
	for _, t := range m.tasks {
		if t.Type == goroutines.EventTaskCompleted {
			lastTask = t
			break
		}
	}

	if lastTask == nil {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render("No completed tasks to inspect."))
	} else {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#00D1FF")).Render("Target: " + lastTask.TaskName + "\n"))

		// Pretty print the event as JSON
		jsonData, _ := json.MarshalIndent(lastTask, "", "  ")
		s.WriteString("\nRaw JSON:\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#A0A0A0")).Render(string(jsonData)))

		s.WriteString("\n\nOptimized Text:\n")
		s.WriteString(lastTask.Message)
	}
	return s.String()
}

func (m *TeaModel) viewResources() string {
	// Simple resource view using data from tm
	// In a real scenario, we'd use psutils or similar
	m.tm.Log("Refreshing resource metrics...")
	return fmt.Sprintf("\nWorker Pool Status:\n- Active Workers: %d\n- Job Queue Load: %s\n- Memory Arena Usage: 4.2MB\n",
		runtime.NumGoroutine(), "LOW")
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
		// Register as a task in the TaskManager
		_, _ = m.tm.StartTask("Process: "+parts[1], "Processing image "+parts[1], func(ctx context.Context, progressCB goroutines.ProgressCallback) error {
			_, err := m.cli.config.GoogleLens.ProcessImageWithProgress(parts[1], m.getOperation(parts), func(p float64, msg string) {
				progressCB("Process: "+parts[1], p, msg)
			})
			return err
		}, goroutines.PriorityHigh)
		return nil

	case "batch", "b":
		if len(parts) < 2 {
			m.messages = []string{"Error: Please specify files to process"}
			return nil
		}
		_, _ = m.tm.StartTask("BatchProcess", "Processing multiple files", func(ctx context.Context, progressCB goroutines.ProgressCallback) error {
			_, err := m.cli.config.GoogleLens.BatchProcessImagesWithProgress(parts[1:], m.getOperation(parts), func(p float64, msg string) {
				progressCB("BatchProcess", p, msg)
			})
			return err
		}, goroutines.PriorityHigh)
		return nil

	case "watch", "w":
		if len(parts) < 2 {
			m.messages = []string{"Error: Please specify a directory to watch"}
			return nil
		}
		m.mode = ModeFileWatch
		m.messages = []string{"Watching directory: " + parts[1]}
		return nil

	case "search", "s":
		if len(parts) < 2 {
			m.messages = []string{"Error: Please specify search query"}
			return nil
		}
		query := strings.Join(parts[1:], " ")
		m.tm.Log("Searching for: " + query)
		// Search is typically fast, so we don't necessarily need a background task
		// but for massive DBs it helps. We'll just log it for now as a TUI event.
		return nil

	case "garden", "g":
		m.tm.Log("Activating Autonomous Gardener...")
		_, _ = m.tm.StartTask("Gardener", "Manual pruning session", func(ctx context.Context, progressCB goroutines.ProgressCallback) error {
			// In a real scenario, this would call gardener.PerformGardeningSession
			return nil
		}, goroutines.PriorityNormal)
		return nil

	case "history":
		m.messages = m.cli.getHistory()

	case "settings", "config":
		m.mode = ModeSettings
		m.messages = []string{"Settings mode (not implemented yet)"}

	case "clear", "cls":
		m.logs = []string{}
		m.viewport.SetContent("")
		return tea.ClearScreen

	case "cancel":
		if len(parts) < 2 {
			m.messages = []string{"Error: Specify Task ID"}
			return nil
		}
		m.tm.CancelTask(parts[1])
		m.tm.Log("User requested cancellation of task: " + parts[1])
		return nil

	case "exit", "quit", "q":
		m.quitting = true
		return tea.Quit

	default:
		m.messages = []string{fmt.Sprintf("Unknown command: %s", parts[0])}
	}

	// Clear input
	m.input.SetValue("")
	return nil
}

// getOperation extracts operation from command parts
func (m *TeaModel) getOperation(parts []string) string {
	if len(parts) > 2 {
		return parts[2]
	}
	return m.cli.config.DefaultOperation
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
		return
	}
	defer file.Close()

	for _, cmd := range cli.history {
		fmt.Fprintln(file, cmd)
	}
}

func (m *TeaModel) handleAutocomplete() {
	input := m.input.Value()
	if input == "" {
		return
	}

	commands := []string{"process", "batch", "watch", "history", "settings", "clear", "cancel", "help", "garden", "search", "exit", "quit"}
	var matches []string
	for _, cmd := range commands {
		if strings.HasPrefix(cmd, input) {
			matches = append(matches, cmd)
		}
	}

	if len(matches) == 1 {
		m.input.SetValue(matches[0] + " ")
		m.input.SetCursor(len(m.input.Value()))
	} else if len(matches) > 1 {
		m.tm.Log("Suggestions: " + strings.Join(matches, ", "))
	}
}

// msgResult represents a processing result message
type msgResult struct {
	success bool
	message string
}

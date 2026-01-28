package goroutines

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slog"
)

// TaskStatus represents the status of a background task
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
	StatusCancelled TaskStatus = "cancelled"
)

// TaskPriority represents task priority levels
type TaskPriority int

const (
	PriorityLow    TaskPriority = 1
	PriorityNormal TaskPriority = 5
	PriorityHigh   TaskPriority = 10
	PriorityUrgent TaskPriority = 20
)

// Task represents a background task
type Task struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Function    TaskFunc           `json:"-"`
	Priority    TaskPriority       `json:"priority"`
	Status      TaskStatus         `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	StartedAt   *time.Time         `json:"started_at,omitempty"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
	Error       string             `json:"error,omitempty"`
	Progress    float64            `json:"progress"`
	MaxRetries  int                `json:"max_retries"`
	Retries     int                `json:"retries"`
	Timeout     time.Duration      `json:"timeout"`
	CancelFunc  context.CancelFunc `json:"-"`
	ProgressCB  ProgressCallback   `json:"-"`
}

// TaskFunc represents the function to execute for a task
type TaskFunc func(ctx context.Context, progressCB ProgressCallback) error

// ProgressCallback represents a progress update callback
type ProgressCallback func(taskID string, progress float64, message string)

// TaskManager manages background goroutines and tasks
type TaskManager struct {
	tasks         map[string]*Task
	tasksMutex    sync.RWMutex
	workers       int
	queue         chan *Task
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	maxConcurrent int
	running       int32
	completed     int64
	failed        int64
	progressCB    ProgressCallback
	runningPool   []*Task
	poolMutex     sync.RWMutex
}

// ManagerConfig contains configuration for the task manager
type ManagerConfig struct {
	MaxConcurrent    int              `json:"max_concurrent"`
	QueueSize        int              `json:"queue_size"`
	DefaultTimeout   time.Duration    `json:"default_timeout"`
	ProgressCallback ProgressCallback `json:"-"`
}

// DefaultManagerConfig returns default configuration for the task manager
func DefaultManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		MaxConcurrent:  runtime.NumCPU(),
		QueueSize:      1000,
		DefaultTimeout: 30 * time.Minute,
	}
}

// NewTaskManager creates a new task manager instance
func NewTaskManager(config *ManagerConfig) *TaskManager {
	if config == nil {
		config = DefaultManagerConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	tm := &TaskManager{
		tasks:         make(map[string]*Task),
		workers:       config.MaxConcurrent,
		queue:         make(chan *Task, config.QueueSize),
		ctx:           ctx,
		cancel:        cancel,
		maxConcurrent: config.MaxConcurrent,
		progressCB:    config.ProgressCallback,
		runningPool:   make([]*Task, 0, config.MaxConcurrent),
	}

	// Start worker pool
	tm.startWorkers()

	slog.Info("Task manager initialized", "workers", tm.workers, "queue_size", config.QueueSize)
	return tm
}

// StartTask starts a new background task
func (tm *TaskManager) StartTask(name, description string, function TaskFunc, priority TaskPriority) (string, error) {
	return tm.StartTaskWithConfig(name, description, function, priority, 0, 3, 0)
}

// StartTaskWithConfig starts a new task with detailed configuration
func (tm *TaskManager) StartTaskWithConfig(name, description string, function TaskFunc, priority TaskPriority, timeout time.Duration, maxRetries int, initialProgress float64) (string, error) {
	taskID := generateTaskID(name)

	// Set default timeout if not specified
	if timeout == 0 {
		timeout = 30 * time.Minute
	}

	task := &Task{
		ID:          taskID,
		Name:        name,
		Description: description,
		Function:    function,
		Priority:    priority,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
		Progress:    initialProgress,
		MaxRetries:  maxRetries,
		Timeout:     timeout,
		ProgressCB:  tm.progressCB,
	}

	// Add to tasks map
	tm.tasksMutex.Lock()
	tm.tasks[taskID] = task
	tm.tasksMutex.Unlock()

	// Add to queue
	select {
	case tm.queue <- task:
		slog.Info("Task queued", "id", taskID, "name", name, "priority", priority)
		return taskID, nil
	default:
		// Queue is full, remove task from map
		tm.tasksMutex.Lock()
		delete(tm.tasks, taskID)
		tm.tasksMutex.Unlock()
		return "", fmt.Errorf("task queue is full")
	}
}

// CancelTask cancels a running or pending task
func (tm *TaskManager) CancelTask(taskID string) error {
	tm.tasksMutex.Lock()
	defer tm.tasksMutex.Unlock()

	task, exists := tm.tasks[taskID]
	if !exists {
		return fmt.Errorf("task not found: %s", taskID)
	}

	if task.Status == StatusCompleted || task.Status == StatusFailed || task.Status == StatusCancelled {
		return fmt.Errorf("task cannot be cancelled, current status: %s", task.Status)
	}

	if task.CancelFunc != nil {
		task.CancelFunc()
	}

	task.Status = StatusCancelled
	now := time.Now()
	task.CompletedAt = &now
	task.Error = "Task cancelled by user"

	// Remove from running pool if it's there
	tm.removeFromRunningPool(taskID)

	slog.Info("Task cancelled", "id", taskID, "name", task.Name)
	return nil
}

// GetTask returns information about a specific task
func (tm *TaskManager) GetTask(taskID string) (*Task, error) {
	tm.tasksMutex.RLock()
	defer tm.tasksMutex.RUnlock()

	task, exists := tm.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	// Return a copy to avoid race conditions
	taskCopy := *task
	return &taskCopy, nil
}

// GetAllTasks returns all tasks
func (tm *TaskManager) GetAllTasks() map[string]*Task {
	tm.tasksMutex.RLock()
	defer tm.tasksMutex.RUnlock()

	result := make(map[string]*Task)
	for id, task := range tm.tasks {
		taskCopy := *task
		result[id] = &taskCopy
	}
	return result
}

// GetRunningTasks returns currently running tasks
func (tm *TaskManager) GetRunningTasks() []*Task {
	tm.poolMutex.RLock()
	defer tm.poolMutex.RUnlock()

	result := make([]*Task, len(tm.runningPool))
	for i, task := range tm.runningPool {
		taskCopy := *task
		result[i] = &taskCopy
	}
	return result
}

// GetPendingTasks returns pending tasks in the queue
func (tm *TaskManager) GetPendingTasks() []*Task {
	tm.tasksMutex.RLock()
	defer tm.tasksMutex.RUnlock()

	var pending []*Task
	for _, task := range tm.tasks {
		if task.Status == StatusPending {
			taskCopy := *task
			pending = append(pending, &taskCopy)
		}
	}
	return pending
}

// GetStats returns statistics about the task manager
func (tm *TaskManager) GetStats() map[string]interface{} {
	tm.tasksMutex.RLock()
	totalTasks := len(tm.tasks)
	tm.tasksMutex.RUnlock()

	running := atomic.LoadInt32(&tm.running)
	completed := atomic.LoadInt64(&tm.completed)
	failed := atomic.LoadInt64(&tm.failed)

	return map[string]interface{}{
		"total_tasks":     totalTasks,
		"running_tasks":   running,
		"pending_tasks":   len(tm.queue),
		"completed_tasks": completed,
		"failed_tasks":    failed,
		"max_workers":     tm.maxConcurrent,
		"queue_size":      cap(tm.queue),
		"queue_available": cap(tm.queue) - len(tm.queue),
	}
}

// startWorkers starts the worker pool
func (tm *TaskManager) startWorkers() {
	for i := 0; i < tm.workers; i++ {
		tm.wg.Add(1)
		go tm.worker(i)
	}
	slog.Info("Worker pool started", "workers", tm.workers)
}

// worker processes tasks from the queue
func (tm *TaskManager) worker(id int) {
	defer tm.wg.Done()
	slog.Debug("Worker started", "worker_id", id)

	for {
		select {
		case <-tm.ctx.Done():
			slog.Debug("Worker stopping", "worker_id", id)
			return

		case task := <-tm.queue:
			tm.executeTask(task, id)
		}
	}
}

// executeTask executes a single task
func (tm *TaskManager) executeTask(task *Task, workerID int) {
	atomic.AddInt32(&tm.running, 1)
	defer atomic.AddInt32(&tm.running, -1)

	// Update task status
	task.Status = StatusRunning
	now := time.Now()
	task.StartedAt = &now

	// Add to running pool
	tm.addToRunningPool(task)

	// Create context with timeout and cancellation
	ctx, cancel := context.WithTimeout(tm.ctx, task.Timeout)
	task.CancelFunc = cancel
	defer cancel()

	slog.Info("Task started", "id", task.ID, "name", task.Name, "worker_id", workerID)

	// Execute the task function
	err := task.Function(ctx, func(taskID string, progress float64, message string) {
		tm.updateTaskProgress(taskID, progress, message)
	})

	// Update task status based on result
	completeNow := time.Now()
	task.CompletedAt = &completeNow

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			task.Status = StatusFailed
			task.Error = fmt.Sprintf("Task timed out after %v", task.Timeout)
			atomic.AddInt64(&tm.failed, 1)
		} else if ctx.Err() == context.Canceled {
			task.Status = StatusCancelled
			task.Error = "Task was cancelled"
		} else {
			task.Status = StatusFailed
			task.Error = err.Error()
			atomic.AddInt64(&tm.failed, 1)

			// Check if we should retry
			if task.Retries < task.MaxRetries {
				task.Retries++
				task.Status = StatusPending
				task.Error = fmt.Sprintf("Retry %d/%d: %v", task.Retries, task.MaxRetries, err)

				// Reset for retry
				task.StartedAt = nil
				task.CompletedAt = nil

				// Add back to queue
				select {
				case tm.queue <- task:
					slog.Info("Task queued for retry", "id", task.ID, "retry", task.Retries, "max_retries", task.MaxRetries)
				default:
					task.Status = StatusFailed
					task.Error = fmt.Sprintf("Failed to queue for retry: %v", err)
				}
			}
		}
	} else {
		task.Status = StatusCompleted
		task.Progress = 100.0
		atomic.AddInt64(&tm.completed, 1)
	}

	// Remove from running pool
	tm.removeFromRunningPool(task.ID)

	slog.Info("Task completed",
		"id", task.ID,
		"name", task.Name,
		"status", task.Status,
		"worker_id", workerID,
		"duration", time.Since(*task.StartedAt))
}

// updateTaskProgress updates the progress of a task
func (tm *TaskManager) updateTaskProgress(taskID string, progress float64, message string) {
	tm.tasksMutex.Lock()
	defer tm.tasksMutex.Unlock()

	task, exists := tm.tasks[taskID]
	if exists && task.Status == StatusRunning {
		task.Progress = progress

		// Call global progress callback if set
		if tm.progressCB != nil {
			tm.progressCB(taskID, progress, message)
		}

		// Call task-specific progress callback if set
		if task.ProgressCB != nil {
			task.ProgressCB(taskID, progress, message)
		}
	}
}

// addToRunningPool adds a task to the running pool
func (tm *TaskManager) addToRunningPool(task *Task) {
	tm.poolMutex.Lock()
	defer tm.poolMutex.Unlock()
	tm.runningPool = append(tm.runningPool, task)
}

// removeFromRunningPool removes a task from the running pool
func (tm *TaskManager) removeFromRunningPool(taskID string) {
	tm.poolMutex.Lock()
	defer tm.poolMutex.Unlock()

	for i, task := range tm.runningPool {
		if task.ID == taskID {
			tm.runningPool = append(tm.runningPool[:i], tm.runningPool[i+1:]...)
			break
		}
	}
}

// CleanupOldTasks removes completed tasks older than the specified duration
func (tm *TaskManager) CleanupOldTasks(maxAge time.Duration) int {
	tm.tasksMutex.Lock()
	defer tm.tasksMutex.Unlock()

	cutoff := time.Now().Add(-maxAge)
	removed := 0

	for id, task := range tm.tasks {
		if (task.Status == StatusCompleted || task.Status == StatusFailed || task.Status == StatusCancelled) &&
			task.CompletedAt != nil && task.CompletedAt.Before(cutoff) {
			delete(tm.tasks, id)
			removed++
		}
	}

	slog.Info("Cleaned up old tasks", "removed", removed, "max_age", maxAge)
	return removed
}

// Shutdown gracefully shuts down the task manager
func (tm *TaskManager) Shutdown(timeout time.Duration) error {
	slog.Info("Shutting down task manager", "timeout", timeout)

	// Cancel all operations
	tm.cancel()

	// Wait for workers to finish or timeout
	done := make(chan struct{})
	go func() {
		tm.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("Task manager shutdown complete")
		return nil
	case <-time.After(timeout):
		slog.Warn("Task manager shutdown timeout")
		return fmt.Errorf("shutdown timeout after %v", timeout)
	}
}

// generateTaskID generates a unique task ID
func generateTaskID(name string) string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d", name, timestamp)
}

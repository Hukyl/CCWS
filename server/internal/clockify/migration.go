package clockify

// PERSONAL SCRIPT - NOT FOR GENERAL USE
// This migration service is specifically designed for my personal Clockify workspace migration needs.
// It parses task names in format "<project>/TASK<number>" and reorganizes them into a new workspace
// structure with proper client/project/task hierarchy. This is NOT a general-purpose migration tool
// and should not be used for other Clockify migration scenarios without significant modifications.

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

// MigrationConfig holds all configuration parameters for the migration
type MigrationConfig struct {
	// Source configuration
	SourceWorkspaceName string `json:"sourceWorkspaceName"`
	SourceProjectName   string `json:"sourceProjectName"`

	// Target configuration
	TargetWorkspaceName string `json:"targetWorkspaceName"`

	// Client mapping (optional - if empty, will create clients based on project names)
	ClientMapping map[string]string `json:"clientMapping,omitempty"`

	// Default client name for projects without specific mapping
	DefaultClientName string `json:"defaultClientName"`

	// Migration options
	DryRun        bool `json:"dryRun"`        // If true, only log what would be done
	BatchSize     int  `json:"batchSize"`     // Number of time entries to process at once
	SkipExisting  bool `json:"skipExisting"`  // Skip if target already has time entries
	CreateClients bool `json:"createClients"` // Whether to create new clients automatically
}

// MigrationStats tracks progress and results
type MigrationStats struct {
	TimeEntriesProcessed int
	TimeEntriesCreated   int
	ProjectsCreated      int
	TasksCreated         int
	ClientsCreated       int
	Errors               []string
	StartTime            time.Time
	EndTime              time.Time
}

// ProjectTaskMapping represents the parsed task information
type ProjectTaskMapping struct {
	OriginalTaskName string
	ProjectName      string
	TaskNumber       string
	NewTaskName      string
	ClientName       string
}

// MigrationService handles the workspace migration process
type MigrationService struct {
	client *APIClient
	config *MigrationConfig
	stats  *MigrationStats

	// Caches to avoid repeated API calls
	sourceWorkspace *Workspace
	targetWorkspace *Workspace
	sourceProject   *Project
	targetProjects  map[string]*Project // projectName -> Project
	targetTasks     map[string]*Task    // projectName/taskName -> Task
	targetClients   map[string]*Client  // clientName -> Client
	currentUser     *User
}

// NewMigrationService creates a new migration service with dependency injection
func NewMigrationService(client *APIClient, config *MigrationConfig) *MigrationService {
	if config.BatchSize <= 0 {
		config.BatchSize = 50 // Default batch size
	}

	if config.DefaultClientName == "" {
		config.DefaultClientName = "Default Client"
	}

	return &MigrationService{
		client:         client,
		config:         config,
		stats:          &MigrationStats{StartTime: time.Now()},
		targetProjects: make(map[string]*Project),
		targetTasks:    make(map[string]*Task),
		targetClients:  make(map[string]*Client),
	}
}

// ExecuteMigration runs the complete migration process
func (m *MigrationService) ExecuteMigration() (*MigrationStats, error) {
	log.Printf("Starting migration from %s/%s to %s",
		m.config.SourceWorkspaceName, m.config.SourceProjectName, m.config.TargetWorkspaceName)

	// Step 1: Initialize workspaces and cache data
	if err := m.initializeWorkspaces(); err != nil {
		return m.stats, fmt.Errorf("failed to initialize workspaces: %w", err)
	}

	// Step 2: Get source time entries
	timeEntries, err := m.client.GetProjectTimeEntries(m.sourceWorkspace.ID, m.sourceProject.ID, m.currentUser.ID)
	if err != nil {
		return m.stats, fmt.Errorf("failed to get source time entries: %w", err)
	}

	log.Printf("Found %d time entries to migrate", len(timeEntries))

	// Step 3: Process time entries in batches
	if err := m.processTimeEntries(timeEntries); err != nil {
		return m.stats, fmt.Errorf("failed to process time entries: %w", err)
	}

	m.stats.EndTime = time.Now()
	m.logMigrationSummary()

	return m.stats, nil
}

// initializeWorkspaces sets up source and target workspaces
func (m *MigrationService) initializeWorkspaces() error {
	// Get current user
	user, err := m.client.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	m.currentUser = user

	// Find source workspace
	sourceWs, err := m.client.FindWorkspaceByName(m.config.SourceWorkspaceName)
	if err != nil {
		return fmt.Errorf("failed to find source workspace '%s': %w", m.config.SourceWorkspaceName, err)
	}
	m.sourceWorkspace = sourceWs

	// Find source project
	sourceProj, err := m.client.FindProjectByName(sourceWs.ID, m.config.SourceProjectName)
	if err != nil {
		return fmt.Errorf("failed to find source project '%s': %w", m.config.SourceProjectName, err)
	}
	m.sourceProject = sourceProj

	// Get or create target workspace
	targetWs, err := m.getOrCreateTargetWorkspace()
	if err != nil {
		return fmt.Errorf("failed to get/create target workspace: %w", err)
	}
	m.targetWorkspace = targetWs

	// Cache existing target clients
	if err := m.cacheTargetClients(); err != nil {
		return fmt.Errorf("failed to cache target clients: %w", err)
	}

	return nil
}

// getOrCreateTargetWorkspace gets existing or creates new target workspace
func (m *MigrationService) getOrCreateTargetWorkspace() (*Workspace, error) {
	// Try to find existing workspace first
	ws, err := m.client.FindWorkspaceByName(m.config.TargetWorkspaceName)
	if err == nil {
		log.Printf("Using existing target workspace: %s", ws.Name)
		return ws, nil
	}

	// Note: Workspace creation might not be available in free tier
	// For now, we'll require the target workspace to exist
	return nil, fmt.Errorf("target workspace '%s' not found - please create it manually first", m.config.TargetWorkspaceName)
}

// cacheTargetClients loads existing clients in target workspace
func (m *MigrationService) cacheTargetClients() error {

	for clients, err := range m.client.IterClients(m.targetWorkspace.ID) {
		if err != nil {
			return err
		}

		for _, client := range clients {
			m.targetClients[client.Name] = &client
		}
	}

	log.Printf("Cached %d existing clients in target workspace", len(m.targetClients))
	return nil
}

// processTimeEntries processes all time entries in batches
func (m *MigrationService) processTimeEntries(timeEntries []TimeEntry) error {
	for i := 0; i < len(timeEntries); i += m.config.BatchSize {
		end := i + m.config.BatchSize
		end = min(end, len(timeEntries))

		batch := timeEntries[i:end]
		log.Printf("Processing batch %d-%d of %d time entries", i+1, end, len(timeEntries))

		if err := m.processBatch(batch); err != nil {
			return fmt.Errorf("failed to process batch %d-%d: %w", i+1, end, err)
		}
	}

	return nil
}

// processBatch processes a batch of time entries
func (m *MigrationService) processBatch(timeEntries []TimeEntry) error {
	for _, entry := range timeEntries {
		if err := m.processTimeEntry(&entry); err != nil {
			m.stats.Errors = append(m.stats.Errors, fmt.Sprintf("Failed to process entry %s: %v", entry.ID, err))
			log.Printf("Error processing time entry %s: %v", entry.ID, err)
			continue
		}
		m.stats.TimeEntriesProcessed++
	}

	return nil
}

// processTimeEntry processes a single time entry
func (m *MigrationService) processTimeEntry(entry *TimeEntry) error {
	// Get the task information to parse project/task names
	task, err := m.getSourceTask(entry.TaskID)
	if err != nil {
		return fmt.Errorf("failed to get source task: %w", err)
	}

	// Parse the task name to extract project and task information
	mapping, err := m.ParseTaskName(task.Name)
	if err != nil {
		return fmt.Errorf("failed to parse task name '%s': %w", task.Name, err)
	}

	// Get or create target client
	targetClient, err := m.getOrCreateClient(mapping.ClientName)
	if err != nil {
		return fmt.Errorf("failed to get/create client '%s': %w", mapping.ClientName, err)
	}

	// Get or create target project
	targetProject, err := m.getOrCreateProject(mapping.ProjectName, targetClient.ID)
	if err != nil {
		return fmt.Errorf("failed to get/create project '%s': %w", mapping.ProjectName, err)
	}

	// Get or create target task
	targetTask, err := m.getOrCreateTask(targetProject.ID, mapping.NewTaskName)
	if err != nil {
		return fmt.Errorf("failed to get/create task '%s': %w", mapping.NewTaskName, err)
	}

	// Create the time entry in target workspace
	if err := m.createTargetTimeEntry(entry, targetProject.ID, targetTask.ID); err != nil {
		return fmt.Errorf("failed to create target time entry: %w", err)
	}

	return nil
}

// ParseTaskName parses the old task format and returns mapping information
func (m *MigrationService) ParseTaskName(taskName string) (*ProjectTaskMapping, error) {
	// Expected format: "<real-world project name>/TASK<task number>"
	// Extract using regex
	re := regexp.MustCompile(`^(.+)/TASK(\d+)$`)
	matches := re.FindStringSubmatch(taskName)

	if len(matches) != 3 {
		return nil, fmt.Errorf("task name '%s' does not match expected format '<project>/TASK<number>'", taskName)
	}

	projectName := strings.TrimSpace(matches[1])
	taskNumber := matches[2]
	newTaskName := fmt.Sprintf("TASK %s", taskNumber) // Note the space

	// Determine client name
	clientName := m.config.DefaultClientName
	if m.config.ClientMapping != nil {
		if mappedClient, exists := m.config.ClientMapping[projectName]; exists {
			clientName = mappedClient
		}
	} else if m.config.CreateClients {
		// Use project name as client name if auto-creating clients
		clientName = projectName + " Client"
	}

	return &ProjectTaskMapping{
		OriginalTaskName: taskName,
		ProjectName:      projectName,
		TaskNumber:       taskNumber,
		NewTaskName:      newTaskName,
		ClientName:       clientName,
	}, nil
}

// getSourceTask retrieves a task from the source workspace
func (m *MigrationService) getSourceTask(taskID string) (*Task, error) {
	if taskID == "" {
		return nil, fmt.Errorf("empty task ID")
	}

	for tasks, err := range m.client.IterProjectTasks(m.sourceWorkspace.ID, m.sourceProject.ID) {
		if err != nil {
			return nil, err
		}

		for _, task := range tasks {
			if task.ID == taskID {
				return &task, nil
			}
		}
	}

	return nil, fmt.Errorf("task with ID %s not found", taskID)
}

// getOrCreateClient gets existing or creates new client
func (m *MigrationService) getOrCreateClient(clientName string) (*Client, error) {
	// Check cache first
	if client, exists := m.targetClients[clientName]; exists {
		return client, nil
	}

	// Create new client if enabled
	if m.config.CreateClients && !m.config.DryRun {
		client, err := m.client.CreateClient(m.targetWorkspace.ID, clientName)
		if err != nil {
			return nil, err
		}

		m.targetClients[clientName] = client
		m.stats.ClientsCreated++
		log.Printf("Created client: %s", clientName)
		return client, nil
	}

	if m.config.DryRun {
		log.Printf("DRY RUN: Would create client: %s", clientName)
		// Return a dummy client for dry run
		dummyClient := &Client{ID: "dummy", Name: clientName}
		return dummyClient, nil
	}

	return nil, fmt.Errorf("client '%s' not found and auto-creation disabled", clientName)
}

// getOrCreateProject gets existing or creates new project
func (m *MigrationService) getOrCreateProject(projectName, clientID string) (*Project, error) {
	// Check cache first
	if project, exists := m.targetProjects[projectName]; exists {
		return project, nil
	}

	// Try to find existing project
	for projects, err := range m.client.IterProjects(m.targetWorkspace.ID) {
		if err != nil {
			return nil, err
		}

		for _, proj := range projects {
			if proj.Name == projectName {
				m.targetProjects[projectName] = &proj
				return &proj, nil
			}
		}
	}

	// Create new project
	if m.config.DryRun {
		log.Printf("DRY RUN: Would create project: %s", projectName)
		dummyProject := &Project{ID: "dummy", Name: projectName, ClientID: clientID}
		m.targetProjects[projectName] = dummyProject
		return dummyProject, nil
	}

	project, err := m.client.CreateProject(m.targetWorkspace.ID, projectName)
	if err != nil {
		return nil, err
	}

	m.targetProjects[projectName] = project
	m.stats.ProjectsCreated++
	log.Printf("Created project: %s", projectName)
	return project, nil
}

// getOrCreateTask gets existing or creates new task
func (m *MigrationService) getOrCreateTask(projectID, taskName string) (*Task, error) {
	cacheKey := fmt.Sprintf("%s/%s", projectID, taskName)

	// Check cache first
	if task, exists := m.targetTasks[cacheKey]; exists {
		return task, nil
	}

	// Try to find existing task
	for tasks, err := range m.client.IterProjectTasks(m.targetWorkspace.ID, projectID) {
		if err != nil {
			return nil, err
		}

		for _, task := range tasks {
			if task.Name == taskName {
				m.targetTasks[cacheKey] = &task
				return &task, nil
			}
		}
	}

	// Create new task
	if m.config.DryRun {
		log.Printf("DRY RUN: Would create task: %s", taskName)
		dummyTask := &Task{ID: "dummy", Name: taskName, ProjectID: projectID}
		m.targetTasks[cacheKey] = dummyTask
		return dummyTask, nil
	}

	task, err := m.client.CreateTask(m.targetWorkspace.ID, projectID, taskName)
	if err != nil {
		return nil, err
	}

	m.targetTasks[cacheKey] = task
	m.stats.TasksCreated++
	log.Printf("Created task: %s in project %s", taskName, projectID)
	return task, nil
}

// createTargetTimeEntry creates a time entry in the target workspace
func (m *MigrationService) createTargetTimeEntry(sourceEntry *TimeEntry, targetProjectID, targetTaskID string) error {
	if m.config.DryRun {
		log.Printf("DRY RUN: Would create time entry: %s (%v to %v)",
			sourceEntry.Description,
			sourceEntry.TimeInterval.Start,
			sourceEntry.TimeInterval.End)
		return nil
	}

	// Create the new time entry request
	request := NewTimeEntryRequest{
		Start:       sourceEntry.TimeInterval.Start,
		End:         sourceEntry.TimeInterval.End,
		Billable:    sourceEntry.Billable,
		Description: sourceEntry.Description,
		ProjectID:   targetProjectID,
		TaskID:      targetTaskID,
		TagIDs:      sourceEntry.TagIDs, // Keep original tags
	}

	_, err := m.client.CreateTimeEntryForUser(m.targetWorkspace.ID, m.currentUser.ID, request)
	if err != nil {
		return err
	}

	m.stats.TimeEntriesCreated++
	return nil
}

// logMigrationSummary logs the final migration statistics
func (m *MigrationService) logMigrationSummary() {
	duration := m.stats.EndTime.Sub(m.stats.StartTime)

	log.Printf("=== MIGRATION COMPLETED ===")
	log.Printf("Duration: %v", duration)
	log.Printf("Time Entries Processed: %d", m.stats.TimeEntriesProcessed)
	log.Printf("Time Entries Created: %d", m.stats.TimeEntriesCreated)
	log.Printf("Projects Created: %d", m.stats.ProjectsCreated)
	log.Printf("Tasks Created: %d", m.stats.TasksCreated)
	log.Printf("Clients Created: %d", m.stats.ClientsCreated)
	log.Printf("Errors: %d", len(m.stats.Errors))

	if len(m.stats.Errors) > 0 {
		log.Printf("Error details:")
		for _, err := range m.stats.Errors {
			log.Printf("  - %s", err)
		}
	}
}

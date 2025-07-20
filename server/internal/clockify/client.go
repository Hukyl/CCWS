package clockify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type APIClient struct {
	apiKey   string
	client   *http.Client
	pageSize int
}

const baseURL = "https://api.clockify.me/api/v1"

func NewDefaultClient(apiKey string) *APIClient {
	return &APIClient{
		apiKey:   apiKey,
		client:   &http.Client{},
		pageSize: 5000, // max possible page size
	}
}

// * HTTP methods utilities

func (c *APIClient) get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.apiKey)

	return c.client.Do(req)
}

func (c *APIClient) post(url string, data any) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return c.client.Do(req)
}

func (c *APIClient) put(url string, data any) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return c.client.Do(req)
}

func (c *APIClient) delete(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.apiKey)

	return c.client.Do(req)
}

func (c *APIClient) patch(url string, data any) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return c.client.Do(req)
}

// * Actual API methods

// GetWorkspaces retrieves all workspaces for the authenticated user
func (c *APIClient) GetWorkspaces() ([]Workspace, error) {
	url := fmt.Sprintf("%s/workspaces", baseURL)

	resp, err := c.get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var workspaces []Workspace
	if err := json.NewDecoder(resp.Body).Decode(&workspaces); err != nil {
		return nil, err
	}

	return workspaces, nil
}

// GetCurrentUser retrieves the currently authenticated user
func (c *APIClient) GetCurrentUser() (*User, error) {
	url := fmt.Sprintf("%s/user", baseURL)

	resp, err := c.get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetWorkspaceUsers retrieves a page of users in a workspace
func (c *APIClient) GetWorkspaceUsers(workspaceID string, page int) ([]User, error) {
	url := fmt.Sprintf("%s/workspaces/%s/users", baseURL, workspaceID)

	resp, err := c.get(url + "?page=" + strconv.Itoa(page) + "&page-size=" + strconv.Itoa(c.pageSize))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

// GetProjects retrieves a page of projects in a workspace
func (c *APIClient) GetProjects(workspaceID string, page int) ([]Project, error) {
	url := fmt.Sprintf("%s/workspaces/%s/projects", baseURL, workspaceID)

	resp, err := c.get(url + "?page=" + strconv.Itoa(page) + "&page-size=" + strconv.Itoa(c.pageSize))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

// CreateProject creates a new project in a workspace
func (c *APIClient) CreateProject(workspaceID, name string) (*Project, error) {
	url := fmt.Sprintf("%s/workspaces/%s/projects", baseURL, workspaceID)

	project := map[string]any{
		"name":     name,
		"billable": true,
		"public":   false,
	}

	resp, err := c.post(url, project)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var createdProject Project
	if err := json.NewDecoder(resp.Body).Decode(&createdProject); err != nil {
		return nil, err
	}

	return &createdProject, nil
}

// GetClients retrieves a page of clients in a workspace
func (c *APIClient) GetClients(workspaceID string, page int) ([]Client, error) {
	url := fmt.Sprintf("%s/workspaces/%s/clients", baseURL, workspaceID)

	resp, err := c.get(url + "?page=" + strconv.Itoa(page) + "&page-size=" + strconv.Itoa(c.pageSize))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var clients []Client
	if err := json.NewDecoder(resp.Body).Decode(&clients); err != nil {
		return nil, err
	}

	return clients, nil
}

// CreateClient creates a new client in a workspace
func (c *APIClient) CreateClient(workspaceID, name string) (*Client, error) {
	url := fmt.Sprintf("%s/workspaces/%s/clients", baseURL, workspaceID)

	client := map[string]any{
		"name": name,
	}

	resp, err := c.post(url, client)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var createdClient Client
	if err := json.NewDecoder(resp.Body).Decode(&createdClient); err != nil {
		return nil, err
	}

	return &createdClient, nil
}

// GetTags retrieves a page of tags in a workspace
func (c *APIClient) GetTags(workspaceID string, page int) ([]Tag, error) {
	url := fmt.Sprintf("%s/workspaces/%s/tags", baseURL, workspaceID)

	resp, err := c.get(url + "?page=" + strconv.Itoa(page) + "&page-size=" + strconv.Itoa(c.pageSize))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var tags []Tag
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return nil, err
	}

	return tags, nil
}

// CreateTag creates a new tag in a workspace
func (c *APIClient) CreateTag(workspaceID, name string) (*Tag, error) {
	url := fmt.Sprintf("%s/workspaces/%s/tags", baseURL, workspaceID)

	tag := map[string]any{
		"name": name,
	}

	resp, err := c.post(url, tag)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var createdTag Tag
	if err := json.NewDecoder(resp.Body).Decode(&createdTag); err != nil {
		return nil, err
	}

	return &createdTag, nil
}

// GetTimeEntries retrieves a page of time entries for a user in a workspace with optional filters
func (c *APIClient) GetTimeEntries(workspaceID, userID string, start, end *time.Time, page int) ([]TimeEntry, error) {
	urlStr := fmt.Sprintf("%s/workspaces/%s/user/%s/time-entries", baseURL, workspaceID, userID)

	// Add query parameters for filtering
	params := url.Values{}
	if start != nil {
		params.Add("start", start.Format(time.RFC3339))
	}
	if end != nil {
		params.Add("end", end.Format(time.RFC3339))
	}

	if len(params) > 0 {
		urlStr += "?" + params.Encode()
	}

	resp, err := c.get(urlStr + "?page=" + strconv.Itoa(page) + "&page-size=" + strconv.Itoa(c.pageSize))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var timeEntries []TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntries); err != nil {
		return nil, err
	}

	return timeEntries, nil
}

// GetTimeEntry retrieves a specific time entry by ID
func (c *APIClient) GetTimeEntry(workspaceID, timeEntryID string) (*TimeEntry, error) {
	url := fmt.Sprintf("%s/workspaces/%s/time-entries/%s", baseURL, workspaceID, timeEntryID)

	resp, err := c.get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var timeEntry TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntry); err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

// CreateTimeEntry creates a new time entry in a workspace
func (c *APIClient) CreateTimeEntry(workspaceID string, request NewTimeEntryRequest) (*TimeEntry, error) {
	url := fmt.Sprintf("%s/workspaces/%s/time-entries", baseURL, workspaceID)

	resp, err := c.post(url, request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var timeEntry TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntry); err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

// CreateTimeEntryForUser creates a new time entry for a specific user in a workspace
func (c *APIClient) CreateTimeEntryForUser(workspaceID, userID string, request NewTimeEntryRequest) (*TimeEntry, error) {
	url := fmt.Sprintf("%s/workspaces/%s/user/%s/time-entries", baseURL, workspaceID, userID)

	resp, err := c.post(url, request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var timeEntry TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntry); err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

// UpdateTimeEntry updates an existing time entry
func (c *APIClient) UpdateTimeEntry(workspaceID, timeEntryID string, request UpdateTimeEntryRequest) (*TimeEntry, error) {
	url := fmt.Sprintf("%s/workspaces/%s/time-entries/%s", baseURL, workspaceID, timeEntryID)

	resp, err := c.put(url, request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var timeEntry TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntry); err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

// StopTimeEntry stops a currently running time entry for a user
func (c *APIClient) StopTimeEntry(workspaceID, userID string, endTime time.Time) (*TimeEntry, error) {
	url := fmt.Sprintf("%s/workspaces/%s/user/%s/time-entries", baseURL, workspaceID, userID)

	request := map[string]any{
		"end": endTime.Format(time.RFC3339),
	}

	resp, err := c.patch(url, request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var timeEntry TimeEntry
	if err := json.NewDecoder(resp.Body).Decode(&timeEntry); err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

// DeleteTimeEntry deletes a time entry
func (c *APIClient) DeleteTimeEntry(workspaceID, timeEntryID string) error {
	url := fmt.Sprintf("%s/workspaces/%s/time-entries/%s", baseURL, workspaceID, timeEntryID)

	resp, err := c.delete(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete time entry, status: %d", resp.StatusCode)
	}

	return nil
}

// GetProjectTasks retrieves a page of tasks for a project
func (c *APIClient) GetProjectTasks(workspaceID, projectID string, page int) ([]Task, error) {
	url := fmt.Sprintf("%s/workspaces/%s/projects/%s/tasks", baseURL, workspaceID, projectID)

	resp, err := c.get(url + "?page=" + strconv.Itoa(page) + "&page-size=" + strconv.Itoa(c.pageSize))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var tasks []Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// IterProjectTasks iterates over all tasks for a project, page by page
func (c *APIClient) IterProjectTasks(workspaceID, projectID string) iter.Seq2[[]Task, error] {
	return func(yield func([]Task, error) bool) {
		page := 1
		for {
			tasks, err := c.GetProjectTasks(workspaceID, projectID, page)
			if err != nil {
				yield(nil, err)
				return
			}

			if len(tasks) == 0 {
				return
			}

			if !yield(tasks, nil) {
				return
			}

			page++
		}
	}
}

// CreateTask creates a new task in a project
func (c *APIClient) CreateTask(workspaceID, projectID, name string) (*Task, error) {
	url := fmt.Sprintf("%s/workspaces/%s/projects/%s/tasks", baseURL, workspaceID, projectID)

	task := map[string]any{
		"name":   name,
		"status": "ACTIVE",
	}

	resp, err := c.post(url, task)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var createdTask Task
	if err := json.NewDecoder(resp.Body).Decode(&createdTask); err != nil {
		return nil, err
	}

	return &createdTask, nil
}

// * Helper methods to simplify common operations

// IterWorkspaceUsers iterates over all users in a workspace, page by page
func (c *APIClient) IterWorkspaceUsers(workspaceID string) iter.Seq2[[]User, error] {
	return func(yield func([]User, error) bool) {
		page := 1
		for {
			users, err := c.GetWorkspaceUsers(workspaceID, page)
			if err != nil {
				yield(nil, err)
				return
			}

			if len(users) == 0 {
				return
			}

			if !yield(users, nil) {
				return
			}

			page++
		}
	}
}

// IterTimeEntries iterates over all time entries for a user in a workspace, page by page
func (c *APIClient) IterTimeEntries(workspaceID, userID string, start, end *time.Time) iter.Seq2[[]TimeEntry, error] {
	return func(yield func([]TimeEntry, error) bool) {
		page := 1
		for {
			timeEntries, err := c.GetTimeEntries(workspaceID, userID, start, end, page)
			if err != nil {
				yield(nil, err)
				return
			}

			if len(timeEntries) == 0 {
				return
			}

			if !yield(timeEntries, nil) {
				return
			}

			page++
		}
	}
}

// IterTags iterates over all tags in a workspace, page by page
func (c *APIClient) IterTags(workspaceID string) iter.Seq2[[]Tag, error] {
	return func(yield func([]Tag, error) bool) {
		page := 1
		for {
			tags, err := c.GetTags(workspaceID, page)
			if err != nil {
				yield(nil, err)
				return
			}

			if len(tags) == 0 {
				return
			}

			if !yield(tags, nil) {
				return
			}

			page++
		}
	}
}

// IterClients iterates over all clients in a workspace, page by page
func (c *APIClient) IterClients(workspaceID string) iter.Seq2[[]Client, error] {
	return func(yield func([]Client, error) bool) {
		page := 1
		for {
			clients, err := c.GetClients(workspaceID, page)
			if err != nil {
				yield(nil, err)
				return
			}

			if len(clients) == 0 {
				return
			}

			if !yield(clients, nil) {
				return
			}

			page++
		}
	}
}

// IterProjects iterates over all projects in a workspace, page by page
func (c *APIClient) IterProjects(workspaceID string) iter.Seq2[[]Project, error] {
	return func(yield func([]Project, error) bool) {
		page := 1
		for {
			projects, err := c.GetProjects(workspaceID, page)
			if err != nil {
				yield(nil, err)
				return
			}

			if len(projects) == 0 {
				return
			}

			if !yield(projects, nil) {
				return
			}

			page++
		}
	}
}

// StartTimer starts a new timer for a user (creates a time entry without end time)
func (c *APIClient) StartTimer(workspaceID, userID, description string, projectID *string, taskID *string, tagIDs []string) (*TimeEntry, error) {
	request := NewTimeEntryRequest{
		Start:       time.Now(),
		Billable:    true,
		Description: description,
		TagIDs:      tagIDs,
	}

	if projectID != nil {
		request.ProjectID = *projectID
	}

	if taskID != nil {
		request.TaskID = *taskID
	}

	if tagIDs == nil {
		request.TagIDs = make([]string, 0)
	}

	return c.CreateTimeEntryForUser(workspaceID, userID, request)
}

// CreatePastTimeEntry creates a completed time entry for a specific date and duration
func (c *APIClient) CreatePastTimeEntry(workspaceID, userID string, startTime time.Time, duration time.Duration, description string, projectID *string, taskID *string, tagIDs []string, billable bool) (*TimeEntry, error) {
	endTime := startTime.Add(duration)

	request := NewTimeEntryRequest{
		Start:       startTime,
		End:         &endTime,
		Billable:    billable,
		Description: description,
		TagIDs:      tagIDs,
	}

	if projectID != nil {
		request.ProjectID = *projectID
	}

	if taskID != nil {
		request.TaskID = *taskID
	}

	if tagIDs == nil {
		request.TagIDs = make([]string, 0)
	}

	return c.CreateTimeEntryForUser(workspaceID, userID, request)
}

// CreateTimeEntryWithDates creates a time entry with specific start and end times
func (c *APIClient) CreateTimeEntryWithDates(workspaceID, userID string, startTime, endTime time.Time, description string, projectID *string, taskID *string, tagIDs []string, billable bool) (*TimeEntry, error) {
	request := NewTimeEntryRequest{
		Start:       startTime,
		End:         &endTime,
		Billable:    billable,
		Description: description,
		TagIDs:      tagIDs,
	}

	if projectID != nil {
		request.ProjectID = *projectID
	}

	if taskID != nil {
		request.TaskID = *taskID
	}

	if tagIDs == nil {
		request.TagIDs = make([]string, 0)
	}

	return c.CreateTimeEntryForUser(workspaceID, userID, request)
}

// CreateHistoricalWorkday creates multiple time entries for a past workday
func (c *APIClient) CreateHistoricalWorkday(workspaceID, userID string, date time.Time, entries []HistoricalEntry) ([]*TimeEntry, error) {
	var results []*TimeEntry
	var errors []error

	for _, entry := range entries {
		startTime := time.Date(date.Year(), date.Month(), date.Day(),
			entry.StartHour, entry.StartMinute, 0, 0, date.Location())

		timeEntry, err := c.CreatePastTimeEntry(
			workspaceID, userID, startTime, entry.Duration,
			entry.Description, entry.ProjectID, entry.TaskID, entry.TagIDs, entry.Billable,
		)

		if err != nil {
			errors = append(errors, fmt.Errorf("failed to create entry '%s': %w", entry.Description, err))
			continue
		}

		results = append(results, timeEntry)
	}

	if len(errors) > 0 {
		return results, fmt.Errorf("some entries failed: %v", errors)
	}

	return results, nil
}

// LogPastWorkSession creates a time entry for past work with common defaults
func (c *APIClient) LogPastWorkSession(workspaceID, userID string, date time.Time, startHour, startMinute int, durationHours float64, description string, projectID string) (*TimeEntry, error) {
	startTime := time.Date(date.Year(), date.Month(), date.Day(), startHour, startMinute, 0, 0, date.Location())
	duration := time.Duration(durationHours * float64(time.Hour))

	return c.CreatePastTimeEntry(workspaceID, userID, startTime, duration, description, &projectID, nil, nil, true)
}

// FindWorkspaceByName finds a workspace by name. Returns nil if not found.
func (c *APIClient) FindWorkspaceByName(name string) (*Workspace, error) {
	workspaces, err := c.GetWorkspaces()
	if err != nil {
		return nil, err
	}

	for _, ws := range workspaces {
		if ws.Name == name {
			return &ws, nil
		}
	}

	return nil, fmt.Errorf("workspace '%s' not found", name)
}

// FindProjectByName finds a project by name in a workspace. Returns nil if not found.
func (c *APIClient) FindProjectByName(workspaceID, name string) (*Project, error) {
	for projects, err := range c.IterProjects(workspaceID) {
		if err != nil {
			return nil, err
		}

		for _, proj := range projects {
			if proj.Name == name {
				return &proj, nil
			}
		}
	}

	return nil, fmt.Errorf("project '%s' not found in workspace", name)
}

// GetProjectTimeEntries retrieves all time entries from a project
func (c *APIClient) GetProjectTimeEntries(workspaceID, projectID string, userID string) ([]TimeEntry, error) {
	// TODO: make a generator (iter.Seq2)
	var filteredEntries []TimeEntry

	for timeEntries, err := range c.IterTimeEntries(workspaceID, userID, nil, nil) {
		if err != nil {
			return nil, err
		}

		for _, entry := range timeEntries {
			if entry.ProjectID == projectID {
				filteredEntries = append(filteredEntries, entry)
			}
		}
	}

	return filteredEntries, nil
}

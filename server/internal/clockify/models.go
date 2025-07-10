package clockify

import (
	"fmt"
	"time"
)

// Workspace represents a Clockify workspace
type Workspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (w Workspace) String() string {
	return fmt.Sprintf("Workspace <%s>: %s", w.ID, w.Name)
}

// Client represents a client/customer in Clockify
type Client struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WorkspaceID string `json:"workspaceId"`
	Archived    bool   `json:"archived"`
	Note        string `json:"note,omitempty"`
}

func (c Client) String() string {
	return c.Name
}

func NewClient(id, name, workspaceId string) Client {
	return Client{
		ID:          id,
		Name:        name,
		WorkspaceID: workspaceId,
		Archived:    false,
	}
}

// User represents a user in Clockify
type User struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	Name             string `json:"name"`
	ProfilePicture   string `json:"profilePicture,omitempty"`
	ActiveWorkspace  string `json:"activeWorkspace,omitempty"`
	DefaultWorkspace string `json:"defaultWorkspace,omitempty"`
	Status           string `json:"status,omitempty"`
}

func (u User) String() string {
	if u.Name != "" {
		return u.Name
	}
	return u.Email
}

func NewUser(id, email, name string) User {
	return User{
		ID:    id,
		Email: email,
		Name:  name,
	}
}

// Tag represents a tag in Clockify for categorizing time entries
type Tag struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	WorkspaceID string `json:"workspaceId"`
	Archived    bool   `json:"archived"`
}

func (t Tag) String() string {
	return t.Name
}

func NewTag(id, name, workspaceId string) Tag {
	return Tag{
		ID:          id,
		Name:        name,
		WorkspaceID: workspaceId,
		Archived:    false,
	}
}

// Project represents a project in Clockify
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ClientID    string `json:"clientId,omitempty"`
	ClientName  string `json:"clientName,omitempty"`
	WorkspaceID string `json:"workspaceId"`
	Billable    bool   `json:"billable"`
	Public      bool   `json:"public"`
	Archived    bool   `json:"archived"`
	Color       string `json:"color,omitempty"`
	Note        string `json:"note,omitempty"`
	// Simplified for free plan - avoiding complex memberships and estimates
}

func (p Project) String() string {
	return p.Name
}

func NewProject(id, name, workspaceId string) Project {
	return Project{
		ID:          id,
		Name:        name,
		WorkspaceID: workspaceId,
		Billable:    true,
		Public:      false,
		Archived:    false,
	}
}

// Task represents a task within a project
type Task struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ProjectID string `json:"projectId"`
	Status    string `json:"status"`
	Estimate  string `json:"estimate,omitempty"`
}

func (t Task) String() string {
	return t.Name
}

func NewTask(id, name, projectId string) Task {
	return Task{
		ID:        id,
		Name:      name,
		ProjectID: projectId,
		Status:    "ACTIVE",
	}
}

// TimeInterval represents the time period for a time entry
type TimeInterval struct {
	Start    time.Time  `json:"start"`
	End      *time.Time `json:"end,omitempty"`
	Duration string     `json:"duration,omitempty"`
}

// TimeEntry represents a time log entry in Clockify
type TimeEntry struct {
	ID           string        `json:"id"`
	Description  string        `json:"description,omitempty"`
	TagIDs       []string      `json:"tagIds,omitempty"`
	UserID       string        `json:"userId"`
	Billable     bool          `json:"billable"`
	TaskID       string        `json:"taskId,omitempty"`
	ProjectID    string        `json:"projectId,omitempty"`
	TimeInterval *TimeInterval `json:"timeInterval"`
	WorkspaceID  string        `json:"workspaceId"`
	IsLocked     bool          `json:"isLocked,omitempty"`
}

func (te TimeEntry) String() string {
	if te.Description != "" {
		return te.Description
	}
	return fmt.Sprintf("TimeEntry %s", te.ID)
}

func NewTimeEntry(userID, workspaceID string, start time.Time) TimeEntry {
	return TimeEntry{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Billable:    true,
		TimeInterval: &TimeInterval{
			Start: start,
		},
		TagIDs: make([]string, 0),
	}
}

// NewTimeEntryRequest represents the structure for creating a new time entry
type NewTimeEntryRequest struct {
	Start       time.Time  `json:"start"`
	End         *time.Time `json:"end,omitempty"`
	Billable    bool       `json:"billable"`
	Description string     `json:"description,omitempty"`
	ProjectID   string     `json:"projectId,omitempty"`
	TaskID      string     `json:"taskId,omitempty"`
	TagIDs      []string   `json:"tagIds,omitempty"`
}

// UpdateTimeEntryRequest represents the structure for updating a time entry
type UpdateTimeEntryRequest struct {
	Start       time.Time  `json:"start"`
	End         *time.Time `json:"end,omitempty"`
	Billable    bool       `json:"billable"`
	Description string     `json:"description,omitempty"`
	ProjectID   string     `json:"projectId,omitempty"`
	TaskID      string     `json:"taskId,omitempty"`
	TagIDs      []string   `json:"tagIds,omitempty"`
}

// HistoricalEntry represents a time entry for bulk historical creation
type HistoricalEntry struct {
	StartHour   int           `json:"startHour"`   // Hour (0-23)
	StartMinute int           `json:"startMinute"` // Minute (0-59)
	Duration    time.Duration `json:"duration"`    // How long the work took
	Description string        `json:"description"`
	ProjectID   *string       `json:"projectId,omitempty"`
	TaskID      *string       `json:"taskId,omitempty"`
	TagIDs      []string      `json:"tagIds,omitempty"`
	Billable    bool          `json:"billable"`
}

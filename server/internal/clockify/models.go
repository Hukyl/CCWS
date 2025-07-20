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

// Webhooks

// WebhookTriggerSourceType represents the type of the source of the webhook trigger
type WebhookTriggerSourceType string

// WebhookTriggerSourceType values
const (
	ProjectIDTrigger    WebhookTriggerSourceType = "PROJECT_ID"
	UserIDTrigger       WebhookTriggerSourceType = "USER_ID"
	TagIDTrigger        WebhookTriggerSourceType = "TAG_ID"
	TaskIDTrigger       WebhookTriggerSourceType = "TASK_ID"
	WorkspaceIDTrigger  WebhookTriggerSourceType = "WORKSPACE_ID"
	UserGroupIDTrigger  WebhookTriggerSourceType = "USER_GROUP_ID"
	InvoiceIDTrigger    WebhookTriggerSourceType = "INVOICE_ID"
	AssignmentIDTrigger WebhookTriggerSourceType = "ASSIGNMENT_ID"
	ExpenseIDTrigger    WebhookTriggerSourceType = "EXPENSE_ID"
)

// WebhookEvent represents the type of the event that triggered the webhook
type WebhookEvent string

// WebhookEvent values
const (
	NewProjectEvent                   WebhookEvent = "NEW_PROJECT"
	NewTaskEvent                      WebhookEvent = "NEW_TASK"
	NewClientEvent                    WebhookEvent = "NEW_CLIENT"
	NewTimerStartedEvent              WebhookEvent = "NEW_TIMER_STARTED"
	TimerStoppedEvent                 WebhookEvent = "TIMER_STOPPED"
	TimeEntryUpdatedEvent             WebhookEvent = "TIME_ENTRY_UPDATED"
	TimeEntryDeletedEvent             WebhookEvent = "TIME_ENTRY_DELETED"
	TimeEntrySplitEvent               WebhookEvent = "TIME_ENTRY_SPLIT"
	NewTimeEntryEvent                 WebhookEvent = "NEW_TIME_ENTRY"
	TimeEntryRestoredEvent            WebhookEvent = "TIME_ENTRY_RESTORED"
	NewTagEvent                       WebhookEvent = "NEW_TAG"
	UserDeletedFromWorkspaceEvent     WebhookEvent = "USER_DELETED_FROM_WORKSPACE"
	UserJoinedWorkspaceEvent          WebhookEvent = "USER_JOINED_WORKSPACE"
	UserDeactivatedOnWorkspaceEvent   WebhookEvent = "USER_DEACTIVATED_ON_WORKSPACE"
	UserActivatedOnWorkspaceEvent     WebhookEvent = "USER_ACTIVATED_ON_WORKSPACE"
	UserEmailChangedEvent             WebhookEvent = "USER_EMAIL_CHANGED"
	UserUpdatedEvent                  WebhookEvent = "USER_UPDATED"
	NewInvoiceEvent                   WebhookEvent = "NEW_INVOICE"
	InvoiceUpdatedEvent               WebhookEvent = "INVOICE_UPDATED"
	NewApprovalRequestEvent           WebhookEvent = "NEW_APPROVAL_REQUEST"
	ApprovalRequestStatusUpdatedEvent WebhookEvent = "APPROVAL_REQUEST_STATUS_UPDATED"
	TimeOffRequestRequestedEvent      WebhookEvent = "TIME_OFF_REQUESTED"
	TimeOffRequestApprovedEvent       WebhookEvent = "TIME_OFF_REQUEST_APPROVED"
	TimeOffRequestRejectedEvent       WebhookEvent = "TIME_OFF_REQUEST_REJECTED"
	TimeOffRequestWithdrawnEvent      WebhookEvent = "TIME_OFF_REQUEST_WITHDRAWN"
	BalanceUpdatedEvent               WebhookEvent = "BALANCE_UPDATED"
	TagUpdatedEvent                   WebhookEvent = "TAG_UPDATED"
	TagDeletedEvent                   WebhookEvent = "TAG_DELETED"
	TaskUpdatedEvent                  WebhookEvent = "TASK_UPDATED"
	ClientUpdatedEvent                WebhookEvent = "CLIENT_UPDATED"
	TaskDeletedEvent                  WebhookEvent = "TASK_DELETED"
	ClientDeletedEvent                WebhookEvent = "CLIENT_DELETED"
	ExpenseRestoredEvent              WebhookEvent = "EXPENSE_RESTORED"
	AssignmentCreatedEvent            WebhookEvent = "ASSIGNMENT_CREATED"
	AssignmentDeletedEvent            WebhookEvent = "ASSIGNMENT_DELETED"
	AssignmentPublishedEvent          WebhookEvent = "ASSIGNMENT_PUBLISHED"
	AssignmentUpdatedEvent            WebhookEvent = "ASSIGNMENT_UPDATED"
	ExpenseCreatedEvent               WebhookEvent = "EXPENSE_CREATED"
	ExpenseDeletedEvent               WebhookEvent = "EXPENSE_DELETED"
	ExpenseUpdatedEvent               WebhookEvent = "EXPENSE_UPDATED"
)

// WebhookRequest represents the structure for creating a new webhook
type WebhookRequest struct {
	Name              string                     `json:"name"`
	TriggerSource     []WebhookTriggerSourceType `json:"triggerSource"`
	TriggerSourceType WebhookTriggerSourceType   `json:"triggerSourceType"`
	TargetURL         string                     `json:"url"`
	Event             WebhookEvent               `json:"webhookEvent"`
}

// Webhook represents a webhook in Clockify
type Webhook struct {
	AuthToken         string                     `json:"authToken"`
	Enabled           bool                       `json:"enabled"`
	ID                string                     `json:"id"`
	Name              string                     `json:"name"`
	TriggerSource     []WebhookTriggerSourceType `json:"triggerSource"`
	TriggerSourceType WebhookTriggerSourceType   `json:"triggerSourceType"`
	TargetURL         string                     `json:"url"`
	UserID            string                     `json:"userId"`
	Event             WebhookEvent               `json:"webhookEvent"`
	WorkspaceID       string                     `json:"workspaceId"`
}

func (w Webhook) String() string {
	return fmt.Sprintf("Webhook <%s>: %s listening for %s at %s", w.ID, w.Name, w.Event, w.TargetURL)
}

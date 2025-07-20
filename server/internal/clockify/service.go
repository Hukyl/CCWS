package clockify

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

// WorkspaceWebhookService is a service for managing webhooks for a workspace.
//
// It is responsible for managing the lifecycle of a webhook. By default, it
// accepts all events regarding the workspace.
//
// *Note*: Clockify allows only one event per webhook. Therefore, to capture different events,
// the service creates multiple webhooks.
type WorkspaceWebhookService struct {
	apiClient *APIClient

	workspace Workspace
	url       string

	webhooks map[WebhookEvent]Webhook
}

func NewWorkspaceWebhookService(apiClient *APIClient, workspace Workspace, url string) *WorkspaceWebhookService {
	return &WorkspaceWebhookService{apiClient: apiClient, workspace: workspace, url: url}
}

var (
	ErrWebhookNotFound = errors.New("webhook not found")
	ErrDeleteWebhook   = errors.New("failed to delete webhook")
)

var eventToObject = map[WebhookEvent]any{
	NewTimerStartedEvent: &TimeEntry{},
	TimerStoppedEvent:    &TimeEntry{},
	NewClientEvent:       &Client{},
	NewProjectEvent:      &Project{},
	NewTagEvent:          &Tag{},
}

// Create creates a new webhook for the workspace.
func (s *WorkspaceWebhookService) Create() error {
	webhooks := make(map[WebhookEvent]Webhook)

	for event := range eventToObject {
		webhook, err := s.apiClient.CreateWebhook(s.workspace.ID, WebhookRequest{
			Name:              makeWebhookName(s.workspace.Name),
			Event:             event,
			TriggerSource:     []WebhookTriggerSourceType{WorkspaceIDTrigger},
			TriggerSourceType: WorkspaceIDTrigger,
			TargetURL:         s.url,
		})
		if err != nil {
			return fmt.Errorf("failed to create webhook: %w", err)
		}
		webhooks[event] = *webhook
	}

	s.webhooks = webhooks

	return nil
}

// Delete deletes the webhook for the workspace.
func (s *WorkspaceWebhookService) Delete() error {
	totalErr := ErrDeleteWebhook
	ok := true

	for _, webhook := range s.webhooks {
		err := s.apiClient.DeleteWebhook(s.workspace.ID, webhook.ID)
		if err != nil {
			totalErr = errors.Join(totalErr, err)
			ok = false
		}
	}

	if !ok {
		return totalErr
	}

	return nil
}

// TODO: webhook returns different schema than the API client uses. Create new models/adapt existing.
func (s *WorkspaceWebhookService) ProcessWebhook(r *http.Request) (WebhookEvent, any, error) {
	eventType := r.Header.Get("Clockify-Webhook-Event-Type")
	if eventType == "" {
		slog.Error("missing_event_type_header")
		return "", nil, errors.New("missing Clockify-Webhook-Event-Type header")
	}

	event := WebhookEvent(eventType)
	slog.Debug("processing_webhook", "event", event)

	objTemplate, ok := eventToObject[event]
	if !ok {
		slog.Error("unsupported_event_type", "event", event)
		return event, nil, fmt.Errorf("unsupported event type: %s", eventType)
	}

	// Signature verification (stub)
	signature := r.Header.Get("Clockify-Signature")
	if signature == "" {
		slog.Error("missing_signature_header")
		return event, nil, errors.New("missing Clockify-Signature header")
	}
	if !verifyClockifySignature(signature, r) {
		slog.Error("invalid_signature")
		return event, nil, errors.New("invalid signature")
	}

	// Read and decode body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("failed_to_read_body", "error", err)
		return event, nil, fmt.Errorf("failed to read body: %w", err)
	}
	defer r.Body.Close()

	obj := cloneObject(objTemplate)
	if err := json.Unmarshal(body, obj); err != nil {
		slog.Error("failed_to_unmarshal_body", "error", err, "obj", obj)
		return event, nil, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	return event, obj, nil
}

// cloneObject returns a new instance of the same type as the template (pointer to struct)
func cloneObject(template any) any {
	switch template.(type) {
	case *TimeEntry:
		return &TimeEntry{}
	case *Client:
		return &Client{}
	case *Project:
		return &Project{}
	case *Tag:
		return &Tag{}
	default:
		return nil
	}
}

// verifyClockifySignature is a stub for signature verification
func verifyClockifySignature(signature string, r *http.Request) bool {
	// TODO: Implement signature verification using webhook secret
	return true // Always valid for now
}

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hukyl/CCWS/internal/clockify"
	"github.com/Hukyl/CCWS/internal/config"
)

func makeWebhookHandler(webhookService *clockify.WorkspaceWebhookService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("webhook_received")

		// Log request headers
		for name, values := range r.Header {
			for _, value := range values {
				slog.Info("request_header", "name", name, "value", value)
			}
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("error_reading_request_body", "error", err)
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		// Output the full request body as text
		slog.Info("request_body", "body", string(body))

		// Return a success response
		w.WriteHeader(http.StatusOK)

		event, obj, err := webhookService.ProcessWebhook(r)
		if err != nil {
			slog.Error("error_processing_webhook", "error", err)
		}

		slog.Info("webhook_processed", "event", event, "obj", obj)
	}
}

var (
	webhookURL    string
	workspaceName string
)

func main() {
	flag.StringVar(&webhookURL, "webhook-url", "http://localhost:8080", "The URL to send the webhook to")
	flag.StringVar(&workspaceName, "workspace-name", "", "The name of the workspace to delete time entries from")
	flag.Parse()

	if workspaceName == "" {
		slog.Error("workspace_name_is_required")
		return
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed_to_load_config", "error", err)
		return
	}

	apiKey := cfg.ClockifyAPIKey
	client := clockify.NewDefaultClient(apiKey)

	workspace, err := client.FindWorkspaceByName(workspaceName)
	if err != nil {
		slog.Error("failed_to_find_workspace", "error", err)
		return
	}
	fmt.Println("Found workspace:", workspace)

	webhookService := clockify.NewWorkspaceWebhookService(
		client,
		*workspace,
		webhookURL,
	)

	err = webhookService.Create()
	if err != nil {
		slog.Error("failed_to_create_webhook", "error", err)
		return
	}
	defer func() {
		err = webhookService.Delete()
		if err != nil {
			slog.Error("failed_to_delete_webhook", "error", err)
			return
		}
		fmt.Println("Webhook deleted")
	}()

	fmt.Println("Webhook created")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	// Create a http server that will receive the webhook
	server := http.Server{
		Addr:    ":8080",
		Handler: makeWebhookHandler(webhookService),
	}

	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error("failed_to_start_server", "error", err)
			return
		}
	}()

	fmt.Println("Server started on http://localhost:8080")

	<-signals

	err = server.Shutdown(context.Background())
	if err != nil {
		slog.Error("failed_to_shutdown_server", "error", err)
		return
	}
	fmt.Println("Server shutdown gracefully")
}

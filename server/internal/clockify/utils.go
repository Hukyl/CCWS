package clockify

import (
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
)

// kebabify converts a string to kebab-case
func kebabify(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "-"))
}

func makeWebhookName(workspaceName string) string {
	const maxWebhookNameLength = 30
	const randomPartLength = 6
	const suffix = "-wh"
	const allowedRunes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// 1. Cut the workspace name up to 20 chars
	if len(workspaceName) > 20 {
		workspaceName = workspaceName[:20]
	}

	// 2. Strip whitespace and control chars
	stripped := make([]rune, 0, len(workspaceName))
	for _, r := range workspaceName {
		if r > 31 && r != 127 && r != ' ' && r != '\t' && r != '\n' && r != '\r' {
			stripped = append(stripped, r)
		}
	}

	// 3. Kebabify it
	kebabified := kebabify(string(stripped))

	// 4. Add a hyphen and 6 random symbols (A-Z, a-z, 0-9)
	randomPart := make([]rune, randomPartLength)
	for i := range randomPart {
		randomPart[i] = rune(allowedRunes[seededRandInt(len(allowedRunes))])
	}

	// 5. Add a hyphen and 'wh'
	name := fmt.Sprintf("%s-%s%s", kebabified, string(randomPart), suffix)

	// 6. Ensure total length <= 30
	if len(name) > maxWebhookNameLength {
		slog.Warn("webhook_name_too_long", "name", name, "max_length", maxWebhookNameLength)
		name = name[:maxWebhookNameLength]
	}

	return name
}

// seededRandInt returns a random int in [0, n) using math/rand with a seeded source.
func seededRandInt(n int) int {
	// Use a package-level seeded rand for thread safety in real code
	// Here, for simplicity, use rand.Intn
	return rand.Intn(n)
}

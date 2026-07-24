//go:build integration

package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestIntegrationDeliversMessageToMailpit(t *testing.T) {
	host := os.Getenv("SEMREL_TEST_SMTP_HOST")
	port := os.Getenv("SEMREL_TEST_SMTP_PORT")
	apiURL := os.Getenv("SEMREL_TEST_MAILPIT_API")
	if host == "" || port == "" || apiURL == "" {
		t.Fatal("SEMREL_TEST_SMTP_HOST, SEMREL_TEST_SMTP_PORT, and SEMREL_TEST_MAILPIT_API are required")
	}

	var stderr bytes.Buffer
	code := run(context.Background(), func(key string) string {
		return map[string]string{
			"SEMREL_PLUGIN_SMTP_HOST": host,
			"SEMREL_PLUGIN_SMTP_PORT": port,
			"SEMREL_PLUGIN_FROM":      "semrel@example.test",
			"SEMREL_PLUGIN_TO":        "release@example.test",
			"SEMREL_VERSION":          "v1.2.3",
			"SEMREL_CHANGELOG":        "- integration release",
		}[key]
	}, &stderr)
	if code != 0 {
		t.Fatalf("email hook code = %d, stderr = %q", code, stderr.String())
	}

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		response, err := http.Get(apiURL + "/api/v1/messages")
		if err == nil {
			body, readErr := io.ReadAll(response.Body)
			_ = response.Body.Close()
			if readErr == nil && strings.Contains(string(body), "[Release] v1.2.3") &&
				strings.Contains(string(body), "integration release") {
				return
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	t.Fatalf("Mailpit did not receive the release message; stderr = %q", stderr.String())
}

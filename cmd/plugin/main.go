package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	plugin "github.com/SemRels/hook-email/internal/plugin"
)

const pluginSchemaVersion = 1

type mailer interface {
	Notify(context.Context, string, string, string) error
}

var newMailer = func(cfg plugin.EmailConfig) mailer {
	return plugin.NewMailer(cfg)
}

func run(ctx context.Context, getenv func(string) string, stderr io.Writer) int {
	_, _ = fmt.Fprintf(stderr, "plugin_schema_version=%d\n", pluginSchemaVersion)
	host := getenv("SEMREL_PLUGIN_SMTP_HOST")
	from := getenv("SEMREL_PLUGIN_FROM")
	to := plugin.ParseRecipients(getenv("SEMREL_PLUGIN_TO"))
	version := firstNonEmpty(getenv("SEMREL_VERSION"), getenv("SEMREL_TAG_NAME"), getenv("SEMREL_NEXT_VERSION"))

	if host == "" || from == "" || len(to) == 0 {
		fmt.Fprintln(stderr, "hook-email: SEMREL_PLUGIN_SMTP_HOST, SEMREL_PLUGIN_FROM, and SEMREL_PLUGIN_TO are required")
		return 1
	}
	if version == "" {
		fmt.Fprintln(stderr, "hook-email: SEMREL_VERSION, SEMREL_TAG_NAME, or SEMREL_NEXT_VERSION is required")
		return 1
	}
	if getenv("SEMREL_DRY_RUN") == "true" {
		return 0
	}

	cfg := plugin.EmailConfig{
		SMTPHost: host,
		SMTPPort: getenv("SEMREL_PLUGIN_SMTP_PORT"),
		From:     from,
		To:       to,
		Username: getenv("SEMREL_PLUGIN_USERNAME"),
		Password: getenv("SEMREL_PLUGIN_PASSWORD"),
	}

	if err := newMailer(cfg).Notify(ctx, version, getenv("SEMREL_CHANGELOG"), getenv("SEMREL_TAG_NAME")); err != nil {
		fmt.Fprintln(stderr, "hook-email:", err)
		return 1
	}
	return 0
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	os.Exit(run(ctx, os.Getenv, os.Stderr))
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

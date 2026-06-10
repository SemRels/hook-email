package main

import (
	"bytes"
	"context"
	"errors"
	"testing"

	plugin "github.com/SemRels/hook-email/internal/plugin"
)

type fakeMailer struct {
	version   string
	changelog string
	tagName   string
	err       error
}

func (f *fakeMailer) Notify(_ context.Context, version, changelog, tagName string) error {
	f.version = version
	f.changelog = changelog
	f.tagName = tagName
	return f.err
}

func env(kv map[string]string) func(string) string {
	return func(key string) string { return kv[key] }
}

func TestRun_Success(t *testing.T) {

	fake := &fakeMailer{}
	old := newMailer
	newMailer = func(cfg plugin.EmailConfig) mailer {
		if cfg.SMTPHost != "smtp.example.com" {
			t.Fatalf("unexpected host: %s", cfg.SMTPHost)
		}
		return fake
	}
	defer func() { newMailer = old }()

	var stderr bytes.Buffer
	code := run(context.Background(), env(map[string]string{
		"SEMREL_PLUGIN_SMTP_HOST": "smtp.example.com",
		"SEMREL_PLUGIN_FROM":      "bot@example.com",
		"SEMREL_PLUGIN_TO":        "team@example.com",
		"SEMREL_VERSION":          "v1.2.3",
		"SEMREL_CHANGELOG":        "- feature",
		"SEMREL_TAG_NAME":         "v1.2.3",
	}), &stderr)

	if code != 0 || stderr.String() != "plugin_schema_version=1\n" {
		t.Fatalf("unexpected result: code=%d stderr=%q", code, stderr.String())
	}
	if fake.version != "v1.2.3" || fake.changelog != "- feature" || fake.tagName != "v1.2.3" {
		t.Fatalf("unexpected notify args: %+v", fake)
	}
}

func TestRun_DryRun(t *testing.T) {

	called := false
	old := newMailer
	newMailer = func(plugin.EmailConfig) mailer {
		called = true
		return &fakeMailer{}
	}
	defer func() { newMailer = old }()

	var stderr bytes.Buffer
	code := run(context.Background(), env(map[string]string{
		"SEMREL_PLUGIN_SMTP_HOST": "smtp.example.com",
		"SEMREL_PLUGIN_FROM":      "bot@example.com",
		"SEMREL_PLUGIN_TO":        "team@example.com",
		"SEMREL_VERSION":          "v1.2.3",
		"SEMREL_DRY_RUN":          "true",
	}), &stderr)

	if code != 0 || called {
		t.Fatalf("unexpected result: code=%d called=%v", code, called)
	}
}

func TestRun_ValidationError(t *testing.T) {

	var stderr bytes.Buffer
	code := run(context.Background(), env(map[string]string{}), &stderr)
	if code != 1 || stderr.Len() == 0 {
		t.Fatalf("unexpected result: code=%d stderr=%q", code, stderr.String())
	}
}

func TestRun_NotifyError(t *testing.T) {

	old := newMailer
	newMailer = func(plugin.EmailConfig) mailer {
		return &fakeMailer{err: errors.New("boom")}
	}
	defer func() { newMailer = old }()

	var stderr bytes.Buffer
	code := run(context.Background(), env(map[string]string{
		"SEMREL_PLUGIN_SMTP_HOST": "smtp.example.com",
		"SEMREL_PLUGIN_FROM":      "bot@example.com",
		"SEMREL_PLUGIN_TO":        "team@example.com",
		"SEMREL_VERSION":          "v1.2.3",
	}), &stderr)
	if code != 1 || stderr.Len() == 0 {
		t.Fatalf("unexpected result: code=%d stderr=%q", code, stderr.String())
	}
}

func TestRun_MissingVersion(t *testing.T) {
	var stderr bytes.Buffer
	code := run(context.Background(), env(map[string]string{
		"SEMREL_PLUGIN_SMTP_HOST": "smtp.example.com",
		"SEMREL_PLUGIN_FROM":      "bot@example.com",
		"SEMREL_PLUGIN_TO":        "team@example.com",
	}), &stderr)
	if code != 1 || stderr.Len() == 0 {
		t.Fatalf("unexpected result: code=%d stderr=%q", code, stderr.String())
	}
}

func TestFirstNonEmpty(t *testing.T) {
	if got := firstNonEmpty("", "v1.2.3", "v1.2.4"); got != "v1.2.3" {
		t.Fatalf("unexpected value: %s", got)
	}
}

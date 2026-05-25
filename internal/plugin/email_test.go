package plugin

import (
	"context"
	"errors"
	"net/smtp"
	"strings"
	"testing"
)

func TestParseRecipients(t *testing.T) {
	t.Parallel()

	recipients := ParseRecipients("a@example.com, b@example.com;c@example.com\nd@example.com")
	if got, want := len(recipients), 4; got != want {
		t.Fatalf("got %d recipients, want %d", got, want)
	}
}

func TestMailerValidate(t *testing.T) {
	t.Parallel()

	err := NewMailer(EmailConfig{}).Validate()
	if err == nil || !strings.Contains(err.Error(), "SMTP host") {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestMailerNotify_Success(t *testing.T) {
	t.Parallel()

	var (
		addr string
		auth smtp.Auth
		from string
		to   []string
		msg  string
	)

	m := NewMailerWithSender(EmailConfig{
		SMTPHost: "smtp.example.com",
		SMTPPort: "587",
		From:     "bot@example.com",
		To:       []string{"team@example.com"},
		Username: "user",
		Password: "secret",
	}, func(gotAddr string, gotAuth smtp.Auth, gotFrom string, gotTo []string, gotMsg []byte) error {
		addr = gotAddr
		auth = gotAuth
		from = gotFrom
		to = gotTo
		msg = string(gotMsg)
		return nil
	})

	err := m.Notify(context.Background(), "v1.2.3", "- feature", "v1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "smtp.example.com:587" {
		t.Fatalf("unexpected addr: %s", addr)
	}
	if auth == nil {
		t.Fatal("expected SMTP auth")
	}
	if from != "bot@example.com" || len(to) != 1 || to[0] != "team@example.com" {
		t.Fatalf("unexpected envelope: from=%s to=%v", from, to)
	}
	if !strings.Contains(msg, "Subject: [Release] v1.2.3") || !strings.Contains(msg, "Changelog:") {
		t.Fatalf("unexpected message: %s", msg)
	}
}

func TestMailerNotify_DefaultPort(t *testing.T) {
	t.Parallel()

	var addr string
	m := NewMailerWithSender(EmailConfig{
		SMTPHost: "smtp.example.com",
		From:     "bot@example.com",
		To:       []string{"team@example.com"},
	}, func(gotAddr string, _ smtp.Auth, _ string, _ []string, _ []byte) error {
		addr = gotAddr
		return nil
	})

	if err := m.Notify(context.Background(), "v1.2.3", "", "v1.2.3"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr != "smtp.example.com:25" {
		t.Fatalf("unexpected addr: %s", addr)
	}
}

func TestMailerNotify_ContextCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	m := NewMailer(EmailConfig{
		SMTPHost: "smtp.example.com",
		From:     "bot@example.com",
		To:       []string{"team@example.com"},
	})

	err := m.Notify(ctx, "v1.2.3", "", "")
	if err == nil || !strings.Contains(err.Error(), "context canceled") {
		t.Fatalf("expected context error, got %v", err)
	}
}

func TestMailerNotify_SendError(t *testing.T) {
	t.Parallel()

	m := NewMailerWithSender(EmailConfig{
		SMTPHost: "smtp.example.com",
		From:     "bot@example.com",
		To:       []string{"team@example.com"},
	}, func(string, smtp.Auth, string, []string, []byte) error {
		return errors.New("boom")
	})

	err := m.Notify(context.Background(), "v1.2.3", "", "")
	if err == nil || !strings.Contains(err.Error(), "send mail") {
		t.Fatalf("expected send error, got %v", err)
	}
}

func TestBuildBody_IncludesTagWhenDifferent(t *testing.T) {
	t.Parallel()

	body := buildBody("1.2.3", "- feature", "v1.2.3")
	if !strings.Contains(body, "Tag: v1.2.3") {
		t.Fatalf("expected tag in body, got %q", body)
	}
}

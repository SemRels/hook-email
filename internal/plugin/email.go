package plugin

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

type Sender func(addr string, auth smtp.Auth, from string, to []string, msg []byte) error

type EmailConfig struct {
	SMTPHost string
	SMTPPort string
	From     string
	To       []string
	Username string
	Password string
}

type Mailer struct {
	cfg      EmailConfig
	sendMail Sender
}

func NewMailer(cfg EmailConfig) *Mailer {
	if cfg.SMTPPort == "" {
		cfg.SMTPPort = "25"
	}
	return &Mailer{cfg: cfg, sendMail: smtp.SendMail}
}

func NewMailerWithSender(cfg EmailConfig, sender Sender) *Mailer {
	m := NewMailer(cfg)
	if sender != nil {
		m.sendMail = sender
	}
	return m
}

func ParseRecipients(raw string) []string {
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n'
	})

	recipients := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			recipients = append(recipients, part)
		}
	}
	return recipients
}

func (m *Mailer) Validate() error {
	var errs []string
	if strings.TrimSpace(m.cfg.SMTPHost) == "" {
		errs = append(errs, "SMTP host is required")
	}
	if strings.TrimSpace(m.cfg.From) == "" {
		errs = append(errs, "from address is required")
	}
	if len(m.cfg.To) == 0 {
		errs = append(errs, "at least one recipient is required")
	}
	if len(errs) > 0 {
		return fmt.Errorf("email: %s", strings.Join(errs, "; "))
	}
	return nil
}

func (m *Mailer) Notify(ctx context.Context, version, changelog, tagName string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("email: context canceled: %w", err)
	}
	if err := m.Validate(); err != nil {
		return err
	}

	release := firstNonEmpty(version, tagName)
	if release == "" {
		release = "unknown"
	}

	subject := fmt.Sprintf("[Release] %s", release)
	body := buildBody(release, changelog, tagName)
	message := buildMessage(m.cfg.From, m.cfg.To, subject, body)

	var auth smtp.Auth
	if m.cfg.Username != "" {
		auth = smtp.PlainAuth("", m.cfg.Username, m.cfg.Password, m.cfg.SMTPHost)
	}

	if err := m.sendMail(net.JoinHostPort(m.cfg.SMTPHost, m.cfg.SMTPPort), auth, m.cfg.From, m.cfg.To, []byte(message)); err != nil {
		return fmt.Errorf("email: send mail: %w", err)
	}
	return nil
}

func buildBody(version, changelog, tagName string) string {
	var builder strings.Builder
	builder.WriteString("A new SemRel release is available.\n\n")
	builder.WriteString("Version: ")
	builder.WriteString(version)
	builder.WriteString("\n")
	if tagName != "" && tagName != version {
		builder.WriteString("Tag: ")
		builder.WriteString(tagName)
		builder.WriteString("\n")
	}
	if trimmed := strings.TrimSpace(changelog); trimmed != "" {
		builder.WriteString("\nChangelog:\n")
		builder.WriteString(trimmed)
		builder.WriteString("\n")
	}
	return builder.String()
}

func buildMessage(from string, to []string, subject, body string) string {
	return strings.Join([]string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", strings.Join(to, ", ")),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

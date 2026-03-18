package smtp

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/ramdhanrizki/bytecode-api/configs"
	identityService "github.com/ramdhanrizki/bytecode-api/internal/identity/application/service"
	sharedLogger "github.com/ramdhanrizki/bytecode-api/internal/shared/logger"
)

const smtpDialTimeout = 10 * time.Second

type Sender struct {
	cfg    configs.SMTPConfig
	logger sharedLogger.Logger
}

func NewSender(cfg configs.SMTPConfig, logger sharedLogger.Logger) *Sender {
	return &Sender{cfg: cfg, logger: logger}
}

func (s *Sender) Send(ctx context.Context, message identityService.MailMessage) error {
	address := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))
	dialer := &net.Dialer{Timeout: smtpDialTimeout}
	connection, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("dial smtp server: %w", err)
	}
	defer func() { _ = connection.Close() }()

	client, err := smtp.NewClient(connection, s.cfg.Host)
	if err != nil {
		return fmt.Errorf("create smtp client: %w", err)
	}
	defer func() { _ = client.Close() }()

	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(&tls.Config{ServerName: s.cfg.Host, MinVersion: tls.VersionTLS12}); err != nil {
			return fmt.Errorf("start tls: %w", err)
		}
	}

	if s.cfg.Username != "" {
		auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("authenticate smtp client: %w", err)
		}
	}

	if err := client.Mail(s.cfg.From); err != nil {
		return fmt.Errorf("set smtp sender: %w", err)
	}
	if err := client.Rcpt(message.ToAddress); err != nil {
		return fmt.Errorf("set smtp recipient: %w", err)
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("create smtp data writer: %w", err)
	}

	if _, err := io.WriteString(writer, s.buildMessage(message)); err != nil {
		_ = writer.Close()
		return fmt.Errorf("write smtp message: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("close smtp message writer: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("quit smtp client: %w", err)
	}

	s.logger.Info("smtp email sent", zap.String("to", message.ToAddress), zap.String("subject", message.Subject))
	return nil
}

func (s *Sender) buildMessage(message identityService.MailMessage) string {
	from := (&mail.Address{Name: s.cfg.FromName, Address: s.cfg.From}).String()
	to := (&mail.Address{Name: message.ToName, Address: message.ToAddress}).String()
	headers := []string{
		"From: " + from,
		"To: " + to,
		"Subject: " + sanitizeHeader(message.Subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
	}

	return strings.Join(headers, "\r\n") + "\r\n\r\n" + message.Body
}

func sanitizeHeader(value string) string {
	value = strings.ReplaceAll(value, "\r", "")
	value = strings.ReplaceAll(value, "\n", "")
	return value
}

package provider

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"net/textproto"
	"text/template"

	"github.com/jordan-wright/email"

	"brickwall/internal/common"
)

type ISmtpProvider interface {
	Send(string, string, string, map[string]string) error
}

type SmtpProvider struct {
	ctx context.Context

	smtpServerHost     string
	smtpServerPort     int
	smtpServerUser     string
	smtpServerPassword string
	smtpSenderFrom     string
	smtpTemplates      string
	smtpUseTLS         bool
}

func NewSmtpProvider(ctx context.Context) ISmtpProvider {
	env := ctx.Value(common.KeyEnvProvider).(IEnvProvider)

	return &SmtpProvider{
		ctx: ctx,

		smtpServerHost:     env.GetString("SMTP_SERVER_HOST", DefSmtpServerHost),
		smtpServerPort:     env.GetInt("SMTP_SERVER_PORT", DefSmtpServerPort),
		smtpServerUser:     env.GetString("SMTP_SERVER_USER", DefSmtpServerUser),
		smtpServerPassword: env.GetString("SMTP_SERVER_PASSWORD", DefSmtpServerPassword),
		smtpSenderFrom:     env.GetString("SMTP_SENDER_FROM", DefSmtpSenderFrom),
		smtpTemplates:      env.GetString("SMTP_TEMPLATES", DefSmtpTemplates),
		smtpUseTLS:         env.GetBool("SMTP_USE_TLS", DefSmtpUseTLS),
	}
}

func (rcv *SmtpProvider) Send(to, subject, template string, data map[string]string) error {
	addr := fmt.Sprintf("%s:%d", rcv.smtpServerHost, rcv.smtpServerPort)
	auth := smtp.PlainAuth("", rcv.smtpServerUser, rcv.smtpServerPassword, rcv.smtpServerHost)
	file := fmt.Sprintf("%s/%s", rcv.smtpTemplates, template)

	body, err := loadTemplate(file, data)
	if err != nil {
		return err
	}
	email := &email.Email{
		To:      []string{to},
		From:    rcv.smtpSenderFrom,
		Subject: subject,
		HTML:    []byte(body),
		Headers: textproto.MIMEHeader{},
	}
	if rcv.smtpUseTLS {
		return email.SendWithTLS(addr, auth, &tls.Config{
			ServerName: rcv.smtpServerHost,
		})
	}
	return email.Send(addr, auth)
}

func loadTemplate(templatePath string, data map[string]string) (string, error) {
	var body bytes.Buffer

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("%w: %v", common.ErrEmailTemplateNotFound, err)
	}
	if err := tmpl.Execute(&body, data); err != nil {
		return "", fmt.Errorf("%w: %v", common.ErrEmailTemplateRendering, err)
	}
	return body.String(), nil
}

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
	"github.com/urfave/cli/v3"

	"brickwall/internal/common"
)

type ISmtpProvider interface {
	SendEmail(string, string, string, map[string]string) error
}

type SmtpProvider struct {
	ctx context.Context
}

func NewSmtpProvider(ctx context.Context) ISmtpProvider {
	return &SmtpProvider{ctx: ctx}
}

func (rcv *SmtpProvider) SendEmail(to, subject, template string, data map[string]string) error {
	cli := rcv.ctx.Value(common.KeyCommand).(*cli.Command)

	addr := fmt.Sprintf(
		"%s:%d",
		cli.String("smtp-server-host"),
		cli.Int("smtp-server-port"),
	)
	auth := smtp.PlainAuth(
		"",
		cli.String("smtp-server-user"),
		cli.String("smtp-server-password"),
		cli.String("smtp-server-host"),
	)
	file := fmt.Sprintf("%s/%s",
		cli.String("smtp-templates"), template,
	)
	body, err := loadTemplate(file, data)
	if err != nil {
		return err
	}
	email := &email.Email{
		To:      []string{to},
		From:    cli.String("smtp-sender-from"),
		Subject: subject,
		HTML:    []byte(body),
		Headers: textproto.MIMEHeader{},
	}
	if cli.Bool("smtp-use-tls") {
		return email.SendWithTLS(addr, auth, &tls.Config{
			ServerName: cli.String("smtp-server-host"),
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

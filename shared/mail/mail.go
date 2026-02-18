package mail

import (
	"context"

	"github.com/resend/resend-go/v3"
)

type MailClient struct {
	client *resend.Client
}

func NewMailClient(apiKey string) *MailClient {
	client := resend.NewClient(apiKey)
	return &MailClient{client: client}
}

func (m *MailClient) SendMail(ctx context.Context, from, to, subject, html string) error {
	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Html:    html,
		ReplyTo: from,
	}

	_, err := m.client.Emails.SendWithContext(ctx, params)
	return err
}

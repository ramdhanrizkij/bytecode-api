package service

import "context"

type MailMessage struct {
	ToAddress string
	ToName    string
	Subject   string
	Body      string
}

type MailSender interface {
	Send(ctx context.Context, message MailMessage) error
}

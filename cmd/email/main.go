package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/51st-state/api/pkg/email"
	"github.com/51st-state/api/pkg/event"

	pubsubNSQ "github.com/51st-state/api/pkg/pubsub/nsq"
	"github.com/nsqio/go-nsq"

	"github.com/playnet-public/flagenv"
	"go.uber.org/zap"

	"gopkg.in/gomail.v2"
)

var (
	nsqAddr      = flagenv.String("nsq-addr", "nsq", "the address of a nsq message queue")
	nsqTopic     = flagenv.String("nsq-topic", "events", "the topic for api events")
	smtpEmail    = flagenv.String("smtp-email", "local@example.com", "the email of the sender")
	smtpServer   = flagenv.String("smtp-server", "mail.example.com", "the smtp server to send emails from")
	smtpPort     = flagenv.Int("smtp-port", 587, "the smtp servers port to send emails from")
	smtpUsername = flagenv.String("smtp-username", "local@example.com", "the username for the smtp server")
	smtpPassword = flagenv.String("smtp-password", "1234", "the password for the smtp server")
)

const channelName = "email"

func main() {
	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Fatal(err.Error())
	}

	flagenv.Parse()

	pubsubC, err := pubsubNSQ.NewConsumer(*nsqTopic, channelName, *nsqAddr, nsq.NewConfig())
	if err != nil {
		logger.Fatal(err.Error())
	}

	finalConsumer := event.NewConsumer(pubsubC)

	d := gomail.NewPlainDialer(*smtpServer, *smtpPort, *smtpUsername, *smtpPassword)

	if err := finalConsumer.Consume(context.Background(), func(ctx context.Context, e *event.Event) error {
		l := logger.With(
			zap.String("event_uuid", e.Meta.UUID),
			zap.String("event_id", string(e.Meta.ID)),
			zap.String("event_version", e.Meta.Version),
			zap.String("event_created_at", e.Meta.CreatedAt.Format(time.RFC3339Nano)),
		)
		switch e.Meta.ID {
		case email.NewEvent:
			l.Info("processing new email send event")

			var decEmail email.Email

			if err := json.Unmarshal(e.Payload, &decEmail); err != nil {
				return err
			}

			l = l.With(
				zap.String("email_from", decEmail.Header.From),
				zap.Strings("email_to", decEmail.Header.To),
			)

			l.Info("sending email")
			if err := d.DialAndSend(newMessage(&decEmail)); err != nil {
				l.Error(
					"sending an email errored",
					zap.Error(err),
				)
				return err
			}
		}

		return nil
	}); err != nil {
		logger.Fatal(err.Error())
	}
}

func newMessage(e *email.Email) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", e.Header.From)
	m.SetHeader("To", e.Header.To...)
	m.SetHeader("Subject", e.Header.Subject)
	m.SetBody("text/html", e.Body)

	for _, v := range e.Header.Cc {
		m.SetAddressHeader("Cc", v.Email, v.Name)
	}

	for _, v := range e.Header.Bcc {
		m.SetAddressHeader("Bcc", v.Email, v.Name)
	}

	return m
}

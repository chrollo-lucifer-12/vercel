package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/chrollo-lucifer-12/shared/mail"
	"github.com/hibiken/asynq"
)

const TypeEmailSend = "email:send"

type EmailWorker struct {
	server     *asynq.Server
	mux        *asynq.ServeMux
	mailClient *mail.MailClient
}

func NewEmailWorkerServer(redisAddr string, apiKey string) *EmailWorker {

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "127.0.0.1:6379"},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"emails": 10,
			},
		},
	)

	mailClient := mail.NewMailClient(apiKey)

	mux := asynq.NewServeMux()

	worker := &EmailWorker{
		server:     server,
		mux:        mux,
		mailClient: mailClient,
	}

	worker.registerHandlers()

	return worker
}

func (w *EmailWorker) registerHandlers() {
	w.mux.HandleFunc(TypeEmailSend, func(ctx context.Context, t *asynq.Task) error {
		var payload EmailJob

		if err := json.Unmarshal(t.Payload(), &payload); err != nil {
			return err
		}

		fmt.Println("sending mail", payload)

		return w.mailClient.SendMail(ctx, payload.From, payload.To, payload.Subject, payload.Html)
	})
}

func (w *EmailWorker) Start() {
	fmt.Println("running email worker")
	if err := w.server.Run(w.mux); err != nil {
		log.Fatal(err)
	}
}

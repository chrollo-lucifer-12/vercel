package main

import (
	"github.com/chrollo-lucifer-12/shared/env"
	"github.com/chrollo-lucifer-12/shared/queue"
)

func main() {
	env.Load()

	worker := queue.NewEmailWorkerServer(env.RedisUrl.GetValue(), env.ResendApiKey.GetValue())
	worker.Start()
}

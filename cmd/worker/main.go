package main

import (
	"github.com/lPoltergeist/rinha-backend.git/data"
	worker "github.com/lPoltergeist/rinha-backend.git/workers"
)

func main() {
	data.InitRedis()

	worker.InitWorkers(190)

	select {}
}

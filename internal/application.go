package internal

import (
	"temporal-getting-started/adapter/in/process"
	"temporal-getting-started/adapter/in/rest"
)

type Application struct {
	Worker *process.Worker
	Server *rest.ServerHTTP
}

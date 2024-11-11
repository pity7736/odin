package main

import (
	"raiseexception.dev/odin/src/shared/infrastructure/api"
	"raiseexception.dev/odin/src/shared/infrastructure/repositoryfactory"
)

func main() {
	app := api.NewFiberApplication(repositoryfactory.New())
	app.Start()
}

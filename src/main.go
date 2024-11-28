package main

import (
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
	"raiseexception.dev/odin/src/app"
)

func main() {
	app := app.NewFiberApplication(accountingrepositoryfactory.New())
	app.Start()
}

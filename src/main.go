package main

import (
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/accountsrepositoryfactory"
	"raiseexception.dev/odin/src/app"
)

func main() {
	app := app.NewFiberApplication(
		accountingrepositoryfactory.New(),
		accountsrepositoryfactory.New(),
	)
	app.Start()
}

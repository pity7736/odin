package main

import (
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/accountsrepositoryfactory"
	"raiseexception.dev/odin/src/app"
)

func main() {
	application := app.NewFiberApplication(
		accountingrepositoryfactory.New(),
		accountsrepositoryfactory.New(),
	)
	if err := application.Start(); err != nil {
		panic(err)
	}
}

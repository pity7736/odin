package main

import (
	"raiseexception.dev/odin/src/accounting/infrastructure/repositories/accountingrepositoryfactory"
	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/pgrepositories"
	"raiseexception.dev/odin/src/app"
)

func main() {
	application := app.NewFiberApplication(
		accountingrepositoryfactory.New(),
		pgrepositories.NewPGSessionRepository(),
		pgrepositories.NewPGUserRepository(),
	)
	if err := application.Start(); err != nil {
		panic(err)
	}
}

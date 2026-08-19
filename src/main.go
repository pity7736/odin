package main

import (
	"raiseexception.dev/odin/src/accounts/infrastructure/repositories/inmemory"
	"raiseexception.dev/odin/src/accounts/infrastructure/security/bcrypthasher"
	"raiseexception.dev/odin/src/app"
)

func main() {
	application := app.NewFiberApplication(
		inmemory.NewInMemorySessionRepository(),
		inmemory.NewInMemoryUserRepository(),
		bcrypthasher.New(),
	)
	if err := application.Start(); err != nil {
		panic(err)
	}
}

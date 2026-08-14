package handler

import "github.com/gofiber/fiber/v2"

const SessionName = "__Secure-odin-session"

type localsKey string

const SessionTokenKey localsKey = "sessionToken"

type Handler interface {
	Handle(ctx *fiber.Ctx) error
}

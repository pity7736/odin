package handler

import "github.com/gofiber/fiber/v2"

const SessionName = "__Secure-odin-session"

type Handler interface {
	Handle(ctx *fiber.Ctx) error
}

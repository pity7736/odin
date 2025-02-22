package api

import "github.com/gofiber/fiber/v2"

type Handler interface {
	Handle(ctx *fiber.Ctx) error
}

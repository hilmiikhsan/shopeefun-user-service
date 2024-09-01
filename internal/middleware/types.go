package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Locals struct {
	UserId string
	Role   string
}

func GetLocals(c *fiber.Ctx) *Locals {
	var locals = Locals{}

	userId, ok := c.Locals("user_id").(string)
	if ok {
		locals.UserId = userId
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get user_id from locals")
	}

	role, ok := c.Locals("role").(string)
	if ok {
		locals.Role = role
	} else {
		log.Warn().Msg("middleware::Locals-GetLocals failed to get role from locals")
	}

	return &locals
}

func (l *Locals) GetUserId() string {
	return l.UserId
}

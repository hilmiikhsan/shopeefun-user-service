package route

import (
	"github.com/gofiber/fiber/v2"
	integration "github.com/hilmiikhsan/shopeefun-user-service/internal/integration/oauth2google"
	handlerUser "github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/handler/rest"
	"github.com/hilmiikhsan/shopeefun-user-service/pkg/response"
	"github.com/rs/zerolog/log"
)

func SetupRoutes(app *fiber.App) {
	var (
		api         = app.Group("/api/v1")
		googleOauth = integration.NewOauth2googleIntegration()
	)

	handlerUser.NewUserHandler(googleOauth).Register(api)

	// fallback route
	app.Use(func(c *fiber.Ctx) error {
		var (
			method = c.Method()                       // get the request method
			path   = c.Path()                         // get the request path
			query  = c.Context().QueryArgs().String() // get all query params
			ua     = c.Get("User-Agent")              // get the request user agent
			ip     = c.IP()                           // get the request IP
		)

		log.Info().
			Str("url", c.OriginalURL()).
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Str("ua", ua).
			Str("ip", ip).
			Msg("Route not found.")
		return c.Status(fiber.StatusNotFound).JSON(response.Error("Route not found"))
	})
}

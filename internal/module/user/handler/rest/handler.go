package rest

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/adapter"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/middleware"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/entity"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/ports"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/repository"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/service"
	"github.com/hilmiikhsan/shopeefun-user-service/pkg/errmsg"
	"github.com/hilmiikhsan/shopeefun-user-service/pkg/response"
	"github.com/rs/zerolog/log"
)

type userHandler struct {
	service ports.UserService
}

func NewUserHandler() *userHandler {
	var handler = new(userHandler)

	repo := repository.NewUserRepository(adapter.Adapters.ShopeefunPostgres)
	service := service.NewUserService(repo)

	handler.service = service

	return handler
}

func (h *userHandler) Register(router fiber.Router) {
	router.Post("/register", h.register)
	router.Post("/login", h.login)
	router.Get("/profile", middleware.AuthBearer, h.getProfile)
}

func (h *userHandler) register(c *fiber.Ctx) error {
	var (
		req        = new(entity.RegisterRequest)
		ctx        = c.Context()
		validators = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::register - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := validators.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::register - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Register(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::register - Failed to register user")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *userHandler) login(c *fiber.Ctx) error {
	var (
		req        = new(entity.LoginRequest)
		ctx        = c.Context()
		validators = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := validators.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::login - Invalid request body")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.Login(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::login - Failed to login user")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) getProfile(c *fiber.Ctx) error {
	var (
		req    = new(entity.GetProfileRequest)
		ctx    = c.Context()
		locals = middleware.GetLocals(c)
	)

	req.UserId = locals.GetUserId()

	res, err := h.service.GetProfile(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::getProfile - Failed to get profile")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

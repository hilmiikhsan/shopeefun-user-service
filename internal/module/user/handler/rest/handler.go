package rest

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/adapter"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/infrastructure/config"
	integOauth "github.com/hilmiikhsan/shopeefun-user-service/internal/integration/oauth2google"
	oauth "github.com/hilmiikhsan/shopeefun-user-service/internal/integration/oauth2google/entity"
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
	service     ports.UserService
	integration integOauth.Oauth2googleContract
}

func NewUserHandler(oauth integOauth.Oauth2googleContract) *userHandler {
	var handler = new(userHandler)

	repo := repository.NewUserRepository(adapter.Adapters.ShopeefunPostgres)
	service := service.NewUserService(repo, oauth)

	handler.integration = oauth
	handler.service = service

	return handler
}

func (h *userHandler) Register(router fiber.Router) {
	router.Post("/register", h.register)
	router.Post("/login", h.login)
	router.Get("/profile", middleware.AuthBearer, h.getProfile)
	router.Get("/profile/:user_id", middleware.AuthBearer, h.getProfileByUserId)

	router.Get("/oauth/google/url", h.oauthGoogleUrl)
	router.Get("/signin/callback", h.callbackSigninGoogle)
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

func (h *userHandler) getProfileByUserId(c *fiber.Ctx) error {
	var (
		req        = new(entity.GetProfileRequest)
		ctx        = c.Context()
		validators = adapter.Adapters.Validator
	)

	req.UserId = c.Params("user_id")

	if err := validators.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::profileByUserId - Invalid Request")
		code, errs := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.GetProfile(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("handler::profileByUserId - Failed to get profile")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *userHandler) oauthGoogleUrl(c *fiber.Ctx) error {
	return c.Redirect(h.integration.GetUrl("/"), http.StatusTemporaryRedirect)
}

func (h *userHandler) callbackSigninGoogle(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	state, code := c.FormValue("state"), c.FormValue("code")
	if state == "" && code == "" {
		log.Error().Msg("handler::callbackSigninGoogle - Invalid request")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errmsg.NewCustomErrors(fiber.StatusBadRequest, errmsg.WithMessage("Invalid request"))))
	}

	token, err := h.integration.Exchange(ctx, code)
	if err != nil {
		log.Error().Err(err).Msg("handler::callbackSigninGoogle - Failed to exchange token")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		log.Error().Err(err).Msg("handler::callbackSigninGoogle - Failed to get provider")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: config.Envs.Oauth.Google.ClientId,
	})

	_, err = verifier.Verify(context.Background(), token.Extra("id_token").(string))
	if err != nil {
		log.Error().Err(err).Msg("handler::callbackSigninGoogle - Failed to verify token")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	result, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("handler::callbackSigninGoogle - Failed to get user info")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}
	defer result.Body.Close()

	var userInfo oauth.UserInfoResponse
	if err := json.NewDecoder(result.Body).Decode(&userInfo); err != nil {
		log.Error().Err(err).Msg("handler::callbackSigninGoogle - Failed to decode user info")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.LoginWithGoogle(ctx, &userInfo)
	if err != nil {
		log.Error().Err(err).Any("payload", userInfo).Msg("handler::callbackSigninGoogle - Failed to login with google")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

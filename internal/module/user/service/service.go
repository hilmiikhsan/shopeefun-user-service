package service

import (
	"context"
	"time"

	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/entity"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/ports"
	"github.com/hilmiikhsan/shopeefun-user-service/pkg"
	"github.com/hilmiikhsan/shopeefun-user-service/pkg/errmsg"
	jwthandler "github.com/hilmiikhsan/shopeefun-user-service/pkg/jwt_handler"
	"github.com/rs/zerolog/log"
)

var _ ports.UserService = &userService{}

type userService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *userService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {
	hashed, err := pkg.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to hash password")
		return nil, errmsg.NewCustomErrors(500, errmsg.WithMessage("Gagal menghash password"))
	}

	req.HassedPassword = hashed

	result, err := s.repo.Register(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to register user")
		return nil, err
	}

	return result, nil
}

func (s *userService) Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error) {
	var res = new(entity.LoginResponse)

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to find user")
		return nil, err
	}

	if !pkg.ComparePassword(user.Password, req.Password) {
		log.Warn().Any("payload", req).Msg("service::Login - Password not match")
		return nil, errmsg.NewCustomErrors(401, errmsg.WithMessage("Email atau password salah"))
	}

	token, err := jwthandler.GenerateTokenString(jwthandler.CostumClaimsPayload{
		UserId:          user.Id,
		Role:            user.Role,
		TokenExpiration: time.Now().Add(time.Hour * 24),
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Login - Failed to generate token")
		return nil, err
	}

	res.Id = user.Id
	res.Token = token

	return res, nil
}

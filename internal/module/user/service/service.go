package service

import (
	"context"

	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/entity"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/ports"
	"github.com/hilmiikhsan/shopeefun-user-service/pkg"
	"github.com/hilmiikhsan/shopeefun-user-service/pkg/errmsg"
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

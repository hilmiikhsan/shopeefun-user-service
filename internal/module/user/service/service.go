package service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/infrastructure/config"
	integOauth "github.com/hilmiikhsan/shopeefun-user-service/internal/integration/oauth2google"
	oauthgoogleent "github.com/hilmiikhsan/shopeefun-user-service/internal/integration/oauth2google/entity"
	role "github.com/hilmiikhsan/shopeefun-user-service/internal/module/role/entity"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/entity"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/ports"
	"github.com/hilmiikhsan/shopeefun-user-service/pkg"
	"github.com/hilmiikhsan/shopeefun-user-service/pkg/errmsg"
	jwthandler "github.com/hilmiikhsan/shopeefun-user-service/pkg/jwt_handler"
	"github.com/rs/zerolog/log"
)

var _ ports.UserService = &userService{}

type userService struct {
	repo  ports.UserRepository
	oauth integOauth.Oauth2googleContract
}

func NewUserService(repo ports.UserRepository, oauth integOauth.Oauth2googleContract) *userService {
	return &userService{
		repo:  repo,
		oauth: oauth,
	}
}

func (s *userService) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {
	hashed, err := pkg.HashPassword(req.Password)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to hash password")
		return nil, errmsg.NewCustomErrors(fiber.StatusInternalServerError, errmsg.WithMessage("Gagal menghash password"))
	}

	req.HassedPassword = hashed

	result, err := s.repo.Register(ctx, req)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::Register - Failed to register user")
		return nil, err
	}

	return &entity.RegisterResponse{
		Id:   result.Id,
		Name: result.Name,
		Role: role.Role{
			Id:   result.RoleId,
			Name: result.RoleName,
		},
	}, nil
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
		return nil, errmsg.NewCustomErrors(fiber.StatusUnauthorized, errmsg.WithMessage("Email atau password salah"))
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

func (s *userService) GetProfile(ctx context.Context, req *entity.GetProfileRequest) (*entity.GetProfileResponse, error) {
	user, err := s.repo.FindById(ctx, req.UserId)
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::GetProfile - Failed to get user")
		return nil, err
	}

	return &entity.GetProfileResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role: role.Role{
			Id:   user.RoleId,
			Name: user.RoleName,
		},
	}, nil
}

func (s *userService) LoginWithGoogle(ctx context.Context, req *oauthgoogleent.UserInfoResponse) (*entity.LoginResponse, error) {
	var res = new(entity.LoginResponse)

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errCostum, ok := err.(*errmsg.CustomError); ok {
			if errCostum.Code != fiber.StatusBadRequest {
				log.Error().Err(err).Any("payload", req).Msg("service::LoginWithGoogle[1] - Failed to find user")
				return nil, err
			}

			hashed, err := pkg.HashPassword(config.Envs.Oauth.Google.OauthDefaultPassword)
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("service::LoginWithGoogle - Failed to hash password")
				return nil, errmsg.NewCustomErrors(fiber.StatusInternalServerError, errmsg.WithMessage("Gagal menghash password"))
			}

			result, err := s.repo.Register(ctx, &entity.RegisterRequest{
				Name:           req.Name,
				Email:          req.Email,
				HassedPassword: hashed,
			})
			if err != nil {
				log.Error().Err(err).Any("payload", req).Msg("service::LoginWithGoogle - Failed to register user")
				return nil, err
			}

			user = &entity.UserResult{
				Id:   result.Id,
				Role: result.RoleName,
			}
		}
	}

	token, err := jwthandler.GenerateTokenString(jwthandler.CostumClaimsPayload{
		UserId:          user.Id,
		Role:            user.Role,
		TokenExpiration: time.Now().Add(time.Hour * 24),
	})
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::LoginWithGoogle - Failed to generate token")
		return nil, err
	}

	res.Token = token

	return res, nil
}

func (s *userService) GetOauthGoogleUrl(ctx context.Context) (string, error) {
	url := s.oauth.GetUrl("state")

	return url, nil
}

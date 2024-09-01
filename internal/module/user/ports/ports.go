package ports

import (
	"context"

	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/entity"
)

type UserRepository interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error)
	FindByEmail(ctx context.Context, email string) (*entity.UserResult, error)
	FindById(ctx context.Context, id string) (*entity.GetProfileResponse, error)
}

type UserService interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error)
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error)
	GetProfile(ctx context.Context, req *entity.GetProfileRequest) (*entity.GetProfileResponse, error)
}

package ports

import (
	"context"

	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/entity"
)

type UserRepository interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error)
}

type UserService interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error)
}

package ports

import (
	"context"

	oauthgoogleent "github.com/hilmiikhsan/shopeefun-user-service/internal/integration/oauth2google/entity"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/entity"
)

type UserRepository interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResult, error)
	FindByEmail(ctx context.Context, email string) (*entity.UserResult, error)
	FindById(ctx context.Context, id string) (*entity.GetProfileResult, error)
}

type UserService interface {
	Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error)
	Login(ctx context.Context, req *entity.LoginRequest) (*entity.LoginResponse, error)
	GetProfile(ctx context.Context, req *entity.GetProfileRequest) (*entity.GetProfileResponse, error)
	LoginWithGoogle(ctx context.Context, req *oauthgoogleent.UserInfoResponse) (*entity.LoginResponse, error)
	GetOauthGoogleUrl(ctx context.Context) (string, error)
}

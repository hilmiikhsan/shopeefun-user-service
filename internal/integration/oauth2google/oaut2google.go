package integration

import (
	"context"
	"encoding/json"

	"github.com/hilmiikhsan/shopeefun-user-service/internal/infrastructure/config"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/integration/oauth2google/entity"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Oauth2googleContract interface {
	GetUrl(state string, opts ...oauth2.AuthCodeOption) string
	Exchange(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, token *oauth2.Token) (entity.UserInfoResponse, error)
}

type ouath2google struct {
	cfg oauth2.Config
}

func NewOauth2googleIntegration() *ouath2google {
	var googleOauthCfg = oauth2.Config{
		ClientID:     config.Envs.Oauth.Google.ClientId,
		ClientSecret: config.Envs.Oauth.Google.ClientSecret,
		RedirectURL:  config.Envs.Oauth.Google.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	return &ouath2google{
		cfg: googleOauthCfg,
	}
}

func (o *ouath2google) GetUrl(state string, opts ...oauth2.AuthCodeOption) string {
	return o.cfg.AuthCodeURL(state, opts...)
}

func (o *ouath2google) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return o.cfg.Exchange(ctx, code)
}

func (o *ouath2google) GetUserInfo(ctx context.Context, token *oauth2.Token) (entity.UserInfoResponse, error) {
	var (
		info   = entity.UserInfoResponse{}
		client = o.cfg.Client(ctx, token)
	)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		log.Error().Err(err).Msg("failed to get user info")
		return info, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		log.Error().Err(err).Msg("failed to decode user info")
		return info, err
	}

	return info, nil
}

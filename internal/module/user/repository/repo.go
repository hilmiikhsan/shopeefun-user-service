package repository

import (
	"context"

	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/entity"
	"github.com/hilmiikhsan/shopeefun-user-service/internal/module/user/ports"
	"github.com/hilmiikhsan/shopeefun-user-service/pkg/errmsg"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var _ ports.UserRepository = &userRepository{}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Register(ctx context.Context, req *entity.RegisterRequest) (*entity.RegisterResponse, error) {
	var res = new(entity.RegisterResponse)

	err := r.db.QueryRowContext(ctx, r.db.Rebind(queryInsertUser),
		req.Email,
		req.Name,
		req.HassedPassword,
	).Scan(
		&res.Id,
		&res.Name,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if !ok {
			log.Error().Err(err).Any("payload", req).Msg("repo::Register - Failed to insert user")
			return nil, err
		}

		switch pqErr.Code.Name() {
		case "unique_violation":
			log.Warn().Err(err).Any("payload", req).Msg("repo::Register - Email already registered")
			return nil, errmsg.NewCustomErrors(409, errmsg.WithMessage("Email sudah terdaftar"))
		default:
			log.Error().Err(err).Any("payload", req).Msg("repo::Register - Failed to insert user")
			return nil, err
		}
	}

	return res, nil
}

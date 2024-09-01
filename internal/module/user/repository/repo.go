package repository

import (
	"context"
	"database/sql"

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

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.UserResult, error) {
	var res = new(entity.UserResult)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryGetUserByEmail), email)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Str("email", email).Msg("repo::FindByEmail - User not found")
			return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("Email atau password salah"))
		}
		log.Error().Err(err).Str("email", email).Msg("repo::FindByEmail - Failed to get user")
		return nil, err
	}

	return res, nil
}

func (r *userRepository) FindById(ctx context.Context, id string) (*entity.GetProfileResponse, error) {
	var res = new(entity.GetProfileResponse)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryGetProfile), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Err(err).Str("id", id).Msg("repo::FindById - User not found")
			return nil, errmsg.NewCustomErrors(400, errmsg.WithMessage("User tidak ditemukan"))
		}

		log.Error().Err(err).Str("id", id).Msg("repo::FindById - Failed to get user")
		return nil, err
	}

	return res, nil
}

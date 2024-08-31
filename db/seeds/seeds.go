package seeds

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Seed struct {
	db *sqlx.DB
}

// NewSeed return a Seed with a pool of connection to a dabase.
func newSeed(db *sqlx.DB) Seed {
	return Seed{
		db: db,
	}
}

func Execute(db *sqlx.DB, table string, total int) {
	seed := newSeed(db)
	seed.run(table, total)
}

// Run seeds.
func (s *Seed) run(table string, total int) {
	switch table {
	case "roles":
		s.rolesSeed()
	case "users":
		s.usersSeed(total)
	case "all":
		s.rolesSeed()
		s.usersSeed(total)
	case "delete-all":
		s.deleteAll()
	default:
		log.Warn().Msg("No seed to run")
	}
}

// rolesSeed seeds the roles table.
func (s *Seed) rolesSeed() {
	roleMaps := []map[string]any{
		{"name": "admin"},
		{"name": "end_user"},
	}

	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	_, err = tx.NamedExec(insertRoleName, roleMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating roles")
		return
	}

	log.Info().Msg("roles table seeded successfully")
}

// users
func (s *Seed) usersSeed(total int) {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		}

		err = tx.Commit()
		if err != nil {
			log.Error().Err(err).Msg("Error committing transaction")
		}
	}()

	type generalData struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	var (
		roles    = make([]generalData, 0)
		userMaps = make([]map[string]any, 0)
	)

	err = s.db.Select(&roles, getRoles)
	if err != nil {
		log.Error().Err(err).Msg("Error selecting roles")
		return
	}

	for i := 0; i < total; i++ {
		selectedRole := roles[gofakeit.Number(0, len(roles)-1)]

		dataUserToInsert := make(map[string]any)
		// dataUserToInsert["id"] = ulid.Make().String()
		dataUserToInsert["id"] = uuid.New().String()
		dataUserToInsert["role_id"] = selectedRole.Id
		dataUserToInsert["name"] = gofakeit.Name()
		dataUserToInsert["email"] = gofakeit.Email()
		dataUserToInsert["whatsapp_number"] = gofakeit.Phone()
		dataUserToInsert["password"] = "$2y$10$mVf4BKsfPSh/pjgHjvk.JOlGdkIYgBGyhaU9WQNMWpYskK9MZlb0G" // password

		userMaps = append(userMaps, dataUserToInsert)
	}

	var (
		endUserId   string
		adminUserId string
	)

	// iterate over roles to get service advisor id
	for _, role := range roles {
		if role.Name == "admin" {
			adminUserId = role.Id
			continue
		}
		if role.Name == "end_user" {
			endUserId = role.Id
			continue
		}
	}

	EndUser := map[string]any{
		// "id":              ulid.Make().String(),
		"id":              uuid.New().String(),
		"role_id":         endUserId,
		"name":            "Ikhsan",
		"email":           "ikhsan@fake.com",
		"whatsapp_number": gofakeit.Phone(),
		"password":        "$2y$10$mVf4BKsfPSh/pjgHjvk.JOlGdkIYgBGyhaU9WQNMWpYskK9MZlb0G", // password
	}

	AdminUser := map[string]any{
		"id":              uuid.New().String(),
		"role_id":         adminUserId,
		"name":            "Hilmi",
		"email":           "hilmi@fake.com",
		"whatsapp_number": gofakeit.Phone(),
		"password":        "$2y$10$mVf4BKsfPSh/pjgHjvk.JOlGdkIYgBGyhaU9WQNMWpYskK9MZlb0G", // password
	}

	userMaps = append(userMaps, EndUser)
	userMaps = append(userMaps, AdminUser)

	_, err = tx.NamedExec(insertRoles, userMaps)
	if err != nil {
		log.Error().Err(err).Msg("Error creating users")
		return
	}

	log.Info().Msg("users table seeded successfully")
}

func (s *Seed) deleteAll() {
	tx, err := s.db.BeginTxx(context.Background(), nil)
	if err != nil {
		log.Error().Err(err).Msg("Error starting transaction")
		return
	}

	defer func() {
		if err != nil {
			err = tx.Rollback()
			log.Error().Err(err).Msg("Error rolling back transaction")
			return
		} else {
			err = tx.Commit()
			if err != nil {
				log.Error().Err(err).Msg("Error committing transaction")
			}
		}
	}()

	_, err = tx.Exec(deleteUsers)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting users")
		return
	}
	log.Info().Msg("users table deleted successfully")

	_, err = tx.Exec(deleteRoles)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting roles")
		return
	}
	log.Info().Msg("roles table deleted successfully")

	log.Info().Msg("=== All tables deleted successfully ===")
}

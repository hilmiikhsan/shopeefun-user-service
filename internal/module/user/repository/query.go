package repository

const (
	queryInsertUser = `
		INSERT INTO users
		(
			role_id,
			email,
			name,
			password
		) VALUES (
			(SELECT id FROM roles WHERE name = 'end_user'),
			?, ?, ?
		)

		RETURNING id, name
	`
)

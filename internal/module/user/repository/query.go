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

	queryGetUserByEmail = `
		SELECT
			u.id,
			r.name as role,
			u.name,
			u.email,
			u.password
		FROM users u
		LEFT JOIN roles r ON r.id = u.role_id
		WHERE u.email = ?
	`

	queryGetProfile = `
		SELECT
			u.id,
			r.name AS role,
			u.name,
			u.email
		FROM users u
		LEFT JOIN roles r ON r.id = u.role_id
		WHERE u.id = ?
	`
)

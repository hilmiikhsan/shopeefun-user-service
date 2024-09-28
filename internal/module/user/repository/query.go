package repository

const (
	queryInsertUser = `
		WITH inserted_user AS (
			INSERT INTO users (
				role_id,
				email,
				name,
				password
			) VALUES (
				(SELECT id FROM roles WHERE name = 'end_user'),
				?, ?, ?
			)
			RETURNING id, role_id, name
		)
		SELECT 
			inserted_user.id, 
			inserted_user.name, 
			roles.id AS role_id,
			roles.name AS role_name
		FROM 
			inserted_user
		JOIN 
			roles ON roles.id = inserted_user.role_id
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
			r.id AS role_id,
			r.name AS role_name,
			u.name,
			u.email
		FROM users u
		LEFT JOIN roles r ON r.id = u.role_id
		WHERE u.id = ?
	`
)

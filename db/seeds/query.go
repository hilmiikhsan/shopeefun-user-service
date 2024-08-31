package seeds

const (
	insertRoleName = `
		INSERT INTO roles (name) VALUES (:name)
	`

	getRoles = `
		SELECT id, name FROM roles
	`

	insertRoles = `
		INSERT INTO users (
			id, 
			role_id, 
			name, 
			email, 
			whatsapp_number, 
			password
		) VALUES (:id, :role_id, :name, :email, :whatsapp_number, :password)
	`

	deleteUsers = `
		DELETE FROM users
	`

	deleteRoles = `
		DELETE FROM roles
	`
)

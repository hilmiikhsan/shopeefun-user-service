package entity

type UserResult struct {
	Id       string `db:"id"`
	Role     string `db:"role"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

type RegisterResult struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	RoleId   string `db:"role_id"`
	RoleName string `db:"role_name"`
}

type GetProfileResult struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	RoleId   string `db:"role_id"`
	RoleName string `db:"role_name"`
}

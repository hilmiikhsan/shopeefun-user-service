package entity

type UserResult struct {
	Id       string `db:"id"`
	Role     string `db:"role"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	Password string `db:"password"`
}

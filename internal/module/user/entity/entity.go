package entity

type RegisterRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`

	HassedPassword string
}

type RegisterResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

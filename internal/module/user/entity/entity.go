package entity

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`

	HassedPassword string
}

type RegisterResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

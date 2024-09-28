package entity

import "github.com/hilmiikhsan/shopeefun-user-service/internal/module/role/entity"

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`

	HassedPassword string
}

type RegisterResponse struct {
	Id   string      `json:"id"`
	Name string      `json:"name"`
	Role entity.Role `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

type GetProfileRequest struct {
	UserId string `validate:"required"`
}

type GetProfileResponse struct {
	Id    string      `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
	Role  entity.Role `json:"role"`
}

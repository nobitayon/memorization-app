package model

import "github.com/google/uuid"

// UserService defines methods the handler layer
// expects any service t interacts with to implement
type UserService interface {
	Get(uid uuid.UUID) (*User, error)
}

type UserRepository interface {
	FindByID(uid uuid.UUID) (*User, error)
}

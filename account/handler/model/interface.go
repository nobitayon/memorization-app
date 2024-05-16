package model

import (
	"context"

	"github.com/google/uuid"
)

// UserService defines methods the handler layer
// expects any service t interacts with to implement
type UserService interface {
	Get(ctx context.Context, uid uuid.UUID) (*User, error)
}

type UserRepository interface {
	FindByID(ctx context.Context, uid uuid.UUID) (*User, error)
}

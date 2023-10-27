package models

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	CreateCompany(ctx context.Context, ne Company) (Company, error)
	ViewCompany(ctx context.Context, userId string) (float64, error)
	// CreateEmployee(ctx context.Context, ne newEmp) (Employee, error)
	CreateUser(ctx context.Context, nu NewUser) (User, error)
	Authenticate(ctx context.Context, email, password string) (jwt.RegisteredClaims,
		error)
}

type Store struct {
	Service
}

func NewStore(s Service) Store {
	return Store{Service: s}
}

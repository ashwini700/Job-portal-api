package models

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	JobbyCompId(jobs Job, compId string) (Job, error)
	GetJobById(ctx context.Context, jobId string) (Job, error)
	FetchJobByCompanyId(ctx context.Context, companyId string) ([]Job, error)
	GetAllJobs(ctx context.Context) ([]Job, error)
	CreateCompany(ctx context.Context, ne Company) (Company, error)
	ViewCompany(ctx context.Context) ([]Company, error)
	FetchcompanyID(ctx context.Context, companyID string) (Company, error)
	CreateUser(ctx context.Context, nu NewUser) (User, error) //for signup and login
	Authenticate(ctx context.Context, email, password string) (jwt.RegisteredClaims,
		error)
}

type Store struct {
	Service
}

func NewStore(s Service) Store {
	return Store{Service: s}
}

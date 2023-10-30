package models

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Conn is our main struct, including the database instance for working with data.
type Conn struct {

	// db is an instance of the SQLite database.
	db *gorm.DB
}

// Add jobs

func (s *Conn) JobbyCompId(j Job,compid string) (Job, error) {
	compId, _ := strconv.ParseUint(compid, 10, 64)

	job := Job{
		Name : j.Name,
		Field: j.Field,
		Experience:j.Experience, 
		CompanyId:  compId,
	}
	err := s.db.Create(&job).Error
	if err != nil {
		return Job{}, err
	}
	return job, nil
}

// func (s *Conn) JobbyCompId(jobs []Job, compId string) ([]Job, error) {
// 	companyId, err := strconv.ParseUint(compId, 10, 64)
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, j := range jobs {
// 		job := Job{
// 			Name:       j.Name,
// 			Field:      j.Field,
// 			Experience: j.Experience,
// 			CompanyId:  companyId,
// 		}
// 		err := s.db.Create(&job).Error
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return jobs, nil
// }

// job by company id
func (s *Conn) FetchJobByCompanyId(ctx context.Context, companyId string) ([]Job, error) {
	var listOfJobs []Job
	tx := s.db.WithContext(ctx).Where("company_id = ?", companyId)
	err := tx.Find(&listOfJobs).Error
	if err != nil {
		return nil, err
	}

	return listOfJobs, nil
}

// Job by id
func (s *Conn) GetJobById(ctx context.Context, jobId string) (Job, error) {
	var jobData Job
	tx := s.db.WithContext(ctx).Where("ID = ?", jobId)
	err := tx.Find(&jobData).Error
	if err != nil {
		return Job{}, err
	}

	return jobData, nil
}

func (s *Conn) GetAllJobs(ctx context.Context) ([]Job, error) {
	var listJobs []Job
	tx := s.db.WithContext(ctx)
	err := tx.Find(&listJobs).Error
	if err != nil {
		return nil, err
	}

	return listJobs, nil
}

// ViewCompany implements Service.
func (s *Conn) ViewCompany(ctx context.Context) ([]Company, error) {
	var companies []Company
	tx := s.db.WithContext(ctx)
	err := tx.Find(&companies).Error
	if err != nil {
		return nil, err
	}

	return companies, nil
}

func (s *Conn) FetchcompanyID(ctx context.Context, companyId string) (Company, error) {
	var comp Company
	tx := s.db.WithContext(ctx).Where("id = ?", companyId)
	err := tx.Find(&comp).Error
	if err != nil {
		return Company{}, err
	}
	return comp, nil
}

// NewService is the constructor for the Conn struct.
func NewService(db *gorm.DB) (*Conn, error) {

	// We check if the database instance is nil, which would indicate an issue.
	if db == nil {
		return nil, errors.New("please provide a valid connection")
	}

	// We initialize our service with the passed database instance.
	s := &Conn{db: db}
	return s, nil
}

// CreateUser is a method that creates a new user record in the database.
func (s *Conn) CreateUser(ctx context.Context, nu NewUser) (User, error) {

	// We hash the user's password for storage in the database.
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	// We prepare the User record.
	u := User{
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: string(hashedPass),
	}

	// We attempt to create the new User record in the database.
	err = s.db.Create(&u).Error
	if err != nil {
		return User{}, err
	}

	// Successfully created the record, return the user.
	return u, nil
}

// CreateCompany implements Service.
func (s *Conn) CreateCompany(ctx context.Context, ne Company) (Company, error) {
	//prepare company record
	c := Company{
		CompanyName: ne.CompanyName,
		Location:    ne.Location,
	}

	//create new company database
	err := s.db.Create(&c).Error
	if err != nil {
		return Company{}, err
	}
	return c, nil
}

// Authenticate is a method that checks a user's provided email and password against the database.
func (s *Conn) Authenticate(ctx context.Context, email, password string) (jwt.RegisteredClaims,
	error) {

	// We attempt to find the User record where the email
	// matches the provided email.
	var u User
	tx := s.db.Where("email = ?", email).First(&u)
	if tx.Error != nil {
		return jwt.RegisteredClaims{}, tx.Error
	}

	// We check if the provided password matches the hashed password in the database.
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return jwt.RegisteredClaims{}, err
	}

	// Successful authentication! Generate JWT claims.
	c := jwt.RegisteredClaims{
		Issuer:    "service project",
		Subject:   strconv.FormatUint(uint64(u.ID), 10),
		Audience:  jwt.ClaimStrings{"employees"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	// And return those claims.
	return c, nil
}

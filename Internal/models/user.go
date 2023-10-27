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

// ViewCompany implements Service.
func (*Conn) ViewCompany(ctx context.Context, userId string) (float64, error) {
	panic("unimplemented")
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
		JobRole:     ne.JobRole,
	}

	//create new company database
	err := s.db.Create(&c).Error
	if err != nil {
		return Company{}, err
	}
	return c, nil
}

// ViewCompany implements Service.
// func (com *Conn) ViewCompany(ctx context.Context, Id string) ([]Company, float64, error) {
// 	var comp = make([]Company, 0, 10)
// 	tx := com.db.WithContext(ctx).Where("Id = ?", Id)
// 	err := tx.Find(&comp).Error
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	totalCompanies, err := CalculateTotalComp(comp, "companies")
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return comp, totalCompanies, nil
// }

// func CalculateTotalComp(CompanyName []Company, category string) (float64, error) {
// 	if CompanyName == nil {
// 		return 0, errors.New("company doesn't exist")
// 	}
// 	// Compute the total cost
// 	var totalCompanies float64
// 	for _, company := range CompanyName {
// 		totalCompanies += float64(company.ID)
// 	}
// 	return totalCompanies, nil
// }

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

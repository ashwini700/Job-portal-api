package handlers

import (
	"Job-portal-api/Internal/auth"
	"Job-portal-api/Internal/middleware"
	"Job-portal-api/Internal/models"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type handler struct {
	s models.Store
	a *auth.Auth
}

// function to add or register companies
func (h *handler) AddCompany(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		// If the traceId isn't found in the request, log an error and return
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	var ne models.Company
	err := json.NewDecoder(c.Request.Body).Decode(&ne)
	if err != nil {
		// If there is an error in decoding, log the error and return
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	//validate company and details
	validate := validator.New()
	err = validate.Struct(&ne)
	if err != nil {
		// If validation fails, log the error and return
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide Name, ID and Location"})
		return
	}

	//attempt to add new company
	usr, err := h.s.CreateCompany(ctx, ne)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("user company signup problem")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Company signup failed"})
		return
	}

	c.JSON(http.StatusOK, usr)
}

// Signup is a method for the handler struct which handles user registration
func (h *handler) Signup(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		// If the traceId isn't found in the request, log an error and return
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// Define a NewUser variable
	var nu models.NewUser

	// Attempt to decode JSON from the request body into the NewUser variable
	err := json.NewDecoder(c.Request.Body).Decode(&nu)
	if err != nil {
		// If there is an error in decoding, log the error and return
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// Create a new validator and validate the NewUser variable
	validate := validator.New()
	err = validate.Struct(&nu)
	if err != nil {
		// If validation fails, log the error and return
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide Name, Email and Password"})
		return
	}

	// Attempt to create the user
	usr, err := h.s.CreateUser(ctx, nu)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Msg("user signup problem")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "user signup failed"})
		return
	}

	// If everything goes right, respond with the created user
	c.JSON(http.StatusOK, usr)
}

// Login is a method for the handler struct which handles user login
func (h *handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// Define a new struct for login data
	var login struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	// Attempt to decode JSON from the request body into the login variable
	err := json.NewDecoder(c.Request.Body).Decode(&login)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// Create a new validator and validate the login variable
	validate := validator.New()
	err = validate.Struct(login)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide Email and Password"})
		return
	}

	// Attempt to authenticate the user with the email and password
	claims, err := h.s.Authenticate(ctx, login.Email, login.Password)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId).Send()
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"msg": "login failed"})
		return
	}

	// Define a new struct for the token
	var tkn struct {
		Token string `json:"token"`
	}

	// Generate a new token and put it in the Token field of the token struct
	tkn.Token, err = h.a.GenerateToken(claims)
	if err != nil {
		log.Error().Err(err).Msg("generating token")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	// If everything goes right, respond with the token
	c.JSON(http.StatusOK, tkn)

}

// fetch all the companies
func (h *handler) ViewCompany(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("traceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	companies, err := h.s.ViewCompany(ctx)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "problem in viewing companies"})
		return
	}
	// m := gin.H{"c": companies}
	c.JSON(http.StatusOK, companies)
}

// fetch company by id
func (h *handler) companyID(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		// If the traceId isn't found in the request, log an error and return
		log.Error().Msg("TrackerId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	companyId := c.Param("id")
	companySources, err := h.s.FetchcompanyID(ctx, companyId)
	if err != nil {
		log.Error().Err(err).Str("Tracker Id", traceId)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "problem in viewing company by id"})
		return
	}
	c.JSON(http.StatusOK, companySources)

}

// add jobs by Id
func (h *handler) addJobsById(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("TrackerId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	compId:=c.Param("id")
	var jobs models.Job
	err := json.NewDecoder(c.Request.Body).Decode(&jobs)
	if err != nil {
		// If there is an error in decoding, log the error and return
		log.Error().Err(err).Str("Tracker Id", traceId)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	jobData,err:=h.s.JobbyCompId(jobs,compId)
	if err != nil {
		log.Error().Err(err).Str("Tracker Id", traceId).Msg("Add Job by companyId problem")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "Job creation failed"})
		return
	}
	c.JSON(http.StatusCreated, jobData)
}
		////
func (h *handler) JobbyCompId(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("TrackerId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	companyId := c.Param("companyId")
	listOfJobs, err := h.s.FetchJobByCompanyId(ctx, companyId)
	if err != nil {
		log.Error().Err(err).Str("Tracker Id", traceId)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "problem in viewing list of company by ID"})
		return
	}
	c.JSON(http.StatusOK, listOfJobs)
}

func (h *handler) fetchJobById(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("TraceId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	jobId := c.Param("ID")
	job, err := h.s.GetJobById(ctx, jobId)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceId)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "problem in viewing list of company by ID"})
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *handler) GetAllJobs(c *gin.Context) {
	ctx := c.Request.Context()
	traceId, ok := ctx.Value(middleware.TraceIdKey).(string)
	if !ok {
		log.Error().Msg("TrackerId missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}
	job, err := h.s.GetAllJobs(ctx)
	if err != nil {
		log.Error().Err(err).Str("Tracker Id", traceId)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "problem in viewing list of company by ID"})
		return
	}
	c.JSON(http.StatusOK, job)
}

package handlers

import (
	"Job-portal-api/Internal/auth"
	"Job-portal-api/Internal/middleware"
	"Job-portal-api/Internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func API(a *auth.Auth, c *models.Conn) *gin.Engine {

	// Create a new Gin engine; Gin is a HTTP web framework written in Go
	r := gin.New()
	m, err := middleware.NewMid(a)
	ms := models.NewStore(c)
	h := handler{
		s: ms,
		a: a,
	}
	if err != nil {
		log.Panic().Msg("middlewares not set up")
	}
	r.Use(m.Log(), gin.Recovery())

	r.GET("/check", m.Authenticate(check))
	r.POST("/signup", h.Signup)
	r.POST("/login", h.Login)
	r.POST("/createcompany/add", h.AddCompany)
	r.GET("/viewcompany/view", h.ViewCompany)
	r.GET("/byIdcompany/:id", h.companyID)
	r.POST("/addJobs/:id", h.addJobsById)
	r.GET("/fetchJob/:id", h.fetchJobById)
	r.POST("/JobbyCompId/:id", h.JobbyCompId)
	r.GET("/getAllJob", h.GetAllJobs)

	return r
}

func check(c *gin.Context) {
	//handle panic using recovery function when happening in separate goroutine
	// go func() {
	// 	panic("some kind of panic")
	// }()
	time.Sleep(time.Second * 3)
	select {
	case <-c.Request.Context().Done():
		fmt.Println("user not there")
		return
	default:
		c.JSON(http.StatusOK, gin.H{"msg": "statusOk"})

	}
}

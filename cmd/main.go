package main

import (
	"log"
	"net/http"

	"github.com/AbdulrahmanDaud10/google-0auth2/pkg/app"
	"github.com/AbdulrahmanDaud10/google-0auth2/pkg/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var server *gin.Engine

func init() {
	repository.PostgresDatabaseConnection()

	server = gin.Default()
}

func main() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Implement Google OAuth2 in Golang"})
	})

	auth_router := router.Group("/auth")
	auth_router.POST("/register", app.SignUpUser)
	auth_router.POST("/login", app.SignInUser)
	auth_router.GET("/logout", app.DeserializeUser(), app.LogOutUser)

	router.GET("/sessions/oauth/google", app.GoogleOAuth)
	router.GET("/users/me", app.DeserializeUser(), app.GetAuthenticatedUser)

	router.StaticFS("/images", http.Dir("public"))
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Route Not Found"})
	})

	log.Fatal(server.Run(":" + "8000"))
}

package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AbdulrahmanDaud10/google-0auth2/pkg/api"
	"github.com/AbdulrahmanDaud10/google-0auth2/pkg/repository"
	"github.com/gin-gonic/gin"
)

// JWT middleware guard to ensure that access to protected resources is granted only when a valid JSON Web Token is provided.
func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var token string
		cookie, err := ctx.Cookie("token")

		// extract the JWT token from either the Authorization header or the Cookies object.
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else if err == nil {
			token = cookie
		}

		if token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		// validate the token using the secret key.
		config, err := repository.LoadConfig(".")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		sub, err := api.ValidateToken(token, config.JWTTokenSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		// Query the database to check if the user associated with the token still exists.
		var user api.User
		result := repository.DB.First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belongong to thos token no longer exists"})
			return
		}

		// The function will lastly add the query result to the Gin context object via the ctx.Set() method and forward the request to subsequent middleware.
		ctx.Set("currentUser", user)
		ctx.Next()
	}
}

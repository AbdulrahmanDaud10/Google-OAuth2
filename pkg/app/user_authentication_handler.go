package app

import (
	"net/http"

	"github.com/AbdulrahmanDaud10/google-0auth2/pkg/api"
	"github.com/gin-gonic/gin"
)

// GetAuthenticatedUser route function that will be protected by a JWT middleware guard.
// This route function will be called to return the currently logged-in user’s account information
// when a GET request is made to the /api/users/me endpoint.
func GetAuthenticatedUser(ctx *gin.Context) {
	// When Gin Gonic calls this route handler, it will extract the user’s profile information from
	// the context object using the ctx.MustGet() function and return the credentials in the JSON response.
	currentUser := ctx.MustGet("currentuser").(api.User)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": api.FilteredResponse(&currentUser)}})
}

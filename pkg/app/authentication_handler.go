package app

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AbdulrahmanDaud10/google-0auth2/pkg/api"
	"github.com/AbdulrahmanDaud10/google-0auth2/pkg/repository"
	"github.com/gin-gonic/gin"
)

// When this handler is triggered, it will parse the incoming data, validate the data against the rules defined in the models.
// `RegisterUserInput struct` save the user’s details in the database, and return a sanitized version of the record in the JSON response.
func SignUpUser(ctx *gin.Context) {
	var payload *api.RegisteredUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	newUser := api.User{
		Name:      payload.Name,
		Email:     strings.ToLower(payload.Email),
		Password:  payload.Password,
		Role:      "user",
		Verified:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// In case the email address submitted in the request already exists in the database, a 409 Conflict error response will be returned to the client.
	result := repository.DB.Create(&newUser)
	if result.Error != nil && strings.Contains(result.Error.Error(), "UNIQUE constraint failed: users.email") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User email already exists"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "something bad happened"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": api.FilteredResponse(&newUser)}})
}

// When this route function is triggered, it will parse the incoming request, validate the data based on
// the models.LoginUserInput struct, and check the database for the existence of a user with the email address provided in the request.
func SignInUser(ctx *gin.Context) {
	var payload *api.LoginUserInput

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	var user api.User
	result := repository.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid Email or Password"})
		return
	}

	// Once a matching user is found, the handler will determine if the account was registered via Google OAuth and, if so, a 403 Unauthorized error will be sent to the client.
	if user.Provider == "Google" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": fmt.Sprintf("Use %v Oauth instead", user.Provider)})
		return
	}

	config, err := repository.LoadConfig(".")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	token, err := api.GenerateToken(config.TokenExpiresIn, user.ID, config.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// To sign out a user, I Initiate a process to delete the current JWT token stored in the user’s API client or browser by sending an expired cookie.
func LogOutUser(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

// When this route function is triggered, it will extract the authorization code from the redirect URL and assign it to the code variable.
// After retrieving the authorization code, the route function will then extract the value of the “state” parameter from the redirect URL and store it in a variable named pathUrl.
func GoogleOAuth(ctx *gin.Context) {
	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Authorization code not provided!"})
		return
	}

	// GetGoogleOauthToken() function with the authorization code as its argument to retrieve the access token from the Google OAuth2 token endpoint.
	tokenRes, err := GetGoogleOauthToken(code)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// GetGoogleUser() function will be called to retrieve the user’s Google account information using the acquired access token.
	googleUser, err := GetGoogleUser(tokenRes.AccessToken, tokenRes.Id_token)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	now := time.Now()
	email := strings.ToLower(googleUser.Email)

	userData := api.User{
		Name:      googleUser.Name,
		Email:     email,
		Password:  "",
		Photo:     googleUser.Picture,
		Provider:  "Google",
		Role:      "user",
		Verified:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	//  If the request is successful, the user’s information will be either inserted into the database using GORM or the existing record will be updated with the latest details,
	// depending on whether the user already exists in the database.
	if repository.DB.Model(&userData).Where("email = ?", email).Updates(&userData).RowsAffected == 0 {
		repository.DB.Create(&userData)
	}

	var user api.User
	repository.DB.First(&user, "email = ?", email)

	config, _ := repository.LoadConfig(".")

	token, err := api.GenerateToken(config.TokenExpiresIn, user.ID.String(), config.JWTTokenSecret)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true)

	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.FrontEndOrigin, pathUrl))
}

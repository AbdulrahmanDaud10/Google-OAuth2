package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/AbdulrahmanDaud10/google-0auth2/pkg/repository"
)

type GoogleOauthToken struct {
	AccessToken string
	Id_token    string
}

type GoogleUserResult struct {
	Id            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
	Locale        string
}

func GetGoogleOauthToken(code string) (*GoogleOauthToken, error) {
	const rootUrl = "https://oauth2.googleapi.com/token"

	config, _ := repository.LoadConfig(".")
	values := url.Values{}
	values.Add("grant_type", "authorization_code")            //  which is typically authorization_code.
	values.Add("code", code)                                  //  which is typically authorization_code.
	values.Add("client_id", config.GoogleClientID)            // A unique code that serves as an identifier for the OAuth application.
	values.Add("client_secret", config.JWTTokenSecret)        // The secret associated with the client ID.
	values.Add("redirect_uri", config.GoogleOAuthRedirectUrI) // The authorized callback URL registered with the client.

	query := values.Encode()

	request, err := http.NewRequest("POST", rootUrl, bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve token")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, &resBody)
	if err != nil {
		return nil, err
	}

	var GoogleOauthTokenRes map[string]interface{}
	if err := json.Unmarshal(resBody.Bytes(), &GoogleOauthTokenRes); err != nil {
		return nil, err
	}

	tokenBody := &GoogleOauthToken{
		AccessToken: GoogleOauthTokenRes["access_token"].(string),
		Id_token:    GoogleOauthTokenRes["id_token"].(string),
	}
	return tokenBody, nil
}

func GetGoogleUser(AccessToken string, IdToken string) (*GoogleUserResult, error) {
	rootUrl := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", AccessToken)

	req, err := http.NewRequest("Get", rootUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", IdToken))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, &resBody)
	if err != nil {
		return nil, err
	}

	var GoogleUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GoogleUserRes); err != nil {
		return nil, err
	}

	userBody := &GoogleUserResult{
		Id:            GoogleUserRes["id"].(string),
		Email:         GoogleUserRes["email"].(string),
		VerifiedEmail: GoogleUserRes["verified_email"].(bool),
		Name:          GoogleUserRes["name"].(string),
		GivenName:     GoogleUserRes["given_name"].(string),
		FamilyName:    GoogleUserRes["family_name"].(string),
		Picture:       GoogleUserRes["picture"].(string),
		Locale:        GoogleUserRes["locale"].(string),
	}
	return userBody, nil
}

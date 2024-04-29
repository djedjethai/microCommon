package apiserver

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	e "gitlab.com/grpasr/common/errors/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AuthCredentials struct {
	AppUrlParams string
	JWT_token    string
}

type CodeStruct struct {
	Code string `json:"code"`
}

type APIserverAuth struct {
	authServerURL       string
	path                string
	redirectURI         string
	urlOauthToken       string
	codeVerifier        string
	clientID            string
	clientSecret        string
	scope               string
	state               string
	grantType           string
	role                string
	acceptEncoding      string
	contentType         string
	method              string
	codeChallengeMethod string
	authCredentials     *AuthCredentials
	codeStruct          *CodeStruct
}

func NewAPIserverAuth(authServerURL, path, redirectURI, urlOauthToken, codeVerifier, clientID, clientSecret, scope string) APIserverAuth {
	return APIserverAuth{
		authServerURL:       authServerURL, // "http://localhost:9096/v1"
		path:                path,          // "apiauth"
		redirectURI:         redirectURI,   // "http://localhost:50001"
		urlOauthToken:       urlOauthToken, // "http://localhost:9096/v1/oauth/token"
		codeVerifier:        codeVerifier,  // "exampleCodeVerifier"
		clientID:            clientID,      // "order" // must match the auth_svc
		clientSecret:        clientSecret,  // "orderSecret" // must match the auth_svc
		scope:               scope,         // "read, openid" // openid must be specify
		state:               "xxyyzz",      // or whatever
		grantType:           "authorization_code",
		role:                "APIserver",
		acceptEncoding:      "gzip",
		contentType:         "application/x-www-form-urlencoded",
		method:              "POST",
		codeChallengeMethod: "S256",
		authCredentials:     &AuthCredentials{},
	}
}

func (a APIserverAuth) Run(ctx context.Context, retry, d int8) e.IError {
	err := a.SetURLAndCreateCodeChallenge()
	if err != nil {
		return err
	}

	// retry pattern
	delay := time.Duration(int64(d))
	baseDelay := delay
	for r := int8(0); ; r++ {
		err = a.QueryToken()
		if err == nil || r >= retry {
			return err
		}

		// Calculate the delay using exponential backoff
		delay = time.Duration(baseDelay*time.Duration(r+1)) * time.Second
		log.Printf("Attempt %d failed; retrying in %v", r+1, delay)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return e.NewCustomHTTPStatus(e.StatusBadRequest, "", err.Error())
		}
	}
}

func (a APIserverAuth) QueryToken() e.IError {
	u, err := url.Parse(fmt.Sprintf("%s/%s", a.authServerURL, a.path))
	if err != nil {
		log.Println("Error parsing URL:", err)
		return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
	}

	// Add the form data to the existing query parameters
	queryParams := u.Query()

	// Add the app parmeters
	// appParams, err := url.ParseQuery(appUrlParams)
	appParams, err := url.ParseQuery(a.authCredentials.AppUrlParams)
	if err != nil {
		log.Println("Error parsing query string:", err)
		return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
	}

	for key, values := range appParams {
		for _, value := range values {
			queryParams.Add(key, value)
		}
	}

	u.RawQuery = queryParams.Encode()

	// log.Println("see where the first req is going: ", u.String())

	resp, err := http.Post(
		u.String(),
		"application/x-www-form-urlencoded",
		strings.NewReader(""),
	)
	if err != nil {
		log.Println("Error request the token: ", err)
		return e.NewCustomHTTPStatus(e.StatusBadRequest, "", err.Error())
	}

	defer resp.Body.Close()

	// Handle the response as needed
	if resp.StatusCode != http.StatusOK {
		log.Println("Error StatusCode request the token: ", resp)
		return e.NewCustomHTTPStatus(e.StatusCode(resp.StatusCode))
		// return nil
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
	}

	// Parse the response body into CodeStruct
	a.codeStruct = &CodeStruct{}
	err = json.Unmarshal(body, a.codeStruct)
	if err != nil {
		log.Println("Error parsing response body:", err)
		return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
	}

	// Define the request parameters
	// method := "POST"
	bodyString := url.Values{}
	bodyString.Set("code", a.codeStruct.Code)
	bodyString.Set("code_verifier", a.codeVerifier)
	bodyString.Set("grant_type", a.grantType)
	bodyString.Set("redirect_uri", a.redirectURI)
	bodyString.Set("sub", a.clientID)
	bodyString.Set("role", a.role)
	// bodyString.Set("token_expiration", "60") // will overwrite the default which is 1 month

	bodyByt := strings.NewReader(bodyString.Encode())

	basicEncodedBase64 := base64.StdEncoding.EncodeToString([]byte(a.clientID + ":" + a.clientSecret))
	basicAuth := "Basic " + basicEncodedBase64

	// Create a new HTTP request
	req, err := http.NewRequest(a.method, a.urlOauthToken, bodyByt)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
	}

	// Set request headers
	req.Header.Set("Accept-Encoding", a.acceptEncoding)
	req.Header.Set("Authorization", basicAuth)
	req.Header.Set("Content-Type", a.contentType)

	// Execute the request
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return e.NewCustomHTTPStatus(e.StatusBadRequest, "", err.Error())
	}
	defer resp.Body.Close()

	// Read the response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
	}

	// fmt.Println("order.go see the bodyyyyyy: ", string(body))

	a.authCredentials.JWT_token = string(body)

	return nil
}

// func setURLAndCreateCodeChallenge(codeVerifier, redirectURI, endpointAuthURI, clientID, codeChallengeMethod, scope, state string) string {
func (a APIserverAuth) SetURLAndCreateCodeChallenge() e.IError {
	// Generate the code challenge
	s256 := a.genCodeChallengeS256(a.codeVerifier)

	// Encode redirect URL
	encodedURI := url.QueryEscape(a.redirectURI)
	fmt.Println("See encodedURI:", encodedURI)

	// Encode scope
	encodedScope := url.QueryEscape(a.scope)
	encodedScope = strings.Replace(encodedScope, "%20", "%2C+", -1)

	// Create the URL,
	// note: authServerURL does not matter, the URL only have to be a full URL
	settedURL := fmt.Sprintf("%s?client_id=%s&code_challenge=%s&code_challenge_method=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s&role=%s",
		a.authServerURL, a.clientID, s256, a.codeChallengeMethod, encodedURI, encodedScope, a.state, a.role)

	parsedURL, err := url.Parse(settedURL)
	if err != nil {
		return e.NewCustomHTTPStatus(e.StatusInternalServerError, "", err.Error())
	}

	a.authCredentials.AppUrlParams = parsedURL.RawQuery
	return nil
}

func (a APIserverAuth) genCodeChallengeS256(s string) string {
	s256 := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(s256[:])
}

func (a APIserverAuth) GetToken() string {
	return a.authCredentials.JWT_token
}

package auth

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func SetGoogleOAuth() {
	authHandler := NewGoogleOAuthHandler()

	http.HandleFunc("/auth/google/login", authHandler.Login)

	http.HandleFunc("/auth/google/callback", authHandler.Callback)
}

type OAuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	Callback(w http.ResponseWriter, r *http.Request)
}

type googleOAuthHandler struct {
	googleOauthConfig *oauth2.Config
	jwtHandler        *jwtHandler
}

func NewGoogleOAuthHandler() OAuthHandler {
	clientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	if clientID == "" {
		log.Fatal("You must set the GOOGLE_OAUTH_CLIENT_ID environment variable")
	}

	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	if clientSecret == "" {
		log.Fatal("You must set the GOOGLE_OAUTH_CLIENT_SECRET environment variable")
	}

	var googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	jwtHandler := NewJWTHandler()
	return &googleOAuthHandler{googleOauthConfig: googleOauthConfig, jwtHandler: jwtHandler}
}

func (g *googleOAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	state := "random" // 應該產生一個隨機的狀態值用於防止 CSRF 攻擊
	url := g.googleOauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (g *googleOAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	code := r.URL.Query().Get("code")
	if code == "" {
		log.Println("Code not found")
		httpError(w, "Authorization code not found", http.StatusBadRequest) // 400 Bad Request
		return
	}

	token, err := g.googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange token: %s\n", err)
		httpError(w, "Failed to exchange token", http.StatusInternalServerError) // 500 Internal Server Error
		return
	}

	client := g.googleOauthConfig.Client(ctx, token)
	userinfo, err := getUserInfo(client)
	if err != nil {
		log.Printf("Failed to get user info: %s\n", err)
		httpError(w, "Failed to retrieve user information", http.StatusInternalServerError) // 500 Internal Server Error
		return
	}

	if !userinfo.EmailVerified {
		log.Println("OAuth return unverified email")
		httpError(w, "User email not verified", http.StatusBadRequest) // 400 Bad Request
		return
	}

	jwtToken, err := g.jwtHandler.GenerateJWT(userinfo.Email)
	if err != nil {
		log.Printf("Failed to generate JWT: %s\n", err)
		httpError(w, "Failed to generate JWT", http.StatusInternalServerError) // 500 Internal Server Error
		return
	}

	response := map[string]string{
		"token": jwtToken,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// httpError simplifies sending error messages
func httpError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func getUserInfo(client *http.Client) (*User, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var user User

	err = json.Unmarshal(data, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

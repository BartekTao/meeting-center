package domain

import "context"

type User struct {
	ID            *string `json:"_id,omitempty"`
	Sub           string  `json:"sub"`
	Name          string  `json:"name"`
	GivenName     string  `json:"given_name"`
	FamilyName    string  `json:"family_name"`
	Picture       string  `json:"picture"`
	Locale        string  `json:"locale"`
	Email         string  `json:"email"`
	EmailVerified bool    `json:"email_verified"`
}

type UserRepo interface {
	GetUser(ctx context.Context, sub string) (*User, error)
	SignUp(ctx context.Context, user User) (*string, error)
}

package handlers

import (
	"html/template"

	"github.com/shabinesh/transcription/core/user"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var emailTemplates = template.Must(
	template.New("email").ParseFiles("templates/email/account_verification_email.tpl"),
)

type UserService interface {
	Register(email string, info map[string]string) (*user.User, error)
	VerifyCode(userID string, otp string) error
	SendOTP(email string) error
	Authenticate(email, code string) (*user.User, error)
	GetUser(email string) (*user.User, error)
}

type handlers struct {
	userService  UserService
	sessionStore sessions.Store
}

func NewHandlers(userService UserService) *handlers {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)
	return &handlers{userService: userService, sessionStore: sessions.NewCookieStore(authKeyOne, encryptionKeyOne)}
}

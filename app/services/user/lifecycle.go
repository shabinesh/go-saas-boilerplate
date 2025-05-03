package user

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/shabinesh/app/core/user"
)

var (
	EmailVerifcationSubject = "Welcome! Please verify your email"
	EmailVerificationBody   = `Enter code %s to verify your email address.`

	ErrorUserExists    = errors.New("user already exists")
	ErrInvalidOTP      = errors.New("invalid otp")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserNotActive   = errors.New("user not active")
	ErrOTPVerification = errors.New("failed to verify otp")
	ErrJWTGeneration   = errors.New("failed to generate jwt")
)

type Storer interface {
	FindUser(id string) (*user.User, error)
	FindUserByEmail(email string) (*user.User, error)
	AddUser(*user.User) (*user.User, error)
	UpdateUserStatus(uu *user.User) error
}

type Mailer interface {
	SendEmail(to, subject, body string) error
}

type OTPProcessor interface {
	Generate(userID string, reason string) string
	Verify(userID string, otp string, reason string) (bool, error)
}

type UserService struct {
	userStore    Storer
	mailer       Mailer
	otpProcessor OTPProcessor

	secret []byte
}

func NewUserService(userStore Storer, otpProcessor OTPProcessor, mailer Mailer, secret []byte) *UserService {
	return &UserService{
		userStore:    userStore,
		otpProcessor: otpProcessor,
		mailer:       mailer,
		secret:       secret,
	}
}

func (a *UserService) SendOTP(id user.UserID, reason string) error {
	u, err := a.userStore.FindUser(string(id))
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	otp := a.otpProcessor.Generate(string(u.ID), reason)

	return a.mailer.SendEmail(u.Email, EmailVerifcationSubject, fmt.Sprintf(EmailVerificationBody, otp))
}

func (a *UserService) SendLoginOTP(email string) error {
	user, err := a.userStore.FindUserByEmail(email)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}

	return a.SendOTP(user.ID, "login")
}

func (a *UserService) Authenticate(res gin.ResponseWriter, email string, otp string, otpReason string) (*user.User, error) {
	user, err := a.userStore.FindUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	ok, err := a.otpProcessor.Verify(string(user.ID), otp, otpReason)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOTPVerification, err)
	}

	if !ok {
		return nil, ErrInvalidOTP
	}

	token, err := genereteJWT(user, a.secret, map[string]interface{}{
		"email": email,
	})

	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrJWTGeneration, err)
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(res, cookie)

	return user, nil
}

func (a *UserService) Register(email string, info map[string]string) (*user.User, error) {
	_, err := a.userStore.FindUserByEmail(email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	u, err := a.userStore.AddUser(&user.User{Email: email})
	if err != nil {
		return nil, err
	}

	err = a.SendOTP(u.ID, "register")
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (a *UserService) VerifyCode(userID, otp, reason string) error {
	ok, err := a.otpProcessor.Verify(userID, otp, reason)
	if err != nil {
		return err
	}

	if !ok {
		return ErrInvalidOTP
	}

	user, err := a.userStore.FindUser(userID)
	if err != nil {
		return err
	}

	user.IsVerified = true
	user.IsActive = true

	return a.userStore.UpdateUserStatus(user)
}

func (a *UserService) GetUser(email string) (*user.User, error) {
	user, err := a.userStore.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func genereteJWT(user *user.User, secret []byte, claims map[string]interface{}) (string, error) {

	mapclaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	for k, v := range claims {
		mapclaims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapclaims)

	return token.SignedString(secret)
}

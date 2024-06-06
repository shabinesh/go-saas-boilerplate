package user

import (
	"errors"
	"fmt"
)

var (
	EmailVerifcationSubject = "Welcome! Please verify your email"
	EmailVerificationBody   = `Enter code %s to verify your email address.`

	ErrorUserExists = errors.New("user already exists")
)

type UserStore interface {
	FindUser(id string) (*User, bool, error)
	FindUserByEmail(email string) (*User, bool, error)
	AddUser(*User) (*User, error)
	UpdateUserStatus(uu *User) error
}

type Mailer interface {
	SendEmail(to, subject, body string) error
}

type OTPProcessor interface {
	Generate(userID string) string
	Verify(userID string, otp string) (bool, error)
}

type UserService struct {
	userStore    UserStore
	mailer       Mailer
	otpProcessor OTPProcessor
}

func NewUserService(userStore UserStore, otpProcessor OTPProcessor, mailer Mailer) *UserService {
	return &UserService{userStore: userStore, otpProcessor: otpProcessor, mailer: mailer}
}

func (a *UserService) Login(email, code string) (*User, error) {
	otp := a.otpProcessor.Generate(string(u.ID))

	err := a.mailer.SendEmail(email, EmailVerifcationSubject, fmt.Sprintf(EmailVerificationBody, otp))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (a *UserService) Authenticate(email string, otp string) (*User, error) {
	user, found, err := a.userStore.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("user not active")
	}

	if ok, err := a.otpProcessor.Verify(string(user.ID), otp); err != nil || !ok {
		return nil, errors.New("invalid otp")
	}

	return user, nil
}

func (a *UserService) Register(email string, info map[string]string) (*User, error) {
	_, found, err := a.userStore.FindUserByEmail(email)
	if !found {
		u, err := a.userStore.AddUser(&User{Email: email})
		if err != nil {
			return nil, err
		}

		otp := a.otpProcessor.Generate(string(u.ID))

		err = a.mailer.SendEmail(email, EmailVerifcationSubject, fmt.Sprintf(EmailVerificationBody, otp))
		if err != nil {
			return nil, err
		}

		return u, nil
	} else if err != nil {
		return &User{}, err
	}

	return nil, fmt.Errorf("User already exists with this email")
}

func (a *UserService) VerifyCode(userID, otp string) error {
	ok, err := a.otpProcessor.Verify(userID, otp)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("invalid otp")
	}

	user, found, err := a.userStore.FindUser(userID)
	if err != nil {
		return err
	}

	if !found {
		return errors.New("user not found")
	}

	user.IsVerified = true
	user.IsActive = true

	return a.userStore.UpdateUserStatus(user)
}

func (a *UserService) GetUser(email string) (*User, error) {
	user, found, err := a.userStore.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !found {
		return nil, errors.New("user not found")
	}

	return user, nil
}

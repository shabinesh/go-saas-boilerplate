package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	registrationVerifiedMessage = "Your email has been verified successfully. Please login to continue."
	sessionKey                  = "session"
)

func (h handlers) RegisterPage(r *gin.Context) {
	r.HTML(http.StatusOK, "register", gin.H{})
}

func (h handlers) Register(r *gin.Context) {
	//sess, _ := h.sessionStore.Get(r.Request, sessionKey)
	email := r.PostForm("email")
	fullName := r.PostForm("fullname")

	user, err := h.userService.Register(email, map[string]string{
		"full_name": fullName,
	})

	if err != nil {
		slog.Error(err.Error())
		r.Writer.WriteHeader(http.StatusInternalServerError)
		r.Writer.Write([]byte(err.Error()))

		return
	}

	r.HTML(http.StatusFound, "get_otp", gin.H{"user_id": user.ID})
}

func (h handlers) Authenticate(r *gin.Context) {
	sess, _ := h.sessionStore.Get(r.Request, sessionKey)

	email := r.PostForm("email")
	code := r.PostForm("code")

	if email == "" || len(code) < 6 {
		r.HTML(http.StatusBadRequest, "login_get_code", gin.H{
			"email": email,
			"error": "email and OTP are required",
		})

		return
	}

	if _, err := h.userService.Authenticate(email, code); err != nil {
		slog.Error(err.Error())
		r.HTML(http.StatusForbidden, "login_get_code", gin.H{
			"email": email,
			"error": err.Error(),
		})

		return
	}

	sess.Values["authenticated"] = true
	err := sess.Save(r.Request, r.Writer)
	if err != nil {
		slog.Error(err.Error())

		return
	}

	r.Redirect(http.StatusFound, "/app/home")
}

func (h handlers) GetOTP(r *gin.Context) {
	userID := r.PostForm("user_id")

	if userID == "" {
		r.Writer.WriteHeader(http.StatusInternalServerError)
		r.Writer.Write([]byte("User ID not found in session"))

		return
	}

	codes := r.PostFormArray("code[]")

	if err := h.userService.VerifyCode(userID, strings.Join(codes, "")); err != nil {
		slog.Error(err.Error())
		r.Writer.WriteHeader(http.StatusBadRequest)
		r.HTML(http.StatusOK, "get_otp", gin.H{
			"Error": err.Error(),
		})

		r.HTML(http.StatusOK, "get_otp", gin.H{"user_id": userID, "message": "failed to verify OTP"})
	}

	r.HTML(http.StatusOK, "message", gin.H{
		"message": registrationVerifiedMessage,
	})
}

func (h handlers) LoginPage(r *gin.Context) {
	if r.Request.Method == http.MethodGet {
		r.HTML(http.StatusOK, "login_get_email", gin.H{})
		return
	}

	email := r.PostForm("email")

	if email == "" {
		r.HTML(http.StatusBadRequest, "login_get_email", gin.H{"error": "email is required"})

		return
	}

	user, err := h.userService.GetUser(email)
	if err != nil {
		slog.Error(err.Error())
		r.HTML(http.StatusInternalServerError, "message", gin.H{"message": "An error occurred while authenticating you."})

		return
	}

	err = h.userService.SendOTP(user.ID)
	if err != nil {
		slog.Error("Error sending OTP", err.Error(), nil)
		r.HTML(http.StatusOK, "login_get_email", gin.H{"error": err.Error()})

		return
	}

	r.HTML(http.StatusOK, "login_get_code", gin.H{"email": email})
}

func (h handlers) Logout(r *gin.Context) {
	sess, _ := h.sessionStore.Get(r.Request, sessionKey)
	delete(sess.Values, "authenticated")
	sess.Options.MaxAge = -1
	sess.Save(r.Request, r.Writer)

	r.Redirect(http.StatusFound, "/login")
}

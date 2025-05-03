package handlers

import (
	"log/slog"
	"net/http"

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

	r.HTML(http.StatusFound, "get_otp", gin.H{"user_id": user.ID, "reason": "register"})
}

func (h handlers) LoginPage(r *gin.Context) {
	r.HTML(http.StatusOK, "login_get_email", gin.H{})
	return
}

func (h handlers) Login(r *gin.Context) {
	email := r.PostForm("email")

	err := h.userService.SendLoginOTP(email)
	if err != nil {
		slog.Error(err.Error())
		r.HTML(http.StatusOK, "login_get_email", gin.H{"email": email, "error": err.Error()})
		return
	}

	r.HTML(http.StatusFound, "login_get_code", gin.H{"email": email})
}

func (h handlers) VerifyCode(r *gin.Context) {
	userID := r.PostForm("user_id")
	code := r.PostForm("code")
	reason := r.PostForm("reason")

	err := h.userService.VerifyCode(userID, code, reason)
	if err != nil {
		slog.Error(err.Error())

		r.HTML(http.StatusOK, "get_otp", gin.H{
			"user_id": userID,
			"reason":  reason,
			"message": err.Error(),
		})

		return
	}

	r.Header("HX-Redirect", "/login")
	r.Status(http.StatusOK)
}

func (h handlers) Authenticate(r *gin.Context) {
	email := r.PostForm("email")
	code := r.PostForm("code")

	_, err := h.userService.Authenticate(r.Writer, email, code, "login")
	if err != nil {
		slog.Error(err.Error())
		r.HTML(http.StatusOK, "login_get_code", gin.H{"email": email, "error": err.Error()})
		return
	}

	r.Header("HX-Redirect", "/app/home")
	r.Status(http.StatusOK)
}

func (h handlers) Logout(r *gin.Context) {
	r.SetCookie("token", "", -1, "/", "", false, true)
	r.Header("HX-Redirect", "/login")

	r.Redirect(http.StatusFound, "/login")
}

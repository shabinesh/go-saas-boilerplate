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
	sess, _ := h.sessionStore.Get(r.Request, sessionKey)
	email := r.PostForm("email")
	fullName := r.PostForm("fullname")

	u, err := h.userService.Register(email, map[string]string{
		"full_name": fullName,
	})

	if err != nil {
		slog.Error(err.Error())
		r.Writer.WriteHeader(http.StatusInternalServerError)
		r.Writer.Write([]byte(err.Error()))

		return
	}

	sess.Values["intent"] = "verify_email"
	sess.Values["user_id"] = string(u.ID)
	err = sess.Save(r.Request, r.Writer)
	if err != nil {
		slog.Error(err.Error())
		r.Writer.WriteHeader(http.StatusInternalServerError)
		r.Writer.Write([]byte(err.Error()))

		return
	}

	r.HTML(http.StatusOK, "get_otp", gin.H{})
}

func (h handlers) GetOTP(r *gin.Context) {
	sess, _ := h.sessionStore.Get(r.Request, sessionKey)

	intent := sess.Values["intent"]

	if intent == "verify_email" {
		userID := sess.Values["user_id"]

		if userID == nil {
			r.Writer.WriteHeader(http.StatusInternalServerError)
			r.Writer.Write([]byte("User ID not found in session"))

			return
		}

		code := r.PostForm("code")

		if err := h.userService.VerifyCode(userID.(string), code); err != nil {
			slog.Error(err.Error())
			r.Writer.WriteHeader(http.StatusBadRequest)
			r.HTML(http.StatusOK, "get_otp", gin.H{
				"Error": err.Error(),
			})

			return
		}
	}

	if intent == "login" {
		email := sess.Values["email"]
		code := r.PostForm("code")

		if _, err := h.userService.Login(email.(string), code); err != nil {
			slog.Error(err.Error())
			r.Writer.WriteHeader(http.StatusBadRequest)
			r.HTML(http.StatusOK, "get_otp", gin.H{
				"Error": err.Error(),
			})

			return
		}

		sess.Values["authenticated"] = true
		err := sess.Save(r.Request, r.Writer)
		if err != nil {
			slog.Error(err.Error())
			r.Writer.WriteHeader(http.StatusInternalServerError)
			r.Writer.Write([]byte(err.Error()))

			return
		}
	}

	r.HTML(http.StatusOK, "message", gin.H{
		"Message": registrationVerifiedMessage,
	})
}

func (h handlers) Login(r *gin.Context) {
	if r.Request.Method == http.MethodGet {
		r.HTML(http.StatusOK, "login", gin.H{})
		return
	}

	email := r.PostForm("email")

	sess, _ := h.sessionStore.Get(r.Request, sessionKey)
	sess.Values["email"] = email
	sess.Values["intent"] = "login"
	err := sess.Save(r.Request, r.Writer)
	if err != nil {
		slog.Error(err.Error())
		r.Writer.WriteHeader(http.StatusInternalServerError)
		r.Writer.Write([]byte(err.Error()))

		return
	}

	r.HTML(http.StatusOK, "get_otp", gin.H{})
}

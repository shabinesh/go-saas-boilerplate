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
	//sess, _ := h.sessionStore.Get(r.Request, sessionKey)
	email := r.PostForm("email")
	fullName := r.PostForm("fullname")

	_, err := h.userService.Register(email, map[string]string{
		"full_name": fullName,
	})

	if err != nil {
		slog.Error(err.Error())
		r.Writer.WriteHeader(http.StatusInternalServerError)
		r.Writer.Write([]byte(err.Error()))

		return
	}

	r.HTML(http.StatusFound, "get_otp", gin.H{})
}

func (h handlers) Authenticate(r *gin.Context) {
	sess, _ := h.sessionStore.Get(r.Request, sessionKey)

	email := r.PostForm("email")
	code := r.PostForm("code")

	if email == "" || code == "" {
		r.Writer.WriteHeader(http.StatusBadRequest)
		r.HTML(http.StatusBadRequest, "login_get_code", gin.H{
			"email": email,
			"error": "email and code are required",
		})
	}

	if _, err := h.userService.Authenticate(email, code); err != nil {
		slog.Error(err.Error())
		r.Writer.WriteHeader(http.StatusBadRequest)
		r.HTML(http.StatusForbidden, "login_get_code", gin.H{
			"email": email,
			"error": err.Error(),
		})

		return
	} else {
		sess.Values["authenticated"] = true
		err := sess.Save(r.Request, r.Writer)
		if err != nil {
			slog.Error(err.Error())

			return
		}

		r.Redirect(http.StatusFound, "/app/home")
	}
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

	//err := h.userService.SendOTP(email)
	//if err != nil {
	//	slog.Error("Error sending OTP", err.Error(), nil)
	//	r.HTML(http.StatusOK, "login_get_email", gin.H{"error": err.Error()})
	//
	//	return
	//}

	r.HTML(http.StatusOK, "login_get_code", gin.H{"email": email})
}

func (h handlers) Logout(r *gin.Context) {
	sess, _ := h.sessionStore.Get(r.Request, sessionKey)
	delete(sess.Values, "authenticated")
	sess.Options.MaxAge = -1
	sess.Save(r.Request, r.Writer)

	r.Redirect(http.StatusFound, "/login")
}

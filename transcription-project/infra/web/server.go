package web

import (
	"database/sql"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"

	"github.com/shabinesh/transcription/core/user"
	"github.com/shabinesh/transcription/infra/repo"
	"github.com/shabinesh/transcription/infra/web/handlers"
	"github.com/shabinesh/transcription/services/mailer"
	"github.com/shabinesh/transcription/services/otp"
)

func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("register", "templates/base.html", "templates/auth/register.html")
	r.AddFromFiles("get_otp", "templates/auth/get_otp.html")
	r.AddFromFiles("message", "templates/message.html").Delims("{%", "%}")
	r.AddFromFiles("login_get_email", "templates/base.html", "templates/auth/login_get_email.html")
	r.AddFromFiles("login_get_code", "templates/base.html", "templates/auth/login_get_code.html")
	r.AddFromFiles("home", "templates/base.html", "templates/app/home.html")

	return r
}

func StartServer(db *sql.DB) {
	server := gin.Default()
	server.HTMLRender = createRender()

	userRepo := repo.NewUserRepo(db)
	otpRepo := repo.NewOTPRepository(db)
	otpService := otp.NewOTPService(otpRepo)
	emailer := mailer.NewMockMailer()
	userService := user.NewUserService(userRepo, otpService, emailer)
	handlers := handlers.NewHandlers(userService)

	server.GET("/register", handlers.RegisterPage)
	server.POST("/register", handlers.Register)
	server.POST("/verify-otp", handlers.GetOTP)
	server.GET("/login", handlers.LoginPage)
	server.POST("/login", handlers.LoginPage)
	server.POST("/authenticate", handlers.Authenticate)
	server.GET("/logout", handlers.Logout)

	// protected routes
	protected := server.Group("/app")
	protected.Use(handlers.RequireAuth())
	{
		protected.GET("/home", handlers.HomePage)
	}

	server.Run()
}

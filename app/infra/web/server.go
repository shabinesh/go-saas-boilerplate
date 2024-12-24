package web

import (
	"github.com/jackc/pgx/v5"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"

	"github.com/shabinesh/app/core/user"
	"github.com/shabinesh/app/infra/repo"
	"github.com/shabinesh/app/infra/web/handlers"
	"github.com/shabinesh/app/services/mailer"
	"github.com/shabinesh/app/services/otp"
)

func createRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("register", "templates/base.html", "templates/auth/register.html")
	r.AddFromFiles("get_otp", "templates/auth/get_otp.html")
	r.AddFromFiles("message", "templates/message.html").Delims("<%", "%>")
	r.AddFromFiles("login_get_email", "templates/base.html", "templates/auth/login_get_email.html")
	r.AddFromFiles("login_get_code", "templates/base.html", "templates/auth/login_get_code.html")
	r.AddFromFiles("home", "templates/base.html", "templates/app/home.html")

	return r
}

func StartServer(db *pgx.Conn) {
	server := gin.Default()
	server.HTMLRender = createRender()

	userRepo := repo.NewUserRepo(db)
	otpRepo := repo.NewOTPRepository(db)
	otpService := otp.NewOTPService(otpRepo)
	emailer := mailer.NewMockMailer()
	userService := user.NewUserService(userRepo, otpService, emailer)
	apiHandler := handlers.NewHandlers(userService)

	server.POST("/register", apiHandler.Register)
	server.GET("/register", apiHandler.RegisterPage)
	server.POST("/verify-otp", apiHandler.GetOTP)
	server.GET("/login", apiHandler.LoginPage)
	server.POST("/login", apiHandler.LoginPage)
	server.POST("/authenticate", apiHandler.Authenticate)
	server.GET("/logout", apiHandler.Logout)

	// protected routes
	protected := server.Group("/app")
	protected.Use(apiHandler.RequireAuth())
	{
		protected.GET("/home", apiHandler.HomePage)
	}

	server.Run(":9090")
}

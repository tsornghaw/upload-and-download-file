package routes

import (
	"upload-and-download-file/controllers"
	"upload-and-download-file/database"
	"upload-and-download-file/models"

	"github.com/gofiber/fiber"
)

type Server struct {
	//logger                 zerolog.Logger
	config models.Config
	gd     *database.GormDatabase
}

func NewServer(config models.Config) *Server {
	database.Connect()

	return &Server{
		config: config,
		gd:     database.GDB,
	}
}

// func Setup(app *fiber.App) {}
func (s *Server) Run() {
	app := fiber.New()

	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)

	app.Listen(":8000")
}

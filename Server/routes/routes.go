package routes

import (
	"net/http"
	"reflect"
	"upload-and-download-file/database"
	"upload-and-download-file/models"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config models.Config
	gd     *database.GormDatabase
}

func NewServer(config models.Config) *Server {
	gd := database.Connect()

	return &Server{
		config: config,
		gd:     gd,
	}
}

func (s *Server) Run() {
	// Create a new gin router with default middleware
	r := gin.Default()

	// CORS configuration
	r.Use(s.CORS())

	// // Logging middleware
	// r.Use(func(c *gin.Context) {
	// 	log.Printf("Received request from %s %s\n", c.Request.Method, c.FullPath())
	// 	c.Next()
	// })

	// User registration and login APIs (no authentication)
	api := r.Group("/api")
	{
		api.POST("/register", s.Register)
		api.POST("/login", s.Login)
		//api.GET("/user", s.User)

		// File upload and download (Authenticated APIs)
		authr := api.Group("/auth")
		{
			authr.GET("/user", s.User)
			authr.POST("/logout", s.Logout)
			authr.POST("/upload", s.Upload)
			authr.GET("/download/:URL", s.CheckDownloadCount, s.Download, s.UpdateDownloadCount)
			authr.GET("/UserSearchAllData", s.UserSearchAllData)

			// adminstrator
			adminr := authr.Group("/admin", s.AdminMiddleware())
			{
				adminr.GET("/SearchAllUsers", s.SearchAllUsers)
				adminr.GET("/SearchAllData", s.SearchAllData)
				adminr.DELETE("/deleteuser", s.DeleteUsers)
				adminr.DELETE("/deletedata", s.DeleteDatas)
			}
		}
	}

	// Run a background goroutine to periodically delete expired files
	// go s.DeleteExpiredFiles()

	r.Run(":8000")
}

func (s *Server) DeleteUsers(c *gin.Context) {
	var frontend_ids models.FrontendRequest
	var users []models.User

	if err := c.ShouldBindJSON(&frontend_ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if len(frontend_ids.FrontendIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided for deletion"})
		return
	}

	for _, fronted_id := range frontend_ids.FrontendIDs {

		if err := s.gd.GetCorresponding(&users, "id = ?", fronted_id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No id existed in database"})
			return
		}

		if reflect.ValueOf(users).IsZero() {
			c.JSON(http.StatusNotFound, gin.H{"message": "User data does not exist!"})
			return
		}

		// Delete users by their IDs
		if err := s.gd.Delete(&users); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Users deleted successfully"})

	}
}

package routes

import (
	"upload-and-download-file/models"

	"github.com/gin-gonic/gin"
)

func (s *Server) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		cookie, err := c.Cookie("jwt")

		if cookie == "" || err != nil {
			c.JSON(401, gin.H{"message": "Unauthenticated 1"})
			c.Abort()
			return
		}

		claims, err := ValidateToken(cookie)

		if err != nil {
			c.JSON(401, gin.H{"message": "Unauthenticated 2"})
			c.Abort()
			return
		}

		// Retrieve user from your database using claims.Issuer
		var user models.User

		if err := s.gd.GetCorresponding(&user, "id = ?", claims.Issuer); err != nil {
			c.JSON(404, gin.H{"message": "User not found"})
			c.Abort()
			return
		}

		c.Set("User", user.Name)
		c.Set("Admin", user.Admin)
		c.Next()
	}
}

func (s *Server) AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if admin, _ := c.Get("Admin"); !admin.(bool) {
			c.JSON(403, gin.H{"message": "you have no permission"})
			c.Abort()
			return
		}

		c.Next()
	}
}

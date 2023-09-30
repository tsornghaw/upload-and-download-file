package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"upload-and-download-file/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const SecretKey = "secret"
const MaxDownloads = 5

func (s *Server) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content for preflight requests
			return
		}
		c.Next()
	}
}

func (s *Server) Register(c *gin.Context) {
	log.Printf("Start Register\n")
	var data map[string]string

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"message": "Bad Request"})
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	if err != nil {
		c.JSON(500, gin.H{"message": "Internal Server Error"})
		return
	}

	// Create and save user in your database
	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	if err := s.gd.Create(&user); err != nil {
		c.JSON(400, gin.H{"message": "Fail to create account"})
		return
	}

	c.JSON(200, user)
}

func (s *Server) Login(c *gin.Context) {
	log.Printf("Start Login\n")
	var data map[string]string

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"message": "Bad Request"})
		return
	}

	PrettyPrint(data)

	// Retrieve user from your database
	var user models.User

	if err := s.gd.GetCorresponding(&user, "email = ?", data["email"]); err != nil {
		c.JSON(400, gin.H{"message": err})
		return
	}

	PrettyPrint(user)

	if user.Id == 0 {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.JSON(400, gin.H{"message": "Incorrect password"})
		return
	}

	token, err := GenerateToken(int(user.Id))
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal Server Error"})
		return
	}

	// TODO: Use Header Authrization to create JWT
	// https://ithelp.ithome.com.tw/articles/10278153
	c.SetCookie("jwt", token, 3600, "/", "http:/localhost:3000/", false, false)
	//c.JSON(200, gin.H{"message": "Success"})
	c.JSON(200, gin.H{"name": user.Name})
}

func (s *Server) User(c *gin.Context) {
	log.Printf("Start User\n")

	cookie, err := c.Cookie("jwt")
	PrettyPrint(cookie)
	if cookie == "" || err != nil {
		PrettyPrint("Unauthenticated1")
		c.JSON(401, gin.H{"message": "Unauthenticated"})
		return
	}

	claims, err := ValidateToken(cookie)
	PrettyPrint("claims : ")
	PrettyPrint(claims)
	if err != nil {
		PrettyPrint("Unauthenticated2")
		c.JSON(401, gin.H{"message": "Unauthenticated"})
		return
	}

	// Retrieve user from your database using claims.Issuer
	var user models.User

	s.gd.GetCorresponding(&user, "id = ?", claims.Issuer)

	c.JSON(200, user)
}

func (s *Server) Logout(c *gin.Context) {
	log.Printf("Start Logout\n")
	c.SetCookie("jwt", "", -1, "/", "http:/localhost:3000/", false, false)
	c.JSON(200, gin.H{"message": "Success"})
}

func (s *Server) SearchAllUsers(c *gin.Context) {
	var users []models.User

	if err := s.gd.GetCorresponding(&users, "admin <> ?", "true"); err != nil {
		c.JSON(400, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GenerateToken(userID int) (string, error) {
	claims := jwt.StandardClaims{
		Issuer:    strconv.Itoa(userID),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(SecretKey))
}

func ValidateToken(cookie string) (*jwt.StandardClaims, error) {

	// Parse the JWT token and validate it
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Provide the secret key for token validation
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		// Token is valid, and you can access the claims
		return claims, nil
	}

	return nil, err
}

// print the contents of the obj
func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}

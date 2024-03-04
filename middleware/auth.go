package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("")
var JwtUsername = ""
var Tokens []string

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GetToken(c *gin.Context) {
	token, err := GenerateJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	Tokens = append(Tokens, token)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		log.Println(bearerToken)
		if bearerToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}

		parts := strings.Split(bearerToken, " ")
		if len(parts) != 2 {
			// Handle error, maybe set an appropriate response
			return
		}
		reqToken := parts[1]
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			log.Println(err)
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			c.Abort()
			return
		}
		if !tkn.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			c.Abort()
			return
		}
	}
}

func GenerateJWT() (string, error) {
	// get jwt key secret from os env
	JwtKey = []byte(os.Getenv("JWT_SECRET"))
	if JwtKey == nil {
		return "", errors.New("JWT_SECRET not found in .env file")
	}

	JwtUsername = os.Getenv("JWT_USERNAME")
	if JwtUsername == "" {
		return "", errors.New("JWT_USERNAME not found in .env file")
	}

	expirationTime := time.Now().Add(30 * time.Second)
	claims := &Claims{
		Username: JwtUsername,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(JwtKey)

}

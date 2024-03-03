package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AllowedOriginsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		referer := c.Request.Header.Get("Referer")
		clientIP := c.ClientIP()

		log.Println("Request Origin:", origin)
		log.Println("Request Referer:", referer)
		log.Println("Client IP:", clientIP)

		var allowed bool

		if origin != "" && matchAllowedOrigin(origin, allowedOrigins) {
			log.Println("Origin allowed:", origin)
			allowed = true
		} else if referer != "" && matchAllowedOrigin(referer, allowedOrigins) {
			log.Println("Referer allowed:", referer)
			allowed = true
		} else if isLocalhost(clientIP) {
			log.Println("Localhost allowed:", clientIP)
			allowed = true
		}

		if allowed {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Origin not allowed"})
		}
	}
}

func matchAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

func isLocalhost(ip string) bool {
	return ip == "127.0.0.1" || ip == "::1"
}

package middleware

import (
	"github.com/gin-gonic/gin"
	s "strings"
)

func WSAuthRewrite() gin.HandlerFunc {
	return func(c *gin.Context) {
		protocolHeaders := s.Split(c.Request.Header.Get("Sec-WebSocket-Protocol"), ",")
		for _, header := range protocolHeaders {
			header := s.TrimSpace(header)
			if s.HasPrefix(header, "Bearer") {
				c.Request.Header.Set("Authorization", header)
				break
			}
		}
	}
}

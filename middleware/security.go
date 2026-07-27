package middleware

import "github.com/gin-gonic/gin"

func SecureCookies() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Add("X-Frame-Options", "DENY")
		c.Writer.Header().Add("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

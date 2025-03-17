package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"brickwall/internal/provider"
)

func AuthMiddleware(jwt provider.IJwtProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		// Заглушка: проверка статуса пользователя и его ролей из БД (можно заменить реальным вызовом)
		userStatus := "active"        // допустим, получили из БД
		userRoles := []string{"user"} // допустим, получили из БД

		if userStatus != "active" {
			c.JSON(http.StatusForbidden, gin.H{"error": "user is blocked"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("roles", userRoles)
		c.Next()
	}
}

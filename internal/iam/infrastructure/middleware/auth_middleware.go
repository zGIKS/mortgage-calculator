package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"finanzas-backend/internal/iam/infrastructure/security"
)

// JWTAuthMiddleware verifica el token JWT en las peticiones
func JWTAuthMiddleware(jwtService *security.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Verificar formato "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Use: Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar el token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Guardar información del usuario en el contexto
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)

		c.Next()
	}
}

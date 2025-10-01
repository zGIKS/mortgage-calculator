package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"finanzas-backend/internal/mortgage/application/acl"
)

// JWTAuthMiddleware verifica el token JWT usando el servicio externo de autenticación (ACL)
func JWTAuthMiddleware(externalAuthService *acl.ExternalAuthenticationService) gin.HandlerFunc {
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

		// Validar el token a través del servicio externo (ACL)
		userID, err := externalAuthService.ValidateTokenAndGetUserID(c.Request.Context(), tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Guardar información del usuario en el contexto
		c.Set("user_id", userID.Value())

		c.Next()
	}
}

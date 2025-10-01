package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService maneja la generación y validación de tokens JWT
type JWTService struct {
	secretKey     string
	issuer        string
	expirationHrs int
}

// JWTClaims representa los claims personalizados del JWT
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// NewJWTService crea una nueva instancia del servicio JWT
func NewJWTService(secretKey string, issuer string, expirationHrs int) *JWTService {
	return &JWTService{
		secretKey:     secretKey,
		issuer:        issuer,
		expirationHrs: expirationHrs,
	}
}

// GenerateToken genera un nuevo token JWT
func (s *JWTService) GenerateToken(userID string, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.expirationHrs))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken valida un token JWT y retorna los claims
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verificar que el método de firma sea HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

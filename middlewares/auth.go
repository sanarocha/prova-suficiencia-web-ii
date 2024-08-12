package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("senha")

// middleware de auth
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// obtendo o header de autenticação
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"erro": "Authorization header é necessário"})
			c.Abort()
			return
		}

		// remove a palavra Bearer e espaços em branco
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		// chama função auxiliar para validar o token
		token, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"erro": "Token inválido"})
			c.Abort()
			return
		}

		// verifica se a claim é válida e tenta acessar o userId
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["user_id"])
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"erro": "Token inválido"})
			c.Abort()
			return
		}

		// continua para próximo handler
		c.Next()
	}
}

// criação do token
func CreateToken(userID string) (string, error) {
	// token JWT usando o método de assinatura HMAC-SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		// claim de validade
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// assina o token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// validar token JWT: assinatura e validade
func ValidateToken(tokenString string) (*jwt.Token, error) {
	// parse do token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validando se o método de assinatura do token é um método HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Método inesperado")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// verifica se o claims é válido
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// verifica se está expirado
		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return nil, errors.New("Token expirado")
			}
		}
		return token, nil
	}

	return nil, errors.New("Token inválido")
}

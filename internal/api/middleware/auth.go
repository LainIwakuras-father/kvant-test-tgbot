package middleware

import (

	"github.com/gin-gonic/gin"
	
)


type AuthMiddlewareSecretKey struct {
	secretKey string
}

func NewAuthMiddleware(secretKey string) *AuthMiddlewareSecretKey {
	return &AuthMiddlewareSecretKey{
		secretKey: secretKey,
	}
}

func (m *AuthMiddlewareSecretKey) Validate(c *gin.Context){
	// Получаем ключ из заголовка Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Authorization header is required",
		})
		return
	}
	// Проверяем формат: Bearer {token}
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Invalid authorization format. Use: Bearer {token}",
		})
		return
	}
	//  Извлекаем токен
	token := authHeader[7:]

	// Сравниваем с секретным ключом
	if token != m.secretKey {
		c.AbortWithStatusJSON(401, gin.H{
			"error": "Invalid authorization token",
		})
		return
	}

	// Если всё ок - пропускаем дальше
	c.Next()
}
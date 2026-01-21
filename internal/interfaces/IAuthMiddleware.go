package interfaces

import "github.com/gin-gonic/gin"

type AuthMiddleware interface{
	Validate(*gin.Context)
}
package middlewares

import (
	"os"

	"github.com/gin-gonic/gin"
	"neft.web/auth"
	"neft.web/models"
)

func RequireAuth() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "neftAuth, XMLHttpRequest, Content-Type, Content-Length")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		tokenString := context.GetHeader("neftAuth")

		if tokenString == "" {
			context.JSON(401, gin.H{"error": models.ERR_JWT_TOKEN_REQUIRED.Error()})
			context.Abort()
			return
		}
		err := auth.ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}

		if os.Getenv("maitenance") == "true" {
			context.JSON(503, gin.H{"maitenance": true})
			context.Abort()
			return
		}

		context.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		context.Writer.Header().Set("Access-Control-Allow-Headers", "neftAuth, XMLHttpRequest, Content-Type, Content-Length")
		context.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(204)
			return
		}
		if os.Getenv("maitenance") == "true" {
			context.AbortWithStatus(503)
			return
		}
		context.Next()
	}
}

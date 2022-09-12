package controllers

import "github.com/gin-gonic/gin"

type Message struct {
	Wiki    string
	Data    interface{}
	Message string
}

func sendAnswer(responseCode int, context *gin.Context, data interface{}) {
	message := Message{
		Wiki:    "https://neftsec.atlassian.net/wiki/spaces/NW/pages/3342337/V1",
		Data:    data,
		Message: "",
	}
	context.JSON(responseCode, message)
}

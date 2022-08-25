package controllers

import "os"

type messageAnswer struct {
	Answer  interface{}
	Message string `json:"message"`
}

func answerV1(middle interface{}) interface{} {
	return messageAnswer{
		Answer:  middle,
		Message: os.Getenv("message"),
	}
}

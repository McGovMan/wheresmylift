package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type Empty struct {
}

type ErrorsArray struct {
	Errors []string `json:"errors" example:"cannot ping database,scheduler offline"`
}

type Error struct {
	Error string `json:"error" example:"a server error was encountered"`
}

type Message struct {
	Message string `json:"message" example:"i just wanted to say hi"`
}

func RespondWithError(c *gin.Context, err error, statusCode int) {
	if err == nil || len(err.Error()) == 0 {
		err = errors.New("unknown error")
	}
	c.JSON(statusCode, Error{Error: err.Error()})
}

func RespondWithString(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{"message": message})
}

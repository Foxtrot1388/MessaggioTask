package server

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type response struct {
	Message string
}

var notCreateMessage = response{
	Message: "failed to add message",
}

var notGetStatistic = response{
	Message: "failed to get statistic",
}

var emptyMessage = errors.New("Message is empty")

func (con *httpServer) logResponseWithError(c *gin.Context, status int, response interface{}, err error) {
	con.log.Debug("response "+c.FullPath(), "body", response, "status", status, "error", err)
	c.JSON(status, response)
}

func (con *httpServer) logResponse(c *gin.Context, status int, response interface{}) {
	con.log.Debug("response "+c.FullPath(), "body", response, "status", status)
	c.JSON(status, response)
}

package server

import (
	"time"

	"github.com/gin-gonic/gin"
)

func getQueryParamDate(c *gin.Context, key string, def time.Time) time.Time {
	str, ok := c.GetQuery(key)
	if ok {
		date, err := time.Parse("2006-01-02", str)
		if err != nil {
			return def
		} else {
			return date
		}
	} else {
		return def
	}
}

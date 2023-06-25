package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllSwitches(c *gin.Context) {
	res, err := GetAllSwitchesWithState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error")

	}
	c.IndentedJSON(200, res)
}

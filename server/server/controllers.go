package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sas604/IOT-fun/server/plug"
)

func GetAllSwitches(c *gin.Context) {
	res, err := GetAllSwitchesWithState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error")

	}
	c.IndentedJSON(200, res)
}

type switchState struct {
	Name  string `json:"name"`
	State string `json:"state"`
	Id    string `json:"id"`
}

func SetSwitch(c *gin.Context) {
	var r switchState
	err := c.ShouldBindJSON(&r)
	if (err) != nil {
		c.JSON(http.StatusInternalServerError, "error processing post request")
	}
	s := plug.NewSwitch(r.Id)
	err = s.SetSwitchState(r.State)
	if (err) != nil {
		c.JSON(http.StatusInternalServerError, "error processing post request")
	}
	c.JSON(200, "Change Succesfull")
}

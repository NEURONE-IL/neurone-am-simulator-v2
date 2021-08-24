package controller

import (
	"fmt"
	"net/http"
	"neurone-am-simulator-v2/generation"
	"neurone-am-simulator-v2/memory"
	"neurone-am-simulator-v2/model"

	"neurone-am-simulator-v2/util"

	"github.com/gin-gonic/gin"
)

type Controller struct {
}

func (controller *Controller) Routes(base *gin.RouterGroup) {

	base.POST("/init/:name", controller.CreateSimulation())
	base.GET("/stop/:name", controller.StopSimulation())
}

func (controller *Controller) CreateSimulation() func(c *gin.Context) {
	return func(c *gin.Context) {
		var configuration model.Configuration

		err := c.Bind(&configuration)

		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo decodificar json", err))
			return
		}

		name := c.Param("name")

		err = generation.SimulateNeurone(configuration, name)

		if err != nil {
			c.JSON(http.StatusInternalServerError, util.GetError("Simulation failed", err))
			return
		}
		c.JSON(http.StatusOK, "ok")
	}
}

func (controller *Controller) StopSimulation() func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")

		fmt.Println(name)

		memory.ActivateChannel(name)

		c.JSON(http.StatusOK, "ok")

	}
}

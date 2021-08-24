package controller

import "github.com/gin-gonic/gin"

var controller Controller

func Routes(base *gin.RouterGroup) {

	controller.Routes(base)
}

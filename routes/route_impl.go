package routes

import (
	"task_queue/controllers"

	"github.com/gin-gonic/gin"
)

type routeImpl struct {
	Controller controllers.QueueController
	Router     *gin.RouterGroup
}

func NewRoute(controller controllers.QueueController, router *gin.RouterGroup) TaskQueueRoute {
	return &routeImpl{Controller: controller, Router: router}
}

func (r *routeImpl) Serve() {
	r.Router.POST("/", r.Controller.SetQueue)
}

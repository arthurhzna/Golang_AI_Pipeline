package controllers

import (
	"github.com/gin-gonic/gin"
)

type QueueController interface {
	SetQueue(*gin.Context)
}

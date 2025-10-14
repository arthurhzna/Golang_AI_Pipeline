package controllers

import (
	"fmt"
	"net/http"
	"task_queue/domain/dto"
	"task_queue/services"

	"task_queue/common/response"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type QueueControllerImpl struct {
	service services.QueueService
}

func NewQueueController(service services.QueueService) QueueController {
	return &QueueControllerImpl{service: service}
}

func (c *QueueControllerImpl) SetQueue(ctx *gin.Context) {
	var request dto.QueueRequest
	err := ctx.ShouldBindWith(&request, binding.FormMultipart)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusBadRequest,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errMessage := err.Error()
		response.HttpResponse(response.ParamHTTPResp{
			Code:    http.StatusBadRequest,
			Message: &errMessage,
			Gin:     ctx,
		})
		return
	}
	fmt.Println("debug request", request)
	result, err := c.service.SetQueue(ctx, &request)
	if err != nil {
		response.HttpResponse(response.ParamHTTPResp{
			Code: http.StatusInternalServerError,
			Err:  err,
			Gin:  ctx,
		})
		return
	}

	response.HttpResponse(response.ParamHTTPResp{
		Code: http.StatusOK,
		Data: result,
		Gin:  ctx,
	})

}

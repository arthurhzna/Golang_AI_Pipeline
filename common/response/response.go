package response

import (
	"net/http"
	"task_queue/constants"
	errConstant "task_queue/constants/error"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message any         `json:"message"`
	Data    interface{} `json:"data"`
}

type ParamHTTPResp struct {
	Code    int
	Err     error
	Message *string
	Gin     *gin.Context
	Data    interface{}
}

func HttpResponse(param ParamHTTPResp) {
	if param.Err == nil {
		param.Gin.JSON(param.Code, Response{
			Status:  constants.Success,
			Message: http.StatusText(http.StatusOK),
			Data:    param.Data,
		})
		return
	}

	message := errConstant.ErrInternalServerError.Error()
	if param.Message != nil {
		message = *param.Message
	} else if param.Err != nil {
		if errConstant.ErrMapping(param.Err) {
			message = param.Err.Error()
		}
	}

	param.Gin.JSON(param.Code, Response{
		Status:  constants.Error,
		Message: message,
		Data:    param.Data,
	})
	return
}

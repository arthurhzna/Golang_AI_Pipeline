package middlewares

import (
	"net/http"
	"os"
	"task_queue/common/response"
	"task_queue/constants"
	errConstant "task_queue/constants/error"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("Recovered from panic: %v", r)
				c.JSON(http.StatusInternalServerError, response.Response{
					Status:  constants.Error,
					Message: errConstant.ErrInternalServerError.Error(),
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func APIKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		validAPIKey := os.Getenv("API_KEY")

		if validAPIKey == "" {
			logrus.Warn("API_KEY not set in environment, skipping validation")
			c.Next()
			return
		}

		apiKey := c.GetHeader("x-api-key")

		if apiKey == "" {
			response.HttpResponse(response.ParamHTTPResp{
				Code: http.StatusUnauthorized,
				Err:  errConstant.ErrInvalidAPIKey,
				Gin:  c,
			})
			c.Abort()
			return
		}

		if apiKey != validAPIKey {
			logrus.Warnf("Invalid API key attempt: %s", apiKey)
			response.HttpResponse(response.ParamHTTPResp{
				Code: http.StatusUnauthorized,
				Err:  errConstant.ErrInvalidAPIKey,
				Gin:  c,
			})
			c.Abort()
			return
		}

		logrus.Infof("Valid API key access from %s", c.ClientIP())
		c.Next()
	}
}

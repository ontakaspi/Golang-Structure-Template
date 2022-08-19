package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/ztrue/tracerr"
	"golang-example/libraries/logger"
	"net/http"
)

type errorHandler struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var errData interface{}
				var httpCode int

				switch err.(type) {
				case map[string]interface{}:
					recoverErr := err.(map[string]interface{})
					httpCode = recoverErr["httpCode"].(int)
					errData = recoverErr["error"]
					var errorData error
					if recoverErr["errorData"] != nil {
						errorData = recoverErr["errorData"].(error)
					} else {
						errorData = nil
					}

					if httpCode >= 500 {
						logger.SetLogFileAndConsole(logger.LogData{
							Message: "unexpected error",
							CustomFields: logrus.Fields{
								"data": errData,
							},
							Level: "ERROR",
						})
						if errorData != nil {
							dataErrr := tracerr.Wrap(errorData)
							tracerr.PrintSourceColor(dataErrr)
						}

					}
				case string:
					httpCode = http.StatusInternalServerError
					errData = errorHandler{
						Status:  "error",
						Message: err,
					}
					if httpCode >= 500 {
						logger.SetLogFileAndConsole(logger.LogData{
							Message: "unexpected error",
							CustomFields: logrus.Fields{
								"data": errData,
							},
							Level: "ERROR",
						})
						dataErrr := tracerr.New(err.(string))
						tracerr.PrintSourceColor(dataErrr)
					}
				default:
					httpCode = http.StatusInternalServerError
					errData = errorHandler{
						Status:  "error",
						Message: err.(error).Error(),
					}
					if httpCode >= 500 {
						logger.SetLogFileAndConsole(logger.LogData{
							Message: "unexpected error",
							CustomFields: logrus.Fields{
								"data": errData,
							},
							Level: "ERROR",
						})

						dataErrr := tracerr.Wrap(err.(error))
						tracerr.PrintSourceColor(dataErrr)
					}
				}

				context.JSON(httpCode, errData)
			}
		}()
		context.Next()
	}
}

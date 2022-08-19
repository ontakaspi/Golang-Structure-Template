package httpResponse

import (
	"github.com/gin-gonic/gin"
	"golang-example/app/helper"
	"net/url"
)

type SuccessDataResp struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type SuccessResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type ErrorResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type ValidationErrorResp struct {
	Status string                 `json:"status"`
	Errors helper.SimplifiedError `json:"errors"`
}
type ValidationErrorRespFormData struct {
	Status string     `json:"status"`
	Errors url.Values `json:"errors"`
}

// ============= Succes Response =============

func HttpCreated(c *gin.Context, message string, data interface{}) {
	c.JSON(201, SuccessDataResp{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Ok(c *gin.Context, message string, data interface{}) {
	c.JSON(201, SuccessDataResp{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func Accepted(c *gin.Context, message string) {
	c.JSON(202, SuccessResp{
		Status:  "Progress",
		Message: message,
	})
}

func NoContent(c *gin.Context) {
	c.Status(204)
}

// ============= Err Response =============

func TooManyRequests() {
	panic(map[string]interface{}{
		"error": ErrorResp{
			Status:  "Too Many Requests",
			Message: "You have made too many requests in a given amount of time. Please try again later.",
		},
		"httpCode": 429,
	})
}
func BadRequest(message helper.SimplifiedError) {
	panic(map[string]interface{}{
		"error": ValidationErrorResp{
			Status: "Bad Request",
			Errors: message,
		},
		"httpCode": 400,
	})
}
func BadRequestFormData(message url.Values) {
	panic(map[string]interface{}{
		"error": ValidationErrorRespFormData{
			Status: "Bad Request",
			Errors: message,
		},
		"httpCode": 400,
	})
}

func Unauthorized() {
	panic(map[string]interface{}{
		"error": ErrorResp{
			Status:  "Unauthorized",
			Message: "You are not authorized to access this resource",
		},
		"httpCode": 401,
	})
}
func Forbidden() {
	panic(map[string]interface{}{
		"error": ErrorResp{
			Status:  "Forbidden",
			Message: "You are not allowed to access this resource",
		},
		"httpCode": 403,
	})
}
func NotFound(message error) {
	panic(map[string]interface{}{
		"error": ErrorResp{
			Status:  "Not Found",
			Message: message.Error(),
		},
		"httpCode": 404,
	})
}
func Conflict(message error) {
	panic(map[string]interface{}{
		"error": ErrorResp{
			Status:  "Conflict",
			Message: message.Error(),
		},
		"httpCode": 409,
	})
}

func InternalServerError(message error) {

	panic(map[string]interface{}{
		"error": ErrorResp{
			Status:  "error",
			Message: message.Error(),
		},
		"httpCode":  500,
		"errorData": message,
	})
}

func ErrorWithHttpCode(error error, httpCode int) {
	switch httpCode {
	case 400:
		BadRequest(helper.SimplifyError(error))
	case 401:
		Unauthorized()
	case 403:
		Forbidden()
	case 404:
		NotFound(error)
	case 409:
		Conflict(error)
	case 500:
		InternalServerError(error)
	case 429:
		TooManyRequests()
	}
}

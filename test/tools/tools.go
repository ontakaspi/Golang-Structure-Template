package testTools

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/ztrue/tracerr"
	"golang-example/libraries/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

type errorHandler struct {
	Status  string      `json:"status"`
	Message interface{} `json:"message"`
}
type Expected struct {
	ExpectedStatusCode int
	ExpectedStatusBody string
	RequestBody        interface{}
	Error              bool
	Params             []gin.Param
	Query              url.Values
	AuthToken          string
}

var DBMock *gorm.DB
var Mock sqlmock.Sqlmock

func SetupDatabase() {

	var (
		db  *sql.DB
		err error
	)
	db, Mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	DBMock, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	os.Setenv("UNIT_TEST", "1")
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

func NewHtppRecorder(method string, params gin.Params, u url.Values, token string, content interface{}) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c := GetTestGinContext(w)
	c.Request.Method = method
	c.Request.Header.Set("Content-Type", "application/json")
	if token != "" {
		c.Request.Header.Set("Authorization", "Bearer "+token)
	}
	// set path params
	c.Params = params
	// set query params
	c.Request.URL.RawQuery = u.Encode()
	// set request body
	if content != nil {
		jsonbytes, err := json.Marshal(content)
		if err != nil {
			panic(err)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(jsonbytes))
	}

	return w, c
}

func CatchPanic(t *testing.T, w *httptest.ResponseRecorder, ScenarioExpected Expected, err interface{}) {

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
				Message: "Unexpected Error",
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
				Message: "Unexpected Error",
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
				Message: "Unexpected Error",
				CustomFields: logrus.Fields{
					"data": errData,
				},
				Level: "ERROR",
			})

			dataErrr := tracerr.Wrap(err.(error))
			tracerr.PrintSourceColor(dataErrr)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write([]byte(fmt.Sprintf("%v", errData)))

	assert.EqualValues(t, ScenarioExpected.ExpectedStatusCode, w.Code)
	assert.EqualValues(t, ScenarioExpected.ExpectedStatusBody, w.Body.String())

}

func RunControllerTest(t *testing.T, testName string, method string,
	expectedScenario Expected, function func(c *gin.Context)) bool {
	return t.Run(testName, func(t *testing.T) {
		var w *httptest.ResponseRecorder
		var c *gin.Context
		var AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

		// check if method is allowed
		if !utils.Contains(AllowedMethods, method) {
			panic("method not allowed")
		}
		w, c = NewHtppRecorder(method, expectedScenario.Params, expectedScenario.Query,
			expectedScenario.AuthToken, expectedScenario.RequestBody)

		if expectedScenario.Error {
			defer func() {
				if err := recover(); err != nil {
					CatchPanic(t, w, expectedScenario, err)
				}
			}()
			function(c)
		} else {
			function(c)
			assert.Equal(t, expectedScenario.ExpectedStatusCode, w.Code)
			assert.Equal(t, expectedScenario.ExpectedStatusBody, w.Body.String())
		}
	})
}

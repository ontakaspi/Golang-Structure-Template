package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang-example/app/helper"
	"golang-example/app/models/request"
	"golang-example/app/service"
	"golang-example/libraries/httpResponse"
	"strconv"
)

type exampleController struct {
	exampleService service.ExampleService
	authService    service.AuthService
	validator      *validator.Validate
}

func NewExampleController(
	exampleService service.ExampleService,
	authService service.AuthService,
	validator *validator.Validate,
) *exampleController {
	return &exampleController{
		exampleService: exampleService,
		authService:    authService,
		validator:      validator,
	}
}

func (global *exampleController) CreateExampleData(c *gin.Context) {

	/*-- Check user permission with decoded JWT Token --*/
	//checkUserRoles := global.authService.UserHasRoles(c, "backend_services")
	//if !checkUserRoles {
	//	httpResponse.Forbidden()
	//	return
	//}

	/*-- Validating project id params from segment --*/

	Request := request.CreateExampleData{}
	Request.BindRequestField(c)
	errorValidation := helper.ValidateFormData(c.Request, Request.Rules, Request.Message)
	if errorValidation != nil {
		httpResponse.BadRequestFormData(errorValidation)
	}
	requestData := Request.RequestCreateExampleData

	/*-- Print data request with Example Service --*/
	ExampleData := global.exampleService.CreateExampleData(requestData)

	httpResponse.HttpCreated(c, "Success create example data", ExampleData)
	return

}

func (global *exampleController) GetExampleDatas(c *gin.Context) {

	/*-- Check user permission with decoded JWT Token --*/
	//checkUserRoles := global.authService.UserHasRoles(c, "backend_services")
	//if !checkUserRoles {
	//	httpResponse.Forbidden()
	//	return
	//}
	ExampleDatas := global.exampleService.GetExampleDatas()
	httpResponse.Ok(c, "Success get example data", ExampleDatas)

}
func (global *exampleController) GetExampleDataById(c *gin.Context) {

	/*-- Check user permission with decoded JWT Token --*/
	//checkUserRoles := global.authService.UserHasRoles(c, "backend_services")
	//if !checkUserRoles {
	//	httpResponse.Forbidden()
	//	return
	//}

	/*-- Validating project id params from segment --*/
	ExampleDataId, err := strconv.Atoi(c.Param("example_data_id"))
	if err != nil {
		httpResponse.BadRequest(helper.SimplifyError(errors.New("paramater project_misconfig_id must be type of integer (/api/v1/example_data/:example_data_id)")))
		return
	}
	ExampleData := global.exampleService.GetExampleDataById(ExampleDataId)
	httpResponse.Ok(c, "Success get example data", ExampleData)

}

func (global *exampleController) EditExampleData(c *gin.Context) {

	/*-- Check user permission with decoded JWT Token --*/
	//checkUserRoles := global.authService.UserHasRoles(c, "backend_services")
	//if !checkUserRoles {
	//	httpResponse.Forbidden()
	//	return
	//}

	/*-- Validating project id params from segment --*/
	ExampleDataId, err := strconv.Atoi(c.Param("example_data_id"))
	if err != nil {
		httpResponse.BadRequest(helper.SimplifyError(errors.New("paramater project_misconfig_id must be type of integer (/api/v1/example_data/:example_data_id)")))
		return
	}

	//*-- Binding and validating data from request body --*/
	var exampleDataRequest request.ExampleRequest
	errorValidation := helper.BindAndValidate(c, &exampleDataRequest, global.validator)
	if errorValidation != nil {
		httpResponse.BadRequest(errorValidation)
	}

	ExampleData := global.exampleService.EditExampleData(ExampleDataId, exampleDataRequest)
	httpResponse.HttpCreated(c, "Success edit example data", ExampleData)

}

func (global *exampleController) DeleteExampleData(c *gin.Context) {

	/*-- Check user permission with decoded JWT Token --*/
	//checkUserRoles := global.authService.UserHasRoles(c, "backend_services")
	//if !checkUserRoles {
	//	httpResponse.Forbidden()
	//	return
	//}
	ExampleDataId, err := strconv.Atoi(c.Param("example_data_id"))
	if err != nil {
		httpResponse.BadRequest(helper.SimplifyError(errors.New("paramater project_misconfig_id must be type of integer (/api/v1/example_data/:example_data_id)")))
		return
	}
	global.exampleService.DeleteExampleData(ExampleDataId)
	httpResponse.NoContent(c)
}

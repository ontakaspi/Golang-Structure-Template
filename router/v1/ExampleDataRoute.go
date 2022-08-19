package route

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang-example/app/controller"
	"golang-example/app/repository"
	"golang-example/app/service"
	"golang-example/config/database"
)

func ExampleRoute(router *gin.RouterGroup) {
	validate := validator.New()
	PostgreDB := database.PostgreDB

	exampleRepository := repository.NewExampleRepository(PostgreDB)
	exampleService := service.NewExampleService(exampleRepository)
	authService := service.NewAuthService()
	exampleController := controller.NewExampleController(exampleService, authService, validate)

	//ExampleData
	router.POST("/example-data", exampleController.CreateExampleData)
	router.GET("/example-data", exampleController.GetExampleDatas)
	router.GET("/example-data/:example_data_id", exampleController.GetExampleDataById)
	router.PUT("/example-data/:example_data_id", exampleController.EditExampleData)
	router.DELETE("/example-data/:example_data_id", exampleController.DeleteExampleData)
}

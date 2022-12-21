This repository is a structured golang project that is used to generate a new backend/API service application by using the golang gin framework.
# Table of contents
1. [Project Requirment](#project-requirment)
2. [Project Structure](#project-structure)
    1. [Root Directory](#root-directory)
    2. [Route Directory](#route-directory)
    3. [Controller Directory](#controller-directory)
    4. [Service Directory](#service-directory)
    5. [Repository Directory](#repository-directory)
    6. [Model Directory](#model-directory)
    7. [Helper Directory](#helper-directory)
    8. [Middleware Directory](#middleware-directory)
    9. [Library Directory](#library-directory)
    10. [Config Directory](#config-directory)
3. [Database Migration](#database-migration)
4. [Run Project](#run-project)

# Project Requirment <a name="project-requirment"></a>

1. Required for this local run in gcc https://jmeubank.github.io/tdm-gcc/
2. Install golang ^v18 here https://go.dev/dl/ and follow the <a href="https://go.dev/doc/install">installation instructions</a>
3. Install postgresql for database transaction (skip if you want to use docker)
4. Install docker (if you want to run this locally)


# Project Structure <a name="project-structure"></a>

### Root Directory <a name="root-directory"></a>
``/`` (root dir)
this folder contain ``.env`` file that store environment data from host, ``_public_key.pem`` for authorize data token JWT (if have jwt Authorization) and file ``main.go`` of The main entrance of the API for setup environment settings, systems,port, etc .
<details><summary>Example .env</summary>

````dotenv
PUBLIC_KEY_PATH="_public_key.pem"
API_PATH_VERSION="/v1"
DB_HOST=${DB_HOST}
DB_PORT=5432
DB_NAME='app_db'
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
PORT=${APP_PORT}


````
</details>
<details><summary>Example main.go</summary>

````go
    package main
        import (
        "fmt"
        "net/http"
        "gopkg.in/gin-gonic/gin.v1"
        "articles/services/mysql"
        "articles/routers/v1"
        "articles/core/models"
    )
    
    var router  *gin.Engine;
    
    func init() {
        mysql.CheckDB()
        router = gin.New();
    router.NoRoute(noRouteHandler())
        version1:=router.Group("/v1")
        v1.InitRoutes(version1)
    
    }
    
    func main() {
        fmt.Println("Server Running on Port: ", 9090)
        http.ListenAndServe(":9090",router)
    }
````
</details>

---
### Route Directory <a name="route-directory"></a>
`/routers` This package will store every routes in your REST API.
The reason separate the handler is, to easy us to manage each routers. So we can create comments about the API , that with apidoc will generate this into structured documentation. Then we will call the function in index.go in current package.
Example:<br>
<details><summary>Example code ExampleRoute.go</summary>

```go
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
```
</details>
<details><summary>Example code routers.go</summary>

```go
package router

import (
	"github.com/gin-gonic/gin"
	"golang-example/app/middleware"
	route "golang-example/router/v1"
)

// InitRoutesJWT function route that use JWT midlleware
func InitRoutesJWT(g *gin.RouterGroup) {
	// Initialize Midlleware
	g.Use(middleware.ErrorHandler())
	g.Use(middleware.JSONMiddleware())
	g.Use(middleware.AuthorizeJWT())
	// Initialize route
	route.ExampleRoute(g)

}

// InitRoutes function route for home or some url not using a JWT Auth
func InitRoutes(g *gin.RouterGroup) {
	g.Use(middleware.ErrorHandler())
	g.Use(middleware.JSONMiddleware())
	// Initialize route
	route.SetHomeRoutes(g)
}

```
</details>

-----
### Controller Directory <a name="controller-directory"></a>
``/app/controllers``
this package will store every controllers in your REST API, and will be used in the routers.
the controller will be used to handle the request and response to the client.
<details><summary>Example Code</summary>

```go
package controller

import (
"errors"
"github.com/gin-gonic/gin"
"github.com/go-playground/validator/v10"
"strconv"
"golang-example/app/helper"
"golang-example/app/models/request"
"golang-example/app/service"
"golang-example/libraries/httpResponse"
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
	checkUserRoles := global.authService.UserHasRoles(c, "backend_services")
	if !checkUserRoles {
		httpResponse.Forbidden()
		return
	}

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
```


</details>


-----
### Service Directory <a name="service-directory"></a>
``/app/service``
this package will store every services in your REST API, and will be used in the controllers.
the service will create a logic for handling the request and response to the client and pass the data to the controller.
<details><summary>Example Code</summary>

```go
package service

import (
	"errors"
	"gorm.io/gorm"
	"golang-example/app/models/entity"
	"golang-example/app/models/request"
	"golang-example/app/repository"
	"golang-example/libraries/httpResponse"
)

type exampleService struct {
	exampleRepository repository.ExampleRepository
}

func NewExampleService(exampleRepository repository.ExampleRepository) ExampleService {
	return &exampleService{
		exampleRepository: exampleRepository,
	}
}

type ExampleService interface {
	CreateExampleData(ExampleData request.RequestCreateExampleData) entity.ExampleData
	GetExampleDatas() []entity.ExampleData
	GetExampleDataById(ExampleDataId int) entity.ExampleData
	EditExampleData(ExampleDataId int, ExampleDataRequest request.ExampleRequest) entity.ExampleData
	DeleteExampleData(ExampleDataId int)
}

func (s *exampleService) CreateExampleData(ExampleDataRequest request.RequestCreateExampleData) (ExampleData entity.ExampleData) {

	//Create Example Data
	ExampleData, err := s.exampleRepository.CreateExampleData(ExampleDataRequest)
	if err != nil {
		httpResponse.InternalServerError(err)
	}
	return ExampleData
}
```
</details>

-----
### Repository Directory <a name="repository-directory"></a>
``/app/repository``
this package will store every repositories in your REST API, and will be used in the services.
the repository will handling data from services and do the CRUD operation to the database.
<details><summary>Example Code</summary>

```go
package repository

import (
	"golang-example/app/models/entity"
	"golang-example/app/models/request"
	"gorm.io/gorm"
)

type exampleRepository struct {
	PostgreDB *gorm.DB
}

func NewExampleRepository(PostgreDB *gorm.DB) ExampleRepository {
	return &exampleRepository{PostgreDB: PostgreDB}
}

type ExampleRepository interface {
	CreateExampleData(ExampleDataRequest request.RequestCreateExampleData) (ExampleData entity.ExampleData, err error)
	GetExampleDatas() (ExampleDatas []entity.ExampleData, err error)
	GetExampleDataById(ExampleDataId int) (ExampleData entity.ExampleData, err error)
	EditExampleData(ExampleData entity.ExampleData) (err error)
	DeleteExampleData(ExampleDataId int) (err error)
}

func (r *exampleRepository) CreateExampleData(ExampleDataRequest request.RequestCreateExampleData) (ExampleData entity.ExampleData, err error) {
	ExampleData = entity.ExampleData{
		Name:    ExampleDataRequest.Name,
		Age:     ExampleDataRequest.Age,
		Address: ExampleDataRequest.Address,
	}
	if err := r.PostgreDB.Create(ExampleData).Error; err != nil {
		return ExampleData, err
	}
	return ExampleData, nil
}
```


</details>

-----
### Models Directory <a name="models-directory"></a>
``/app/models/entity`` This package will store all created model struct for using as data transcaction in database or just as transactional data. we use gorm as ORM library for handling data in database.
<details><summary>Example Code</summary>


```go
package entity

import "gorm.io/gorm"

type ExampleData struct {
	gorm.Model
	Name    string `gorm:"type:varchar(255);not null"`
	Age     int    `gorm:"type:int;not null"`
	Address string `gorm:"type:varchar(255);not null"`
}

```


</details>

``/app/models/request`` This package will store all request as struct data for validate data from API body request. 
We use package validator as validation library for handling data from request.
<ul>
<li>
Reference for validator library is <a href="https://pkg.go.dev/github.com/go-playground/validator/v10">https://pkg.go.dev/github.com/go-playground/validator/v10</a>
</li>
<li>
the reference for validating form-data type is <a href="https://github.com/thedevsaddam/govalidator">https://github.com/thedevsaddam/govalidator</a>.
</li>
</ul>

<details><summary>Example Code</summary>

```go
package request
import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// ExampleRequest is a struct for example using validator v10
type ExampleRequest struct {
	Name    string `json:"name" validate:"required,min=3"`
	Age     int    `json:"age" validate:"required,gte=0,lte=130"`
	Address string `json:"address" validate:"required,min=3"`
}

type RequestCreateExampleData struct {
	Name    string
	Age     int
	Address string
}
// CreateExampleData is a struct for example using thedevsaddam/govalidator
type CreateExampleData struct {
	Rules                    map[string][]string
	Message                  map[string][]string
	RequestCreateExampleData RequestCreateExampleData
}

func (std *CreateExampleData) BindRequestField(c *gin.Context) {

	std.Rules = make(map[string][]string)
	std.Rules["name"] = []string{"required"}
	std.Rules["age"] = []string{"required"}
	std.Rules["address"] = []string{"required"}
	std.RequestCreateExampleData.Name = c.PostForm("name")
	std.RequestCreateExampleData.Age, _ = strconv.Atoi(c.PostForm("age"))
	std.RequestCreateExampleData.Address = c.PostForm("address")

}
```


</details>

``/app/models/response`` This package will store all response as struct data for giving API body response.


----
### Helper Directory <a name="helper-directory"></a>
``/app/helper`` This pacakge will store every function that will reusuable in any function in controller.
included helper:
<ul>
<li>
`ChiperHelper.go` helper for encrypting or decrypting data.
</li>
<li>
`ErrorHelper.go` is helper for error handling, reference for error handling is <a href="https://go.dev/blog/error-handling-and-go">https://go.dev/blog/error-handling-and-go</a>
</li>
<li>
`JWTHelper.go` is helper for JWT token, reference for JWT token is <a href="github.com/dgrijalva/jwt-go">github.com/dgrijalva/jwt-go</a>
</li>
<li>
`SimplifyError.go` helper for convert error response package validator <a href="https://pkg.go.dev/github.com/go-playground/validator/v10">v10</a> to human readable.
</li>
<li>
`validate.go` helper for validate data from request.
</li>
</ul>


----
### Middleware Directory <a name="middleware-directory"></a>
``/app/middlewares``
This package will store every middeleware that will use in routes.
included helper:
<ul>
<li>
`errorHandler.go` middleware for error handling that use in routes(gin routes).
</li>
<li>
`JWTMiddleware.go` middleware for JWT authorization that use in routes(gin routes).
</li>
<li>
`JSONMiddleware.go` middleware for giving JSON response that use in routes(gin routes).
</li>
</ul>

-----
### Library Directory <a name="library-directory"></a>
`/libraries`  This package will store any library that used in projects. But only for manually created/imported library, that not available when using go get package_name commands. Could be your own hashing algorithm, graph, tree etc.
include library:
<ul>
<li>
`httpResponse.go` library for giving response to client based on status code and message. The library will give response in JSON format using gin context default like:
`c.JSON(202, SuccessResp{
		Status:  "Progress",
		Message: message
	})
`` and will return panic if some status code error occur. the panic will be handled by errorHandler middleware.
</li>
<li>
`looger.go` library for logging data to file. this package require logrus library.
</li>
</ul>

-----
### Config Directory <a name="config-directory"></a>
``/config`` This package will store any configuration and setting to used in project from any used service, could be mongodb,redis,mysql, elasticsearch, etc.
`/config/database`  This package for database configuration. File `migrations` is for database migration that use package gorm.io/gorm. the migration will be run when project start.

-----
### Test Directory <a name="config-directory"></a>
``/test/mockDatabase`` This package will store any mockDatabase repository.

`/test/tools/`  This package for any tools that used for unit testing, included tools:
<ul>
<li>
`tools.go` a tools for setting driver mock, run test any manymore.
</li>
</ul>
for more information about mock driver, you can see unit testing documentation in UnitTesting.md

-----
# Database Migration <a name="database-migration"></a>
Database migration is a process of creating and updating database tables to match the current model definitions.

1. for using database migration, you need to import package gorm.io/gorm. 
2. Create entity struct and use gorm to create table in `/models/entity` folder. for example:
    ```go
    package entity
    
    import "gorm.io/gorm"
    
    type ExampleData struct {
        gorm.Model
        Name    string `gorm:"type:varchar(255);not null"`
        Age     int    `gorm:"type:int;not null"`
        Address string `gorm:"type:varchar(255);not null"`
    }
    ```
    *the `gorm.Model` is a for defining table primary key,created_at,updated_at,deleted_at.*
3. Add this line code to `/config/database/migrations` file function `Migrate()`:
    ```go
    package database
    
    import "golang-example/app/models/entity"
    
    func Migrate() {
        db := PostgreDB
        err := db.AutoMigrate(&entity.ExampleData{})
        if err != nil {
            return
        }
    }
    ```
4. Start project it will run migration automatically.
*<br>AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes. It will change existing column’s type if its size, precision, nullable changed. It WON’T delete unused columns to protect your data. for more detail about gorm, please refer to https://gorm.io/docs/migration.html*

# Run Project <a name="run-project"></a>
1. Create database in postgresql.
2. Edit `.env` file and set database configuration based on your database configuration. 
   example:
   ```dotenv
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=postgres
   ```
3. in terminal go to project folder and run command:
    ```shell
    go mod download
    ```
4. then run the project:
    ```shell
    go run main.go
    ```
5. after application run, go to your POSTMAN application, in workspace, click on `import` button and select `File` option.
6. Browser the collection file (`Example Data.postman_collection.json`) in project directory and select it.
7. click on `import` button and all API will show in POSTMAN.
8. Then you can test your API.

### Maintainer : <a href="https://www.linkedin.com/in/kasyfi-assegaf/">Muhammad Kasfi </a>
If you are interested in becoming a maintainer please reach out to me 
https://github.com/ontakaspi 

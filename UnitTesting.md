*   [Create MockDatabase repository](#GolangUnitTesting-CreateMockDatabaserepository)
*   [Create Test Controller](#GolangUnitTesting-CreateTestController)
*   [Create Test Repository](#GolangUnitTesting-CreateTestRepository)
*   [Helper Tools](#GolangUnitTesting-HelperTools)
   *   [Tools model Data(struct) Scenario](#GolangUnitTesting-ToolsmodelData(struct)Scenario)
   *   [Tool Setup Database](#GolangUnitTesting-ToolSetupDatabase)
   *   [Tool Run Controller Test](#GolangUnitTesting-ToolRunControllerTest)
*   [Referensi](#GolangUnitTesting-Referensi)

Create MockDatabase repository
------------------------------
*   Making a mock database is placed in the `test/mockDatabase` folder with the name **MockRepositoryName.java** (because usually the database is called in the repository)


![](/doc-image/48693339.png)

*   Creating a mock database name is adjusted to the method on the repository interface (ex `findById`) plus with the expected return mock (ex \``NotFound`\`) like : `findByIdNotFound`.

*   Creation of mock process with library [sqlMock](https://github.com/DATA-DOG/go-sqlmock): `testTools.Mock.ExpectQuery(StringQuery).WillReturnRows(sqlmock.NewRows(dataRow)).AddRow([]driver.Value{dataReturn}));`, where `StringQuery` filled with the query from the expected, `dataRow` are the fields that will be returned, then `dataReturn` is the dummy data result of the query.

*   Creating a mock with sqlMock is adjusted according to the return from the method, if the method returns a struct `Proect` then we create the dummy data and then we put it in the code `.AddRow(Project);`. **Make sure the dummy data is the data you are looking for/expected with the return from the mock.**

    ```java
    package mockDatabase
    
    import (
    	"github.com/DATA-DOG/go-sqlmock"
    	"regexp"
    	"time"
    	ConstVarString "project-service/app/helper/const"
    	testTools "project-service/test/tools"
    )
    
    var optionsRowProjectInformation = []string{
    	"id",
    	"created_at",
    	"updated_at",
    	"deleted_at",
    	"project_id"
    	}
    
    func GetByProjectIdAndFound() {
    	// Create mock data project_informations
    	timea, _ := time.Parse(ConstVarString.StrinFullTimeFormat, ConstVarString.StrinQueryMockTime)
    	testTools.Mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "project_informations" WHERE project_id = $1 AND "project_informations"."deleted_at" IS NULL ORDER BY "project_informations"."id" LIMIT 1`)).
    		WithArgs(1).
    		WillReturnRows(sqlmock.NewRows(optionsRowProjectInformation).
    			AddRow(1, timea, timea, nil, 1))
    }
    ```


* * *

Create Test Controller
----------------------

*   The test cases are in the same folder, or package, where the code being tested is located.
    //insert picture


*   ![](/doc-image/48595207.png)

    The file name ends in with **\_Test (ProjectController\_Test.go)**

*   The function name start with **Test (**`TestGetScanResult`**)**

*   Make initialization of the controller function call in the file (**The creation of this controller init function is only done once in each controller\_x.go**). When creating, you need passing some data interface service, repository, or a struct into the controller's constructor, such as when setting up a route like this:


```go
func initProjectController() *controller.ProjectController {
	validate := validator.New()
	db := testTools.DBMock

	projectInformationRepository := repository.NewProjectInformationRepository(db)
	projectService := serviceProject.NewProjectService(projectInformationRepository)
	authService := service.NewAuthService()

	services := &controller.ProjectControllerServices{
		ProjectService:                 projectService,
		GeneralService:                 ApiGeneralService,

		AuthService:                    authService,
	}

	projectController := controller.NewProjectController(services, validate)
	return projectController
}
```

-   To make the test function by giving parameters `(t *testing.T)` for each function you want to create. The naming of the function test is the same as the function controller being tested plus a prefix `Test`. **Every test is required to call a function** :

    -   `testTools.SetupDatabase()` : for mock database driver settings

    -   `projectController := initProjectController()` : Calling the initialization of the controller according to the needs of each test.


```go
func TestGetScanResult(t *testing.T) {
	testTools.SetupDatabase()
	projectController := initProjectController()
}
```

*   Before running the test, data requests are required, therefore it is necessary to create a data parameter, query parameter, authToken or request body and also an expected scenario with struct `testTools.Expected`:


_\*IF auth token empty string, it will not set the Authorizaton bearer token_

```go
/* Scenario 1 */
expectedScenario2 := testTools.Expected{
		ExpectedStatusCode: http.StatusAccepted,
		ExpectedStatusBody:       `{"status":"Progress","message":"image scan not implemented yet"}`,
		RequestBody: null,
		Error:      false,
		Params: []gin.Param{{
			Key:   "project_id",
			Value: "1",
		}},
		Query: url.Values{"scan_type": []string{"image"}},
		AuthToken: stringToken,
}
```

*   after creating the parameter data and expected scenario, we run the test by calling the helper `testTools.RunControllerTest` by passing parameters :

    *   `t` : data struct package testing (`t *testing.T`)

    *   `scenarioName` : the name of the test to be run in string form

    *   `"GET"` : the method that will run for this test

    *   `expected` : data struct scenario that was created earlier (`testTools.Expected`)

    *   `projectController.GetScanResult` : is a function controller that has parameters `(c *gin.Context)` **(function must be a controller).**


```go
testTools.RunControllerTest(t, scenarioName,"GET", expected, projectController.GetScanResult)
```

*   If the controller is going to execute a repository, make sure that the repository method is covered by the mock database call as follows :

    *   there is a unit\_test env call, this env is automatically set when calling the test helper `testTools.GetTestGinContext`

    *   **In condition unit test = true** , call the mockdatabase function that corresponds to the repository.

       ```go
       func (r *projectInformationRepository) GetByProjectId(projectId int64) (*entity.ProjectInformation, error) {
           unitTest, _ := strconv.Atoi(os.Getenv("UNIT_TEST"))
           if unitTest == 1 {
               //call test data from test function
               mockDatabase.GetByProjectIdAndFound()
           }
           projectInformation := new(entity.ProjectInformation)
           if err := r.db.Where(ConstVarString.StringQueryProjectId, projectId).First(projectInformation).Error; err != nil {
               return projectInformation, err
           }
       
           return projectInformation, nil
       }
       ```

 When creating a test controller, all executable repository methods must call a mock database or you will get a mock sql error “**all expectations were already fulfilled**”
 
* If there is a function that has other service dependencies (such as apiservice.general) it is made to return mock data according to what is needed in that function.

   *   There is a unit\_test env call, this env is automatically set when calling the test helper `testTools.GetTestGinContext`

   *   **In condition unit test = true** , return from the desired expected data.


Example function get generalService

```go
func (apiGeneralService *ApiGeneralServiceStruct) GetProjectGeneral(JWTToken string, projectId int32) (project *entity.Project, httpCode int, Error error) {
	unitTest, _ := strconv.Atoi(config.GetEnv("UNIT_TEST"))
	if unitTest == 1 {
		//proses data mock
		project = &entity.Project{
			ID:   int32(projectId),
			Url:  "xxxx",
			Key:  "xxxx",
			Organization: entity.Organization{
				ID:  1,
				Key: "xxxx",
			},
		}
		return project, http.StatusOK, nil
	} else {
		//proses data real
		urlGeneralService := "url"
		req, err1 := http.NewRequest("GET", urlGeneralService+"/project/"+strconv.Itoa(int(projectId)), nil)
		if err1 != nil {return nil, 500, err1}
		req.Header.Add("Authorization", JWTToken)
		client := &http.Client{}
		responseGet, err := client.Do(req)
        responseData, errRead := ioutil.ReadAll(responseGet.Body)
		errUnMarshal := json.Unmarshal(responseData, &project)
		return project, http.StatusOK, nil
	}

}
```

If you don't create mock data in this function, testing will depend on it and require the service to be up. So Functions that have other service dependencies (such as apiservice.general) **must be made to return mock data as needed.**

![](images/icons/grey_arrow_down.png)Example full code test: ProjectController\_test.go

```go
package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/url"
	"testing"
	"project-service/app/repository"
	"project-service/app/service"
	"project-service/app/service/apiService"
	serviceProject "project-service/app/service/project"
	testTools "project-service/test/tools"
)
func initProjectController() *controller.ProjectController {
	validate := validator.New()
	db := testTools.DBMock

	projectInformationRepository := repository.NewProjectInformationRepository(db)
	projectService := serviceProject.NewProjectService(projectInformationRepository)
	authService := service.NewAuthService()

	services := &controller.ProjectControllerServices{
		ProjectService:                 projectService,
		GeneralService:                 ApiGeneralService,

		AuthService:                    authService,
	}

	projectController := controller.NewProjectController(services, validate)
	return projectController
}

//auth token
const stringToken = "12345677"

func TestGetScanResult(t *testing.T) {
	testTools.SetupDatabase()
	projectController := initProjectController()
	var ScenarioExpected []testTools.Expected

	/* Scenario 1 */
	expectedScenario1 := testTools.Expected{
		ExpectedStatusCode: http.StatusBadRequest,
		ExpectedStatusBody:       `{Bad Request [paramater project id must be type of integer (/api/v1/project/{project_id}/scan/status?scan_type=xxx)]}`,
		RequestBody: nil,
		Error:      true,
		Params: []gin.Param{{
			Key:   "project_id",
			Value: "sss"}},
		Query:     url.Values{"scan_type": []string{"repository"}},
		AuthToken: stringToken,
	}
	/* Scenario 2 */
	expectedScenario2 := testTools.Expected{
		ExpectedStatusCode: http.StatusAccepted,
		ExpectedStatusBody:       `{"status":"Progress","message":"image scan not implemented yet"}`,
		RequestBody: nil,
		Error:      false,
		Params: []gin.Param{{
			Key:   "project_id",
			Value: "1",
		}},
		Query:     url.Values{"scan_type": []string{"image"}},
		AuthToken: stringToken,
	}

	ScenarioExpected = append(ScenarioExpected, expectedScenario1, expectedScenario2, expectedScenario3)

	for index, expected := range ScenarioExpected {
		scenarioName := fmt.Sprintf("TestGetScanResult %d", index+1)
		testTools.RunControllerTest(t, scenarioName,"GET", expected, projectController.GetScanResult)
	}
}
```

* * *

Create Test Repository
----------------------

*   The test cases are in the same folder, or package, where the code being tested is located in the repository folder.

*   ![](/doc-image/49184820.png)

    The file name ends in with **\_Test (FilterUsedInScanResultRepository\_test.go)**

*   The function name start with **Test** and combined with the name of the repository method.

*   **Every test must call** the init of the database driver with `testTools.SetupDatabase()`

*   In creating the test, just **call the mockdatabase that has been created** then call the method of the repository

*   before call the function we need to s**et the db driver** from helper testtools with `db := testTools.DBMock`

*   to comparing the data we just create the condition if they err or the value not that as expected (**custumizing the logic by your need**)


```go
func TestGetProjectFilterOptionBySeverityValue(t *testing.T) {
	testTools.SetupDatabase()

	type SeverityValues struct {
		SeverityValue string
	}
	severityValues := []SeverityValues{
		{SeverityValue: "CRITICAL"},
		{SeverityValue: "UNKNOWN"},
	}

	for _, severityValue := range severityValues {
		t.Run(severityValue.SeverityValue, func(t *testing.T) {
			mockDatabase.GetProjectFilterOptionBySeverityValueAndFound(severityValue.SeverityValue)
			db := testTools.DBMock

			// Call function
			dataProjectFilterOptions, err := NewFilterUsedInScanResultRepository(db).GetProjectFilterOptionBySeverityValue(1, severityValue.SeverityValue, 1)
			if err != nil {
				t.Errorf(ConstVarString.StringErrorPrintWithVar, err)
			}
			if dataProjectFilterOptions.Value != severityValue.SeverityValue {
				t.Errorf(ConstVarString.StringErrorPrintWithVar, dataProjectFilterOptions)
			}
		})
	}
}
```

* * *

Helper Tools
------------

this file used for helping the test function, the file located in **test/tools/tools.go.**

### Tools model Data(struct) Scenario

![](/doc-image/49315884.png)

this struct is need for creating the scenario test in controller test, the field is :

*   **ExpectedStatusCode**`int` : the expected status code to compare with

*   **ExpectedStatusBody**`string`: return body that expected

*   **RequestBody**`interface{}` : the request body of this scenario, will convert to json data

*   **Error** `bool` : if this scenario produce error, the value is true nor false. This value needed because if the scenario expected error we need to cacth it.

*   **Params** `[]gin.Param`: the parameter data

*   **Query** `url.Values`: the query data

*   **AuthToken** `string`: the string token of bearer, IF auth token empty string, it will not set the Authorizaton bearer token


### Tool Setup Database

![](/doc-image/49446968.png?width=340)

this function will procude setting up the database using _gorm_, and the driver use library _sqlMock._ **Call this everythime mockdatabase will use.**

### Tool Run Controller Test

![](/doc-image/49741869.png?width=448)

this function used for running the test for controller. This function will create a **newHtppRecorder** by gin, **create context** by gin and **will comparing the expected body and status code,** the paramater :

*   `t` : data struct package testing (`t *testing.T`)

*   `scenarioName` : the name of the test to be run in string form

*   `method` : the method that will run for this test

*   `expected` : data struct scenario that was created earlier (`testTools.Expected`)

*   `projectController.GetScanResult` : is a function controller that has parameters `(c *gin.Context)` **(function must be a controller).**


* * *

Reference
---------

[https://blog.canopas.com/golang-unit-tests-with-test-gin-context-80e1ac04adcd](https://blog.canopas.com/golang-unit-tests-with-test-gin-context-80e1ac04adcd)

[https://medium.easyread.co/unit-test-sql-in-golang-5af19075e68e](https://medium.easyread.co/unit-test-sql-in-golang-5af19075e68e)

[https://www.youtube.com/watch?v=t9QJPE5vwhs&t=3859s](https://www.youtube.com/watch?v=t9QJPE5vwhs&t=3859s)

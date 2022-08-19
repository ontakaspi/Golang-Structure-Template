package request

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// ExampleRequest is a struct for example
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

package service

import (
	"errors"
	"golang-example/app/models/entity"
	"golang-example/app/models/request"
	"golang-example/app/repository"
	"golang-example/libraries/httpResponse"
	"gorm.io/gorm"
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

func (s *exampleService) GetExampleDatas() []entity.ExampleData {

	//Get Example Datas
	ExampleDatas, err := s.exampleRepository.GetExampleDatas()
	if err != nil {
		httpResponse.InternalServerError(err)
	}
	return ExampleDatas
}

func (s *exampleService) GetExampleDataById(ExampleDataId int) (ExampleData entity.ExampleData) {

	//Get Example Data
	ExampleData, err := s.exampleRepository.GetExampleDataById(ExampleDataId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httpResponse.NotFound(errors.New("Example Data not found"))
		} else {
			httpResponse.InternalServerError(err)
		}

	}
	return ExampleData
}

func (s *exampleService) EditExampleData(ExampleDataId int, ExampleDataRequest request.ExampleRequest) (ExampleData entity.ExampleData) {

	//Get Example Data
	ExampleData, err := s.exampleRepository.GetExampleDataById(ExampleDataId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			httpResponse.NotFound(errors.New("Example Data not found"))
		} else {
			httpResponse.InternalServerError(err)
		}

	}
	ExampleData.Name = ExampleDataRequest.Name
	ExampleData.Age = ExampleDataRequest.Age
	ExampleData.Address = ExampleDataRequest.Address
	//Update Example Data
	if err := s.exampleRepository.EditExampleData(ExampleData); err != nil {
		httpResponse.InternalServerError(err)
	}
	return ExampleData
}

func (s *exampleService) DeleteExampleData(ExampleDataId int) {

	//Delete Example Data
	if err := s.exampleRepository.DeleteExampleData(ExampleDataId); err != nil {
		httpResponse.InternalServerError(err)
	}
}

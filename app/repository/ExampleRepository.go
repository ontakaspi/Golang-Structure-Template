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
	if err := r.PostgreDB.Create(&ExampleData).Error; err != nil {
		return ExampleData, err
	}
	return ExampleData, nil
}

func (r *exampleRepository) GetExampleDatas() (ExampleDatas []entity.ExampleData, err error) {
	if err := r.PostgreDB.Find(&ExampleDatas).Error; err != nil {
		return ExampleDatas, err
	}
	return ExampleDatas, nil
}

func (r *exampleRepository) GetExampleDataById(ExampleDataId int) (ExampleData entity.ExampleData, err error) {
	if err := r.PostgreDB.First(&ExampleData, ExampleDataId).Error; err != nil {
		return ExampleData, err
	}
	return ExampleData, nil
}

func (r *exampleRepository) EditExampleData(ExampleData entity.ExampleData) (err error) {
	if err := r.PostgreDB.Save(&ExampleData).Error; err != nil {
		return err
	}
	return nil
}

func (r *exampleRepository) DeleteExampleData(ExampleDataId int) (err error) {
	if err := r.PostgreDB.Delete(&entity.ExampleData{}, ExampleDataId).Error; err != nil {
		return err
	}
	return nil
}

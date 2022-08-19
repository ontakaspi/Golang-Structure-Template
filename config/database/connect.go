package database

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"golang-example/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strconv"
)

// DB Declare the variable for the database
var PostgreDB *gorm.DB

// ConnectDB connect to db
func ConnectDB() {
	var err error

	loggers := logrus.New()
	p := config.GetEnv("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		log.Println("Error parse the port")
	}
	if port == 0 {
		loggers.Info("DB_PORT is not set, use default port 5432")
		port = 5432
	}

	// Connection URL to connect to Postgres Database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.GetEnv("DB_HOST"),
		port,
		config.GetEnv("DB_USER"),
		config.GetEnv("DB_PASSWORD"),
		config.GetEnv("DB_NAME"))
	// Connect to the DB and initialize the DB variable
	PostgreDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		loggers.Panic("failed to connect database postgreSQL: " + err.Error())
	}
	loggers.Info("Connection Opened to Database postgreSQL")

}

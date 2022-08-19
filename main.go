package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/itsjamie/gin-cors"
	"github.com/sirupsen/logrus"
	"golang-example/config"
	"golang-example/config/database"
	"golang-example/router"
	"net/http"
	"strconv"
	"time"
)

var route *gin.Engine

func init() {
	// Initialize logger
	config.InitLogFile()
	// Connect to the Database
	database.ConnectDB()
	// database.AutoMigrate()
	database.Migrate()

	//setup main routes
	route = gin.New()
	route.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE, OPTIONS",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		ValidateHeaders: false,
	}))
	route.GET("/", func(c *gin.Context) {
		pathActive := config.GetEnv("API_PATH_VERSION")
		c.JSON(http.StatusOK, gin.H{"data": "this is not valid endpoint, use active API endpoint ('" + pathActive + "')"})
	})

	//setup version 1 routes
	version := route.Group("/api/v1")
	// route Initialize
	router.InitRoutes(version)
	router.InitRoutesJWT(version)

	// Handler if no route define
	route.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "Page not found")
	})

}

func main() {

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		logrus.Error(err)
	}
	time.Local = loc // -> this is setting the global timezone

	loggerMain := logrus.New()
	formatter := &logrus.TextFormatter{
		FullTimestamp: true,
	}
	loggerMain.SetFormatter(formatter)
	loggerMain.Info("Setting global timezone to Asia/Jakarta")

	// Start the cron job
	cronJob()

	p := config.GetEnv("APP_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		loggerMain.Error("Error parse the port")
	}
	if port == 0 {
		loggerMain.Info("APP Port not set, defaulting to 8080")
		port = 8080
	}

	loggerMain.Info("Starting APP Service on Port: " + strconv.Itoa(int(port)))
	err = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(port)), limit(route))
	if err != nil {
		loggerMain.Info("Error starting APP Service on Port: " + strconv.Itoa(int(port)) + "\n" + err.Error())
		return
	}
}

func cronJob() {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	s := gocron.NewScheduler(loc)
	//downloadTrivyDB() every day at 12:00 AM
	go func() {
		// clean up visitor data
		s.Every(1).Minute().Do(cleanupVisitors)
	}()

	s.StartAsync()
}

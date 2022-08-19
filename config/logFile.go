package config

import (
	"fmt"
	"os"
	"time"
)

var LogFile *os.File
var LogFileName string

func InitLogFile() {
	year, month, _ := time.Now().Date()
	LogFileName = fmt.Sprintf("%d-%d.log", year, month)
	curDir, _ := os.Getwd()
	LogFile, _ = os.OpenFile(fmt.Sprintf(curDir)+"/log/"+LogFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return
}

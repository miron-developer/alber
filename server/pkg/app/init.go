/*
	Initialize app
*/

package app

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"alber/pkg/orm"
)

// Code struct for app
type Code struct {
	ExpireMin int
	Value     interface{}
}

// AppConfig - app's configurations
type AppConfig struct {
	PORT              string
	MOBIZON_API_KEY   string
	MAX_REQUEST_COUNT string
}

// Application - app config and items
type Application struct {
	m                   sync.Mutex
	ELog                *log.Logger
	ILog                *log.Logger
	CurrentRequestCount int
	CurrentMin          int // how many minuts pass after start/day
	UsersCode           map[string]*Code
	Config              *AppConfig
}

func checkFatal(eLogger *log.Logger, e error) {
	if e != nil {
		eLogger.Fatal(e)
	}
}

func GetConfigs() (*AppConfig, error) {
	content, e := os.ReadFile(".env")
	if e != nil {
		return nil, e
	}

	conf := &AppConfig{}
	confMap := map[string]interface{}{}
	rows := strings.Split(string(content), "\n")
	for _, row := range rows {
		arr := strings.Split(row, "=")
		confMap[arr[0]] = arr[1]
	}

	e = orm.FillStructFromMap(conf, confMap)
	return conf, e
}

// InitProg initialise
func InitProg() *Application {
	wd, _ := os.Getwd()
	logFile, _ := os.OpenFile(wd+"/logs/log_"+time.Now().Format("2006-01-02")+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	eLog := log.New(logFile, "\033[31m[ERROR]\033[0m\t", log.Ldate|log.Ltime|log.Lshortfile)
	iLog := log.New(logFile, "\033[34m[INFO]\033[0m\t", log.Ldate|log.Ltime|log.Lshortfile)
	iLog.Println("loggers is done!")

	iLog.Println("creating/configuring database")
	checkFatal(eLog, orm.InitDB(iLog))
	iLog.Println("database completed!")

	iLog.Println("configuring app")
	config, e := GetConfigs()
	checkFatal(eLog, e)
	iLog.Println("configuring done")

	return &Application{
		ELog:                eLog,
		ILog:                iLog,
		CurrentRequestCount: 0,
		CurrentMin:          0,
		UsersCode:           map[string]*Code{},
		Config:              config,
	}
}

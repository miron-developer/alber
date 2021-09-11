/*
	Initialize app
*/

package app

import (
	"log"
	"os"
	"strings"
	"sync"
	"zhibek/pkg/orm"
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
	IsHeroku            bool
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
	logFile, _ := os.Create("logs.txt")

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
		IsHeroku:            false,
		UsersCode:           map[string]*Code{},
		Config:              config,
	}
}

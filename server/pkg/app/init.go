/*
	Initialize app
*/

package app

import (
	"log"
	"os"
	"sync"
	"zhibek/pkg/orm"
)

// Code struct for app
type Code struct {
	ExpireMin int
	Value     interface{}
}

// Application this is app struct and items
type Application struct {
	m                   sync.Mutex
	ELog                *log.Logger
	ILog                *log.Logger
	Port                string
	CurrentRequestCount int
	CurrentMin          int // how many minuts pass after start/day
	MaxRequestCount     int
	IsHeroku            bool
	UsersCode           map[string]*Code
}

// InitProg initialise
func InitProg() *Application {
	logFile, _ := os.Create("logs.txt")

	eLog := log.New(logFile, "\033[31m[ERROR]\033[0m\t", log.Ldate|log.Ltime|log.Lshortfile)
	iLog := log.New(logFile, "\033[34m[INFO]\033[0m\t", log.Ldate|log.Ltime|log.Lshortfile)
	iLog.Println("loggers is done!")

	iLog.Println("creating/configuring database")
	orm.InitDB(eLog, iLog)
	iLog.Println("database completed!")

	return &Application{
		ELog:                eLog,
		ILog:                iLog,
		Port:                "4330",
		CurrentRequestCount: 0,
		CurrentMin:          0,
		MaxRequestCount:     1200,
		IsHeroku:            false,
		UsersCode:           map[string]*Code{},
	}
}

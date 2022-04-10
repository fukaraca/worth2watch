package main

import (
	"github.com/fukaraca/worth2watch/db"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
)

func main() {
	defer db.Conn.Close()
	defer db.Cache.Close()
	
	endpoints()
	log.Fatalln("router has encountered an error while main.run: ", r.Run(serverPort))

}

func init() {
	//logger middleware teed to log.file
	logfile, err := os.OpenFile("./logs/log.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Could not create/open log file")
	}
	errlogfile, err := os.OpenFile("./logs/err.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Could not create/open err log file")
	}
	gin.DefaultWriter = io.MultiWriter(logfile, os.Stdout)
	gin.DefaultErrorWriter = io.MultiWriter(errlogfile, os.Stdout)
	//starts with builtin Logger() and Recovery() middlewares
	r = gin.Default()

	db.ConnectDB()
	db.CheckIfInitialized()
	db.CreateRedisClient()
}

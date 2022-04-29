package main

import (
	"github.com/fukaraca/worth2watch/api"
	"github.com/fukaraca/worth2watch/db"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
	"io"
	"log"
	"os"
)

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
	api.R = gin.Default()

	//rate limiter
	rLimit := ratelimit.New(20)
	leakBucket := func(limiter ratelimit.Limiter) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			limiter.Take()
		}
	}
	api.R.Use(leakBucket(rLimit))
	api.R.Use(requestid.New())
	db.ConnectDB()
	db.CheckIfInitialized()
	db.CreateRedisClient()
}

package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

var r *gin.Engine

//os.Getenv("key_string")
var TIMEOUT = 5 * time.Second
var serverHost = "localhost"
var serverPort = ":8080"

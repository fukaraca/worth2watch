package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

//in production, passing the variable to related function or wrapping in a struct is advised inspite of creating global variable like this. see:go.dev

var Conn *pgxpool.Pool
var Cache *redis.Client
var err error
var TIMEOUT = 5 * time.Second

var db_Host = "127.0.0.1"
var db_Port = "5432"
var db_Name = "worth2watchdb"
var db_User = "postgres"
var db_Password = "postgres"
var redis_Host = "localhost"
var redis_Port = ":6379"
var redis_Password = ""
var redis_DB = 0

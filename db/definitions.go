package db

import (
	"github.com/fukaraca/worth2watch/config"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

//in production, passing the variable to related function or wrapping in a struct is advised inspite of creating global variable like this. see:go.dev

var Conn *pgxpool.Pool
var Cache *redis.Client
var err error
var TIMEOUT = config.GetEnv.GetDuration("TIMEOUT")

var db_Host = config.GetEnv.GetString("DB_HOST")
var db_Port = config.GetEnv.GetString("DB_PORT")
var db_Name = config.GetEnv.GetString("DB_NAME")
var db_User = config.GetEnv.GetString("DB_USER")
var db_Password = config.GetEnv.GetString("DB_PASSWORD")
var redis_Host = config.GetEnv.GetString("REDIS_HOST")
var redis_Port = config.GetEnv.GetString("REDIS_PORT")
var redis_Password = config.GetEnv.GetString("REDIS_PASSWORD")
var redis_DB = config.GetEnv.GetInt("REDIS_DB")

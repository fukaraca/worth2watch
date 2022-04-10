package db

import (
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

//in production, passing the variable to related function or wrapping in a struct is advised see:go.dev
var Conn *pgxpool.Pool
var Cache *redis.Client
var err error

var db_Host = "127.0.0.1"
var db_Port = "5432"
var db_Name = "worth2watchDB"
var db_User = "postgres"
var db_Password = "postgres"
var redis_Host = "localhost"
var redis_Port = ":6379"
var redis_Password = ""
var redis_DB = 0

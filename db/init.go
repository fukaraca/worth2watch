package db

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"strings"
)

//ConnectDB function creates a connection pool to PSQL DB.
func ConnectDB() {

	databaseURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_Host, db_Port, db_User, db_Password, db_Name)
	Conn, err = pgxpool.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalln("DB connection error:", err)
	}
	//check whether connection is ok or not
	err = Conn.Ping(context.Background())
	if err != nil {
		log.Fatalln("Ping to DB error:", err)
	}

}

//CreateRedisClient function creates a Redis Client
func CreateRedisClient() {
	client := redis.NewClient(&redis.Options{
		Addr:     redis_Host + redis_Port,
		Password: redis_Password,
		DB:       redis_DB,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln("redis ping error:", err)
	}
	log.Println(pong, " redis activated")
	Cache = client
	//enable notifications from redis
	Cache.ConfigSet(context.Background(), "notify-keyspace-events", "KEA")
}

//CheckIfInitialized functions checks existance of tables and creates if necessary.
func CheckIfInitialized() {
	textByte, err := os.ReadFile("./db/init.sql")
	if err != nil {
		log.Fatalln("init.sql file couldn't be read: ", err)
	}
	statements := strings.Split(string(textByte), ";\n")
	for _, statement := range statements {
		comm, err := Conn.Exec(context.Background(), statement)
		if err != nil {
			log.Println("checked for initial DB structure: ", err, comm.String())
		}
	}

}

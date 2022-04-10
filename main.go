package main

import (
	"github.com/fukaraca/worth2watch/api"
	"github.com/fukaraca/worth2watch/db"
	"github.com/fukaraca/worth2watch/model"
	"log"
)

func main() {
	defer db.Conn.Close()
	defer db.Cache.Close()

	api.Endpoints()
	log.Fatalln("router has encountered an error while main.run: ", model.R.Run(model.ServerPort))
}

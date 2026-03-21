package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"go.com/go-backend/api"
	db "go.com/go-backend/db/sqlc"
	"go.com/go-backend/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load app config: ", err.Error())
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to the database: ", err.Error())
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddr)

	if err != nil {
		log.Fatal("cannot start the server: ", err.Error())
	}
}

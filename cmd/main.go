package main

import (
	"fmt"
	"time"

	"../internal/config"
	"../internal/database"
	"../internal/server"
)

func main() {
	cfg, err := config.ReadConfig("../configs/config.toml")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	db, err := database.New(cfg.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := server.New(cfg.Server, db)

	fmt.Println(db)

	sleepTime := 30 * time.Second
	go server.ListCurrencies(db, sleepTime)

	s.Run()
}

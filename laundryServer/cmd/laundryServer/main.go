package main

import (
	"database/sql"
	"flag"
	"log"

	"laundryServer/internal/config"
	"laundryServer/internal/server"
	"laundryServer/internal/slots"
	"laundryServer/internal/users"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	configPath := flag.String("config", "", "Path to config")
	flag.Parse()

	if *configPath == "" {
		log.Fatalln("No config has been provided.")
	}

	config, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatalln("Failed to create config:", err)
	}

	database, err := sql.Open("sqlite3", config.DatabasePath)
	if err != nil {
		log.Fatalln("Failed to open database.")
	}

	users, err := users.NewUsers(database)
	if err != nil {
		log.Fatalln("Failed to create users:", err)
	}

	slots, err := slots.NewSlots(database)
	if err != nil {
		log.Fatalln("Failed to create slots:", err)
	}

	server, _ := server.NewServer(config.Host, config.Token, users, slots)
	log.Printf("Server started on http://localhost%s\n", config.Host)
	server.ListenAndServe()
}

package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if err := runItemAPI(); err != nil {
		log.Panicf("failed to runItemAPI: %v", err)
	}
}

func runItemAPI() error {
	var repository ItemRepository
	var err error
	dbType := os.Getenv("DB_TYPE")
	switch dbType {
	case "sqlite":
		dbName := mustGetenv("DB_NAME")
		log.Println("use sqlite. database name:", dbName)
		repository, err = NewSQLiteItemRepository(dbName)
		if err != nil {
			return fmt.Errorf("failed to NewSQLiteItemRepository: %v", err)
		}
	// case "mysql":
	// 	log.Println("use mysql")
	// 	repository, err = NewMySQLItemRepository("item_db")
	// 	if err != nil {
	// 		log.Panicf("failed to NewMySQLItemRepository: %v", err)
	// 	}
	default:
		return fmt.Errorf("invalid dbType: %s", dbType)
	}

	cache := NewCache(10000000)

	usecase := NewItemUsecase(repository, cache)

	config := ServerConfig{Port: useEnvOrDefault("SERVER_PORT", "8080")}
	server := NewServer(config, usecase)
	return server.Run()
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("environment variable '%s' not set", k)
	}

	return v
}

func useEnvOrDefault(key, def string) string {
	v := def
	if fromEnv := os.Getenv(key); fromEnv != "" {
		v = fromEnv
	}
	log.Printf("%s environment variable set", key)
	return v
}

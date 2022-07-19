package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if err := runItemAPI(); err != nil {
		log.Panicf("failed to runItemAPI: %v", err)
	}
}

func runItemAPI() error {
	var itemRepository ItemRepository
	var scoreRepository ScoreRepository
	var err error
	dbType := os.Getenv("DB_TYPE")
	switch dbType {
	case "sqlite":
		dbName := MustGetenv("DB_NAME")
		log.Println("use sqlite. database name:", dbName)
		itemRepository, err = NewSQLiteItemRepository(dbName)
		if err != nil {
			return fmt.Errorf("failed to NewSQLiteItemRepository: %v", err)
		}
		scoreRepository, err = NewSQLiteScoreRepository(dbName)
		if err != nil {
			return fmt.Errorf("failed to NewSQLiteScoreRepository: %v", err)
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

	var cache CacheRepository
	cacheType := UseEnvOrDefault("CACHE_TYPE", "memory")
	if cacheType == "memory" {
		log.Println("use in memory cache")
		cache = NewCache(StrToInt(UseEnvOrDefault("CACHE_TTL", "10")))
	} else if cacheType == "redis" {
		log.Println("use redis cache")
		cache = NewRedisCache(StrToInt(UseEnvOrDefault("CACHE_TTL", "10")))
	} else {
		return fmt.Errorf("invalid cache type: %v", cacheType)
	}

	itemUsecase := NewItemUsecase(itemRepository, cache)
	rankingUsecase := NewRankingUsecase(scoreRepository, cache)
	usecase := NewSearchUsecase(itemUsecase, rankingUsecase, cache)

	config := ServerConfig{Port: UseEnvOrDefault("SERVER_PORT", "8080")}
	server := NewServer(config, usecase)
	return server.Run()
}

func MustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("environment variable '%s' not set", k)
	}

	return v
}

func UseEnvOrDefault(key, def string) string {
	v := def
	if fromEnv := os.Getenv(key); fromEnv != "" {
		v = fromEnv
	}
	log.Printf("%s environment variable set", key)
	return v
}

func StrToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Panicf("environment variable '%s' cannot convert to int", s)
	}
	return i
}

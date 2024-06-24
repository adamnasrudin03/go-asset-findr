package configs

import (
	"os"
	"strings"
	"sync"
)

var (
	lock    = &sync.Mutex{}
	configs *Configs
)

func GetInstance() *Configs {
	lock.Lock()
	defer lock.Unlock()

	configs = &Configs{
		App: AppConfig{
			Name: getEnv("APP_NAME", "go-asset-findr"),
			Env:  getEnv("APP_ENV", "dev"),
			Port: getEnv("APP_PORT", "8000"),
		},
		DB: DbConfig{
			Host:        getEnv("DB_HOST", "127.0.0.1"),
			Port:        getEnv("DB_PORT", "5432"),
			DbName:      getEnv("DB_NAME", "my_db"),
			Username:    getEnv("DB_USER", "postgres"),
			Password:    getEnv("DB_PASS", ""),
			DbIsMigrate: getEnv("DB_IS_MIGRATE", "true") == "true",
		},
	}

	return configs
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return strings.TrimSpace(value)
	}
	return strings.TrimSpace(fallback)
}

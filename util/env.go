package util

import (
	"log"
	"os"
)

func ReadEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalln("Environment variable " + key + " not set")
	}
	return value
}

func ReadEnvWithDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

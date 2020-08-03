package main

import (
	"fmt"
	"os"
)

const (
	defaultBasePath string = "/Users/bllli-nuc/Downloads"
	basePathEnvKey  string = "FILE_OBSERVER_BASE_PATH"
)

func GetEnvOrDefault(envName, defaultValue string) string {
	var env = os.Getenv(envName)
	if env == "" {
		fmt.Printf("cannot get env %s use default: %s \n", envName, defaultValue)
		env = defaultValue
	} else {
		fmt.Printf("got %s: %s \n", envName, defaultValue)
	}
	return env
}

func GetBasePath() string {
	return GetEnvOrDefault(basePathEnvKey, defaultBasePath)
}

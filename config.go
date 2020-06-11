package main

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	LoggerConfig string
	Relays string
}

func CollectConfig() (config Config) {
	var missingEnv []string

	// LOG_LEVEL
	var envLoggerLevel = os.Getenv("LOG_LEVEL")
	if envLoggerLevel == "" {
		config.LoggerConfig = "<root>=INFO"
	} else {
		config.LoggerConfig = fmt.Sprintf("<root>=%s", strings.ToUpper(envLoggerLevel))
	}

	// RELAYS
	config.Relays = os.Getenv("RELAYS")
	if config.Relays == "" {
		missingEnv = append(missingEnv, "RELAYS")
	}

	// Validation
	if len(missingEnv) > 0 {
		msg := fmt.Sprintf("Environment variables missing: %v", missingEnv)
		panic(fmt.Sprint(msg))
	}

	return
}

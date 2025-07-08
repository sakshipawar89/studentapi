package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// HTTPServer holds the HTTP server address from config
type HTTPServer struct {
	Address string `yaml:"address" env:"HTTP_ADDRESS" env-default:":8080"`
}

// Config holds the full application configuration
type Config struct {
	Env         string     `yaml:"env" env:"ENV" env-required:"true" env-default:"development"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	HTTPServer  HTTPServer `yaml:"http_server" env-required:"true"`
}

// MustLoad loads the configuration and exits on failure
func MustLoad() *Config {
	var configPath string

	// First, check the CONFIG_PATH environment variable
	configPath = os.Getenv("CONFIG_PATH")

	// If not set, try reading from CLI flag
	if configPath == "" {
		flagPath := flag.String("config", "", "Path to the config file")
		flag.Parse()
		configPath = *flagPath

		if configPath == "" {
			log.Fatal("Config path is not set. Please set the CONFIG_PATH environment variable or use the -config flag.")
		}
	}

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist at path: %s", configPath)
	}

	// Load config using cleanenv
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err.Error())
	}

	return &cfg
}

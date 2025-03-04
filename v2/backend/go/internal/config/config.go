// package config

// import (
//     "os"
//     "github.com/joho/godotenv"
// )

// type Config struct {
//     Port             string
//     PythonServiceURL string
//     LogLevel         string
//     Environment      string
// }

// func Load() *Config {
//     godotenv.Load()

//     return &Config{
//         Port:             getEnv("SERVER_PORT", "8080"),
//         PythonServiceURL: getEnv("PYTHON_SERVICE_URL", "http://localhost:8000"),
//         LogLevel:         getEnv("LOG_LEVEL", "debug"),
//         Environment:      getEnv("GO_ENV", "development"),
//     }
// }

// func getEnv(key, fallback string) string {
//     if value, exists := os.LookupEnv(key); exists {
//         return value
//     }
//     return fallback
// }

// backend/go/internal/config/config.go
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Server contains server-related configuration
type Server struct {
	Port                   int `json:"port"`
	ReadTimeoutSeconds     int `json:"readTimeoutSeconds"`
	WriteTimeoutSeconds    int `json:"writeTimeoutSeconds"`
	IdleTimeoutSeconds     int `json:"idleTimeoutSeconds"`
	ShutdownTimeoutSeconds int `json:"shutdownTimeoutSeconds"`
}

// LLM contains LLM service-related configuration
type LLM struct {
	ServiceURL string `json:"serviceURL"`
}

// Config contains all configuration for the application
type Config struct {
	Server Server `json:"server"`
	LLM    LLM    `json:"llm"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	cfg := &Config{}
	
	// Default server settings
	cfg.Server.Port = 8080
	cfg.Server.ReadTimeoutSeconds = 30
	cfg.Server.WriteTimeoutSeconds = 60
	cfg.Server.IdleTimeoutSeconds = 120
	cfg.Server.ShutdownTimeoutSeconds = 20
	
	// Default LLM service settings
	cfg.LLM.ServiceURL = "http://localhost:8000"
	
	return cfg
}

// Load loads configuration from config.json if it exists, otherwise returns default config
func Load() (*Config, error) {
	cfg := DefaultConfig()
	
	// Try to load from file
	configPath := filepath.Join(".", "config.json")
	if _, err := os.Stat(configPath); err == nil {
		file, err := os.Open(configPath)
		if err != nil {
			return cfg, fmt.Errorf("error opening config file: %w", err)
		}
		defer file.Close()
		
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(cfg); err != nil {
			return cfg, fmt.Errorf("error decoding config file: %w", err)
		}
	}
	
	return cfg, nil
}
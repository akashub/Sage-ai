// // backend/go/internal/config/config.go
// package config

// import (
//     "os"
//     // "github.com/joho/godotenv"
// )

// type Config struct {
//     Port             string
//     PythonServiceURL string
//     LogLevel         string
//     Environment      string
// }

// func Load() *Config {
//     // Load .env file
//     // godotenv.Load()

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
    "os"
    "github.com/joho/godotenv"
)

type Config struct {
    Port             string
    PythonServiceURL string
    LogLevel         string
    Environment      string
}

func Load() *Config {
    godotenv.Load()

    return &Config{
        Port:             getEnv("SERVER_PORT", "8080"),
        PythonServiceURL: getEnv("PYTHON_SERVICE_URL", "http://localhost:8000"),
        LogLevel:         getEnv("LOG_LEVEL", "debug"),
        Environment:      getEnv("GO_ENV", "development"),
    }
}

func getEnv(key, fallback string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return fallback
}
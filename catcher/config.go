package catcher

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Configuration struct {
	HTTPPort     int `json:"http_port"`
	Host         string
	RootHost     string `json:"root_host"`
	FrontendDir  string `json:"frontend_dir"`
	Favicon      string `json:"favicon"`
	RedirectDest string `json:"redirect_dest"`
}

func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func defaultConfig() Configuration {
	httpPort, _ := strconv.Atoi(Getenv("HTTP_PORT", "8080"))
	return Configuration{
		HTTPPort:     httpPort,
		Host:         Getenv("HOST", "127.0.0.1"),
		RootHost:     Getenv("ROOT_HOST", "localhost"),
		FrontendDir:  Getenv("FRONTEND_DIR", "frontend/dist"),
		Favicon:      Getenv("FAVICON", "frontend/favicon.ico"),
		RedirectDest: Getenv("REDIRECT_DEST", ""),
	}
}

func LoadConfiguration(filename string) (*Configuration, error) {
	config := defaultConfig()

	if filename == "" {
		fmt.Println("Using default config+env")
		return &config, nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)

	fmt.Println("Using default config+file")
	return &config, err
}

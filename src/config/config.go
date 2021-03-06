package config

import (
	"encoding/json"
	"fmt"
	"os"

	"go-gcs/src/logger"
	"go-gcs/src/service/googlecloud"
	"go-gcs/src/service/googlecloud/pubsubprovider"
	"go-gcs/src/service/googlecloud/storageprovider"
)

// Config is the structure for server
type Config struct {
	GoogleCloud  *googlecloud.Config      `json:"googlecloud"`
	Storage      *storageprovider.Storage `json:"storage"`
	PubSub       *pubsubprovider.PubSub   `json:"pubsub"`
	JWTSecretKey string                   `json:"jwtSecretKey"`
	Logger       logger.LoggerConfig      `json:"logger"`

	// the version settings of the current application
	Version string `json:"version"`
}

// Read will read config file
func Read(path string) (c Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return c, fmt.Errorf("Failed to open the config file: %v", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c); err != nil {
		return c, fmt.Errorf("Failed to load the config file: %v", err)
	}

	return c, nil
}

// MustRead will read config path
func MustRead(path string) Config {
	c, err := Read(path)
	if err != nil {
		panic(err)
	}
	return c
}

package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Environment string        `yaml:"environment"`
	Storage     StorageConfig `yaml:"storage"`
	HttpConfig  HttpConfig    `yaml:"http_server"`
	AuthConfig  AuthConfig    `yaml:"auth"`
}
type StorageConfig struct {
	Address  string `yaml:"address"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
}
type HttpConfig struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
}
type AuthConfig struct {
	Salt       string        `yaml:"salt"`
	SignedKey  string        `yaml:"signed_key"`
	RefreshTtl time.Duration `yaml:"refresh_ttl"`
	AccessTtl  time.Duration `yaml:"access_ttl"`
}

func NewConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config path does not exist on %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatal(err)
	}
	return &config
}

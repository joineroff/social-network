package config

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type AuthConfig struct {
	Secret     string        `yaml:"token-secret"`
	AccessTTL  time.Duration `yaml:"accessTTL"`
	RefreshTTL time.Duration `yaml:"refreshTTL"`
}

type HttpConfig struct {
	Port   string `yaml:"port"`
	Domain string `yaml:"domain"`
}

type MysqlConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type Config struct {
	Debug   bool       `yaml:"debug"`
	AppName string     `yaml:"app_name"`
	Auth    AuthConfig `yaml:"auth"`
	Http    HttpConfig `yaml:"http"`
	Mysql   struct {
		Master MysqlConfig   `yaml:"master"`
		Slaves []MysqlConfig `yaml:"slaves"`
	} `yaml:"mysql"`
}

func NewConfig(path string) *Config {
	cfg := &Config{}

	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to open config.yml file: %v", err)
	}

	if err := yaml.Unmarshal(b, cfg); err != nil {
		log.Fatalf("failed to parse config.yml: %v", err)
	}

	return cfg
}

func NewDefaultConfig() *Config {
	cfg := &Config{}

	b, err := os.ReadFile("./config/default.yml")
	if err != nil {
		log.Fatalf("failed to open config.yml file: %v", err)
	}

	if err := yaml.Unmarshal(b, cfg); err != nil {
		log.Fatalf("failed to parse config.yml: %v", err)
	}

	return cfg
}

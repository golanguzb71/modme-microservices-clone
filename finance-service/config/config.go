package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type ServerConfig struct {
	Port int `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
	Action   string `yaml:"action"`
}

type GrpcConfig struct {
	UserService struct {
		Address string `yaml:"address"`
	} `yaml:"userService"`
	EducationService struct {
		Address string `yaml:"address"`
	} `yaml:"educationService"`
	FinanceService struct {
		Address string `yaml:"address"`
	} `yaml:"financeService"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Grpc     GrpcConfig     `yaml:"grpc"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

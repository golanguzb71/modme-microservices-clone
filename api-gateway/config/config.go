package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`

	Grpc struct {
		AuditingService struct {
			Address string `yaml:"address"`
		} `yaml:"auditing_service"`

		UserService struct {
			Address string `yaml:"address"`
		} `yaml:"user_service"`

		EducationService struct {
			Address string `yaml:"address"`
		} `yaml:"education_service"`

		LidService struct {
			Address string `yaml:"address"`
		} `yaml:"lid_service"`
		FinanceService struct {
			Address string `yaml:"address"`
		} `yaml:"finance_service"`
	} `yaml:"grpc"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config/config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

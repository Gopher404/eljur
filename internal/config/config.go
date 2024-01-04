package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Bind BindConfig `yaml:"bind"`
	DB   DBConfig   `yaml:"DB"`
}

type BindConfig struct {
	Ip   string `yaml:"ip"`
	Port string `yaml:"port"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
}

func GetConfig(path string) (*Config, error) {
	cnf := &Config{}
	if err := cleanenv.ReadConfig(path, cnf); err != nil {
		return nil, err
	}
	return cnf, nil
}

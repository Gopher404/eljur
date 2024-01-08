package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Bind BindConfig `yaml:"bind"`
	DB   DBConfig   `yaml:"DB"`
	SSO  SSOConfig  `yaml:"SSO"`
	Log  LogConfig  `yaml:"log"`
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

type SSOConfig struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	AppKey string `yaml:"app_key"`
}

type LogConfig struct {
	Type  string `yaml:"type"`
	Out   string `yaml:"out"`
	Level string `yaml:"level"`
}

func GetConfig(path string) (*Config, error) {
	cnf := &Config{}
	if err := cleanenv.ReadConfig(path, cnf); err != nil {
		return nil, err
	}
	return cnf, nil
}

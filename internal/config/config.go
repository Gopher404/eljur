package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Bind     BindConfig     `yaml:"bind"`
	DB       DBConfig       `yaml:"DB"`
	SSO      SSOConfig      `yaml:"SSO"`
	Log      LogConfig      `yaml:"log"`
	Schedule ScheduleConfig `yaml:"schedule"`
}

type BindConfig struct {
	Ip      string        `yaml:"ip"`
	Port    string        `yaml:"port"`
	TimeOut time.Duration `yaml:"time_out"`
}

type DBConfig struct {
	Driver   string `yaml:"driver"`
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

type VKSeverConfig struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type VKAPIConfig struct {
	Version  string        `yaml:"version"`
	GroupId  string        `yaml:"group_id"`
	CacheTTL time.Duration `yaml:"cache_TTL"`
}

type ScheduleConfig struct {
	GroupName string        `yaml:"group_name"`
	VKSever   VKSeverConfig `yaml:"vk_server"`
	VKAPI     VKAPIConfig   `yaml:"vk_api"`
}

func GetConfig(path string) (*Config, error) {
	cnf := &Config{}
	if err := cleanenv.ReadConfig(path, cnf); err != nil {
		return nil, err
	}
	return cnf, nil
}

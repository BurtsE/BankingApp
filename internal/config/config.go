package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ServerPort string `json:"port" yaml:"port"`
	JWTSecret  string `json:"jwt_secret" yaml:"jwt_secret"`
	LogLevel   string `json:"log_level" yaml:"log_level"`
	Postgres   `json:"postgres" yaml:"postgres"`
}

type Postgres struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Database string `json:"database" yaml:"database"`
}

func InitConfig() (*Config, error) {
	c := &Config{}
	data, err := os.ReadFile("./configs/app_config.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}


func GetJWTSecretKey() string {
	return getEnv("SECRET_KEY")
}
func getEnv(param string) string {
	val := os.Getenv(param)
	if val == "" {
		log.Fatalf("env not set: %s", param)
	}
	return val
}
 

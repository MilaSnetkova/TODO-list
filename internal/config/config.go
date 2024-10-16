package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)


type Config struct {
	ServerAddress string        `env:"SERVER_ADDRESS" env-default:"localhost:7540"` 
	Version       string        `env:"VERSION" env-default:"v1.0"`                
	DBFile        string        `env:"TODO_DBFILE" env-default:"./scheduler.db"`    
	TODOPassword  string        `env:"TODO_PASSWORD" env-default:"password"`           
	JWTSecret     string        `env:"JWT_SECRET" env-default:"your_secret_key"`       
	Timeout       time.Duration `env:"TIMEOUT" env-default:"5s"`                   
	Ticker        time.Duration `env:"TICKER" env-default:"1m"`                 
}


func MustLoad() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return &cfg, nil
}
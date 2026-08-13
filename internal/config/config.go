package config

import "github.com/ilyakaznacheev/cleanenv"

type DatabaseConfig struct {
	Path string `yaml:"path" env:"PATH" env-default:"db/flight-routes.db"`

	Mode string `yaml:"mode" env:"MODE" env-default:"rwc"`

	CacheSize int `yaml:"cache_size" env:"CACHE_SIZE" env-default:"2000"`

	SyncMode string `yaml:"sync_mode" env:"SYNC_MODE" env-default:"normal"`

	BusyTimeoutSeconds  int `yaml:"busy_timeout_seconds" env:"BUSY_TIMEOUT_SECONDS" env-default:"5"`
	QueryTimeoutSeconds int `yaml:"query_timeout_seconds" env:"QUERY_TIMEOUT_SECONDS" env-default:"30"`

	MaxOpenConns           int `yaml:"max_open_conns" env:"MAX_OPEN_CONNS" env-default:"1"`
	MaxIdleConns           int `yaml:"max_idle_conns" env:"MAX_IDLE_CONNS" env-default:"1"`
	ConnMaxLifetimeSeconds int `yaml:"conn_max_lifetime_seconds" env:"CONN_MAX_LIFETIME_SECONDS" env-default:"0"`
}

type ServerConfig struct {
	Port int    `yaml:"port" env:"PORT" env-default:"8080"`
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
}

type Config struct {
	DatabaseConfig DatabaseConfig `yaml:"database" env-prefix:"DB_"`
	ServerConfig   ServerConfig   `yaml:"server" env-prefix:"SERVER_"`
}

func New() (Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

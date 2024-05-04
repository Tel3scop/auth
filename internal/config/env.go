package config

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v10"
)

// Provider предоставляет интерфейс для получения конфигурации.
type Provider interface {
	Config() *Config
}

// Config общий конфиг
type Config struct {
	DebugMode   bool   `env:"DEBUG_MODE" envDefault:"false"`
	Environment string `env:"ENV" envDefault:"local"`
	Postgres    Postgres
	GRPC        GRPC
	HTTP        HTTP
	Swagger     Swagger
	Encrypt     Encrypt
}

// Config возвращаем сам конфиг
func (c Config) Config() *Config {
	return &c
}

// Postgres конфиг подключения к БД
type Postgres struct {
	Host               string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port               string `env:"POSTGRES_PORT" envDefault:"5432"`
	User               string `env:"POSTGRES_USER" envDefault:"root"`
	Password           string `env:"POSTGRES_PASSWORD" envDefault:"password"`
	Db                 string `env:"POSTGRES_DB" envDefault:"postgres"`
	SslMode            string `env:"POSTGRES_SSL_MODE" envDefault:"disable"`
	DSN                string `env:"POSTGRES_DSN"`
	MaxOpenConnections int    `env:"POSTGRES_MAX_OPEN_CONNS" envDefault:"100"`
}

// GRPC конфиг подключения к grpc
type GRPC struct {
	Host     string `env:"GRPC_HOST"`
	Port     string `env:"GRPC_PORT"`
	Protocol string `env:"GRPC_PROTOCOL"`
	Address  string
}

// HTTP конфиг подключения к grpc
type HTTP struct {
	Host    string `env:"HTTP_HOST" envDefault:"localhost"`
	Port    string `env:"HTTP_PORT" envDefault:"8080"`
	Address string
}

// Swagger конфиг
type Swagger struct {
	Host    string `env:"SWAGGER_HOST" envDefault:"localhost"`
	Port    string `env:"SWAGGER_PORT" envDefault:"8081"`
	Address string
	Timeout int `env:"SWAGGER_TIMEOUT" envDefault:"5"`
}

// Encrypt конфиг с секретами
type Encrypt struct {
	RefreshTokenSecretKey           string `env:"REFRESH_TOKEN_SECRET"`
	AccessTokenSecretKey            string `env:"ACCESS_TOKEN_SECRET"`
	RefreshTokenExpirationInMinutes int    `env:"REFRESH_TOKEN_EXPIRATION" envDefault:"60"`
	RefreshTokenExpiration          time.Duration
	AccessTokenExpirationInMinutes  int `env:"ACCESS_TOKEN_EXPIRATION" envDefault:"5"`
	AccessTokenExpiration           time.Duration
	AuthPrefix                      string `env:"AUTH_PREFIX" envDefault:"Bearer "`
}

// New создаем новый конфиг
func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("loading config from env is failed: %w", err)
	}
	buildDSN(&cfg.Postgres)
	cfg.GRPC.Address = net.JoinHostPort(cfg.GRPC.Host, cfg.GRPC.Port)
	cfg.HTTP.Address = net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port)
	cfg.Swagger.Address = net.JoinHostPort(cfg.Swagger.Host, cfg.Swagger.Port)
	cfg.Encrypt.AccessTokenExpiration = time.Duration(cfg.Encrypt.AccessTokenExpirationInMinutes) * time.Minute
	cfg.Encrypt.RefreshTokenExpiration = time.Duration(cfg.Encrypt.RefreshTokenExpirationInMinutes) * time.Minute

	return cfg, nil
}

func buildDSN(p *Postgres) {
	p.DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.User, p.Password, p.Host, p.Port, p.Db, p.SslMode)
}

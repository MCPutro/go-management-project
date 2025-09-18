package config

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
)

var (
	configOnce sync.Once
	cfg        config

	applicationConfigOnce sync.Once
	applicationCfg        ApplicationConfig

	databaseConfigOnce sync.Once
	databaseCfg        DatabaseConfig

	jwtConfigOnce sync.Once
	jwtCfg        JwtConfig
)

type Config interface {
	GetApplicationConfig() *ApplicationConfig
	GetDatabaseConfig() *DatabaseConfig
	GetJwtConfig() *JwtConfig
}

type config struct {
	Application ApplicationConfig `mapstructure:"Application"`
	Database    DatabaseConfig    `mapstructure:"Database"`
	Jwt         JwtConfig         `mapstructure:"Jwt"`
}

type ApplicationConfig struct {
	Name string `mapstructure:"Name"`
	Port string `mapstructure:"Port"`
}

type DatabaseConfig struct {
	PostgresSql PostgresConfig `mapstructure:"PostgresSQL"`
	MySql       MySqlConfig    `mapstructure:"MySQL"`
}

type PostgresConfig struct {
	Name         string `mapstructure:"Name"`
	Host         string `mapstructure:"Host"`
	Port         string `mapstructure:"Port"`
	DatabaseName string `mapstructure:"DatabaseName"`
	Username     string `mapstructure:"Username"`
	Password     string `mapstructure:"Password"`
	FormatDSN    string `mapstructure:"FormatDSN"`
	DSN          string
}

type MySqlConfig struct {
	Name         string `mapstructure:"Name"`
	Host         string `mapstructure:"Host"`
	Port         string `mapstructure:"Port"`
	DatabaseName string `mapstructure:"DatabaseName"`
	Username     string `mapstructure:"Username"`
	Password     string `mapstructure:"Password"`
	FormatDNS    string `mapstructure:"FormatDNS"`
	DNS          string
}

type JwtConfig struct {
	SecretKey          string `mapstructure:"SecretKey"`
	ExpirationInSecond int    `mapstructure:"ExpirationInSecond"`
}

func LoadConfig() (Config, error) {
	v := viper.New()
	v.SetConfigName("app-2-local") // nama file tanpa .yml
	v.SetConfigType("yaml")
	v.AddConfigPath("./properties") // lokasi folder properties

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	var dsn string
	// Construct DSN for Postgres
	dsn = fmt.Sprintf(cfg.Database.PostgresSql.FormatDSN,
		cfg.Database.PostgresSql.Username, cfg.Database.PostgresSql.Password,
		cfg.Database.PostgresSql.Host, cfg.Database.PostgresSql.Port, cfg.Database.PostgresSql.DatabaseName)
	cfg.Database.PostgresSql.DSN = dsn

	// Construct DNS for MySQL
	dsn = fmt.Sprintf(cfg.Database.MySql.FormatDNS,
		cfg.Database.MySql.Username, cfg.Database.MySql.Password,
		cfg.Database.MySql.Host, cfg.Database.MySql.Port, cfg.Database.MySql.DatabaseName)
	cfg.Database.MySql.DNS = dsn

	return &cfg, nil
}

func (c *config) GetApplicationConfig() *ApplicationConfig {
	applicationConfigOnce.Do(func() {
		applicationCfg = c.Application
	})
	return &applicationCfg
}

func (c *config) GetDatabaseConfig() *DatabaseConfig {
	databaseConfigOnce.Do(func() {
		databaseCfg = c.Database
	})

	return &databaseCfg
}

func (c *config) GetJwtConfig() *JwtConfig {
	jwtConfigOnce.Do(func() {
		jwtCfg = c.Jwt
	})
	return &jwtCfg
}

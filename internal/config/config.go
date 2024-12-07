package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port int
	}
	Database struct {
		URL      string
		User     string
		Password string
		PoolSize int
	}
	Migration struct {
		Directory string
	}
	JwtAuth struct {
		Key string
	}
}

func NewConfig() *Config {
	env := os.Getenv("APP_ENV")
	config := Config{}

	switch env {
	case "production":
		config.loadFromEnv()
	default:
		env = "development"
		config.loadFromYaml(env)
	}

	return &config
}

func (c *Config) loadFromEnv() {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatalf("Invalid server port value: %v", err)
	}
	c.Server.Port = port

	c.Database.URL = os.Getenv("DB_URL")
	c.Database.User = os.Getenv("DB_USER")
	c.Database.Password = os.Getenv("DB_PASS")
	poolSize, err := strconv.Atoi(os.Getenv("DB_POOLSIZE"))
	if err != nil {
		log.Fatalf("Invalid database pool size value: %v", err)
	}
	c.Database.PoolSize = poolSize

	c.Migration.Directory = os.Getenv("MIGRATIONSDIR")

	c.JwtAuth.Key = os.Getenv("SECRET_KEY")
}

func (c *Config) loadFromYaml(env string) {
	viper.SetConfigName(env)
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	c.Server.Port = viper.GetInt("server.port")

	c.Database.URL = viper.GetString("database.url")
	c.Database.User = viper.GetString("database.user")
	c.Database.Password = viper.GetString("database.password")
	c.Database.PoolSize = viper.GetInt("database.pool_size")

	c.Migration.Directory = viper.GetString("migration.directory")
	c.JwtAuth.Key = viper.GetString("jwt_auth.key")
}

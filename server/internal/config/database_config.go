package config

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadDatabaseConfig loads database configuration from environment variables
func LoadDatabaseConfig() (*DBConfig, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE") // Often "disable" for local dev, "require" for prod

	if dbHost == "" || dbPortStr == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return nil, fmt.Errorf("database environment variables (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME) must be set")
	}

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT value: %w", err)
	}

	if dbSSLMode == "" {
		dbSSLMode = "disable" // Default SSLMode if not set
	}

	return &DBConfig{
		Host:     dbHost,
		Port:     strconv.Itoa(dbPort),
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  dbSSLMode,
	}, nil
}
func LoadDatabase() *bun.DB {
	dbconfig, err := LoadDatabaseConfig()
	if err != nil {
		fmt.Print("unabel to load database")
	}
	dsn := buildDSN(*dbconfig)
	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	return db
}
func buildDSN(config DBConfig) string {
	return "postgres://" + config.User + ":" + config.Password +
		"@" + config.Host + ":" + config.Port + "/" +
		config.DBName + "?sslmode=" + config.SSLMode
}

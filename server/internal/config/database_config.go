
package config
import (
	"database/sql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type DBConfig struct {
 Host string 
 Port string 
 User string 
 Password string 
 DBName string 
 SSLMode string 

}

func LoadDatabaseConfig()*DBConfig{
	dbconfig:=DBConfig{
	Host: "localhost",
	Port : "5432",
	User: "postgres",
	Password: "123",
	DBName: "eduhub",
	SSLMode: "disable",
}
return &dbconfig
}
func LoadDatabase()*bun.DB{
	dbconfig :=LoadDatabaseConfig()

dsn :=buildDSN(*dbconfig)
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
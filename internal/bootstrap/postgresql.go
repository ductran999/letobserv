package bootstrap

import (
	"github.com/ductran999/dbkit"
	"github.com/ductran999/letobserv/configs"
	"github.com/ductran999/letobserv/pkg/dbconn"
	"gorm.io/gorm"
)

func ConnectOrderDB(env *configs.OrdersConfigEnv) (*gorm.DB, error) {
	pgConf := dbkit.PostgreSQLConfig{
		Config: dbkit.Config{
			Host:     env.PgHost,
			Port:     env.PgPort,
			Username: env.PgUser,
			Password: env.PgPass,
			Database: env.PgDB,
			TimeZone: "UTC",
		},
		SSLMode: "disable",
	}

	return dbconn.ConnectDB(pgConf)
}

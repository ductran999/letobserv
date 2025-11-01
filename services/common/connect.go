package common

import (
	"context"
	"log"

	"github.com/ductran999/dbkit"
)

func ConnectDB(pgConf dbkit.PostgreSQLConfig) (dbkit.Connection, error) {
	conn, err := dbkit.NewPostgreSQLConnection(pgConf)
	if err != nil {
		return nil, err
	}
	// Test the connection
	if err := conn.Ping(context.Background()); err != nil {
		conn.Close() //nolint
		return nil, err
	}

	log.Println("Successfully connected to PostgreSQL!")
	return conn, nil
}

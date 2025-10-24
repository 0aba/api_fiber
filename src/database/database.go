package database

import (
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var Psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type DB struct {
	Connection   *sqlx.DB
	QueryBuilder sq.StatementBuilderType
}

var dbInstance *DB

func InitDatabase(host, port, dbname, user, password, sslmode string) (*DB, error) {
	connectString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	connection, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// INFO! останутся константными
	connection.SetMaxOpenConns(25)
	connection.SetMaxIdleConns(5)

	if err := connection.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbInstance = &DB{
		Connection:   connection,
		QueryBuilder: Psql,
	}

	return dbInstance, nil
}

func GetDB() *DB {
	if dbInstance == nil {
		log.Panic("Database not initialized. Call InitDatabase first")
	}

	return dbInstance
}

func Close() error {
	if dbInstance != nil && dbInstance.Connection != nil {
		return dbInstance.Connection.Close()
	}

	return nil
}

package store

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	conn *sqlx.DB
}

type People struct {
	ID   int
	Name string
}

// NewStore creates new database connection
func NewStore(connString string) *Store {
	conn, err := sqlx.Connect("pgx", connString)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil { //Ругается, если ничего не делаю с err, но нужны ли они тут?
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"../../migrations/1_initial.up.sql",
		"postgres", driver)
	if err != nil { //и тут/////////
		panic(err)
	}
	m.Up()

	return &Store{
		conn: conn,
	}
}

func (s *Store) ListPeople() ([]People, error) {

	rows, err := conn.Query(context.Background(), `
	SELECT name, address
	FROM client
	`)
	if err != nil {
		// return fmt.Errorf("client query failed: %w", err)
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		return
	}
	defer rows.Close()

	return nil, nil
}

func (s *Store) GetPeopleByID(id string) (People, error) {
	return People{}, nil
}

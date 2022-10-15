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

	rows, err := s.conn.Query(context.Background(), `
	SELECT name, address
	FROM client
	`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		return nil, nil
	}
	defer rows.Close()
	people := make([]People, 0) //слайс со всеми людьми
	for rows.Next() {
		var (
			name    string
			address string
		)

		if err := rows.Scan(&name, &address); err != nil {
			fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
			return nil, nil
		}
		people = append(people, rows.Scan(&name, &address))

	}
	return people, err
}

func (s *Store) GetPeopleByID(id string) (People, error) {
	return People{}, nil
}

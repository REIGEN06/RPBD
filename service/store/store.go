package store

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	conn *sqlx.DB
}

type People struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// NewStore creates new database connection
func NewStore(connString string) *Store {
	conn, err := sqlx.Connect("pgx", connString)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"../../migrations/1_initial.up.sql",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Up()

	return &Store{
		conn: conn,
	}
}

func (s *Store) ListPeople(ctx context.Context) ([]People, error) {
	//rows, err := s.conn.QueryContext(ctx, `SELECT name, address FROM people`)
	//defer rows.Close()
	people := make([]People, 0) //слайс со всеми людьми
	err := s.conn.SelectContext(ctx, &people, `SELECT name, address FROM people`)
	if err != nil {
		return nil, fmt.Errorf("store ListPeople error: %w", err)
	}

	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
	// 	return nil, nil
	// }
	// for rows.Next() {
	// 	var (
	// 		name    string
	// 		address string
	// 	)

	// 	if err := rows.Scan(&name, &address); err != nil {
	// 		fmt.Fprintf(os.Stderr, "scan failed: %v\n", err)
	// 		return nil, nil
	// 	}
	// 	people = append(people, &name, &address)

	// }
	return people, err
}

func (s *Store) GetPeopleByID(ctx context.Context, id string) (People, error) {
	var pp People
	err := s.conn.SelectContext(ctx, pp, `SELECT name, address FROM people WHERE People.ID = id`)
	if err != nil {
		return nil, fmt.Errorf("store GetPeopleByID error: %w", err)
	}
	return pp, nil
}

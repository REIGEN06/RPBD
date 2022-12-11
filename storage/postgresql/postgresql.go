package postgresql

import (
	"buy-list/product"
	"log"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type BuyList struct {
	conn *sqlx.DB
}

func Connect(connString string) *BuyList {
	conn, err := sqlx.Connect("pgx", connString)
	if err != nil {
		panic(err)
	}
	////////Выдает ошибку panic:no scheme
	// driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	// if err != nil {
	// 	panic(err)
	// }
	// m, err := migrate.NewWithDatabaseInstance(
	// 	"../../migrations/1_initial.up.sql",
	// 	"postgres", driver)
	// if err != nil {
	// 	panic(err)
	// }
	// m.Up()

	return &BuyList{
		conn: conn,
	}
}

func (c *BuyList) AddInList(p *product.Product) error {
	_, err := c.conn.NamedExec(`INSERT INTO products (telegram_user_id, name, weight, inlist,  timerenable)
	VALUES (:telegram_user_id, :name, :weight, :inlist, :timerenable)`, p)
	////////Если попытаться добавить created_at и finished_at, то напишет что не может найти их в структуре Product
	if err != nil {
		log.Printf("AddInList fail: %s", err)
		return err
	}
	return err
}

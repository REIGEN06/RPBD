package postgresql

import (
	"buy-list/product"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Connection struct {
	conn *sqlx.DB
}

func Connect(connString string) *Connection {
	conn, err := sqlx.Connect("pgx", connString)
	if err != nil {
		panic(err)
	}

	driver, err := postgres.WithInstance(conn.DB, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		panic(err)
	}
	//не создает переменную err
	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			println("tables are already migrated!")
		} else {
			panic(err)
		}
	} else {
		println("successfully migrated!")
	}

	return &Connection{
		conn: conn,
	}
}

func (c *Connection) AddIn(p *product.Product) error {
	_, err := c.conn.NamedExec(`INSERT INTO products (telegram_user_id, telegram_chat_id, name, weight, inlist, infridge,  timerenable, created_at, finished_at)
	VALUES (:telegram_user_id, :telegram_chat_id, :name, :weight, :inlist, :infridge, :timerenable, :created_at, :finished_at)`, p)

	if err != nil {
		log.Printf("AddInfail: %s", err)
		return err
	}
	return err
}

func (c *Connection) AddInF(p *product.Product) error {
	_, err := c.conn.NamedExec(`INSERT INTO products (telegram_user_id, telegram_chat_id, name, weight, inlist, infridge,  timerenable, created_at)
	VALUES (:telegram_user_id, :telegram_chat_id, :name, :weight, :inlist, :infridge, :timerenable, :created_at)`, p)

	if err != nil {
		log.Printf("AddInF fail: %s", err)
		return err
	}
	return err
}

func (c *Connection) CreateUser(name string, user_id int64, chat_id int64) string {
	msg := "Теперь вы можете пользоваться ботом, "
	msg += name
	msg += "!😊"
	sqlStatement := `INSERT INTO users  (telegram_user_name, telegram_user_id, telegram_chat_id) VALUES ($1, $2, $3)`
	_, err := c.conn.Exec(sqlStatement, name, user_id, chat_id)

	if err != nil {
		if strings.Contains(err.Error(), "повторяющееся") {
			msg = name
			msg += ", Вы уже зарегистрированы!"
		}
		log.Printf("CreateUser fail: %s", err)
		return msg
	}

	return msg
}

func (c *Connection) GetStatus(user_id int64, chat_id int64) int {
	var status int
	c.conn.QueryRow(`SELECT user_status FROM users WHERE telegram_user_id = $1 AND telegram_chat_id = $2`, user_id, chat_id).Scan(&status)
	return status
}

func (c *Connection) GetList(user_id int64, chat_id int64) ([]product.Product, error) {
	products := make([]product.Product, 0)
	var (
		id          int64
		name        string
		weight      float64
		timerenable bool
	)

	rows, err := c.conn.Query(`SELECT id, name, weight, timerenable FROM products WHERE telegram_user_id = $1 AND telegram_chat_id = $2`, user_id, chat_id)
	if err != nil {
		return nil, fmt.Errorf("GetList error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&id, &name, &weight, &timerenable); err != nil {
			fmt.Fprintf(os.Stderr, "GetList scan failed: %v\n", err)
			return nil, nil
		}

		products = append(products, product.Product{Id: id, Name: name, Weight: weight, TimerEnable: timerenable})
	}
	return products, err
}
func (c *Connection) SetStatus(newStatus int, user_id int64, chat_id int64) {
	sqlStatement := `UPDATE users SET user_status = $1 WHERE telegram_user_id = $2 AND telegram_chat_id = $3`
	_, err := c.conn.Exec(sqlStatement, newStatus, user_id, chat_id)

	if err != nil {
		log.Println(err)
	}
}

func (c *Connection) SetFridge(user_id int64, chat_id int64, from time.Time, to time.Time, name string) error {
	sqlStatement := `UPDATE products SET infridge = $1, inlist = $2, created_at = $3, finished_at = $4, timerenable = $5 WHERE telegram_user_id = $6 AND telegram_chat_id = $7 AND name = $8`
	_, err := c.conn.Exec(sqlStatement, true, false, from, to, true, user_id, chat_id, name)

	if err != nil {
		log.Println(err)
	}

	return err
}

func (c *Connection) OpenProduct(user_id int64, chat_id int64, from time.Time, to time.Time, name string) error {
	//sqlStatement := `UPDATE products SET infridge = $1, inlist = $2, created_at = $3, finished_at = $4, timerenable = $5 WHERE telegram_user_id = $6 AND telegram_chat_id = $7 AND name = $8`
	_, err := c.conn.Exec(sqlStatement, true, false, from, to, true, user_id, chat_id, name)

	if err != nil {
		log.Println(err)
	}

	return err
}

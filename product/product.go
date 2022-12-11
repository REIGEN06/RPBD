package product

import (
	"time"
)

type Product struct {
	Id          int64   `db:"id"`
	User_id     int64   `db:"telegram_user_id"`
	Name        string  `db:"name"`
	Weight      float64 `db:"weight"`
	AlreadyUsed bool    `db:"alreadyused"`
	InList      bool    `db:"inlist"`
	InFridge    bool    `db:"infridge"`
	InTrash     bool    `db:"intrash"`
	created_at  string  `db:"created_at"`
	finished_at string  `db:"finished_at"`
	TimerEnable bool    `db:"timerenable"`
}

func (p *Product) CreateProduct(user_id int64, name string, weight float64, finished_at string) *Product {
	p.User_id = user_id
	p.Name = name
	p.Weight = weight
	p.InList = true
	timenow := time.Now().Format(time.RFC3339)
	p.created_at = timenow
	p.finished_at = finished_at
	p.TimerEnable = true
	return p
}

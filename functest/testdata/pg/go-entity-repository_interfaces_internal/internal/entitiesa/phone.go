package entitiesa

type Phone struct {
	UserID int    `db:"user_id"`
	Number string `db:"number"`
}

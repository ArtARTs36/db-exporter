package entities

type Country struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

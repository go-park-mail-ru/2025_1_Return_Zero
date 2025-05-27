package repository

type Label struct {
	ID      int64    `sql:"id"`
	Name    string   `sql:"name"`
	Members []string `sql:"members"`
}

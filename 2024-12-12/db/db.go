package db

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DB struct {
	users []User
}

func NewDatabase() *DB {
	return &DB{
		users: []Users{}
	}
}
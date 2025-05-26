package models

type User struct {
	UserID   int    `db:"user_id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	PwHash   string `db:"pw_hash"`
}

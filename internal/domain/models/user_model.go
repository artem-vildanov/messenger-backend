package models

type UserModel struct {
	Id           int    `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}

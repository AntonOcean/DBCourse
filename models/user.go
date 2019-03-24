package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	About string `json:"about"`
	Email string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname  string `json:"nickname"`
}

type UserUpdateInsert struct {
	About string `json:"about"`
	Email string `json:"email"`
	Fullname string `json:"fullname"`
}

func (up UserUpdateInsert) Validate() (err error) {
	if up.Fullname == "" || up.Email == "" {
		return errors.New("oops")
	}
	return
}

func (u User) Create(db *sql.DB) (err error) {
	if u.About == "" {
		_, err = db.Exec(
			"INSERT INTO forum_db.user(email, fullname, nickname) VALUES ($1, $2, $3)",
			u.Email, u.Fullname, u.Nickname,
		)
	} else {
		_, err = db.Exec(
			"INSERT INTO forum_db.user(about, email, fullname, nickname) VALUES ($1, $2, $3, $4)",
			u.About, u.Email, u.Fullname, u.Nickname,
		)
	}
	if err != nil {
		fmt.Println("user model Create 1 ",err)
		return
	}
	return
}

func (u User) GetDuplicate(db *sql.DB) (users []User, err error) {
	rows, err := db.Query(
		"SELECT u.about, u.email, u.fullname, u.nickname FROM forum_db.user u WHERE u.nickname=$1 OR u.email=$2",
		u.Nickname, u.Email,
		)
	defer rows.Close()
	for rows.Next() {
		   var user User
		   err = rows.Scan(&user.About, &user.Email, &user.Fullname, &user.Nickname)
		   if err != nil {
			   fmt.Println("user model GetDuplicate 1 ",err)
			   return nil, err
		   }
		   users = append(users, user)
	   }

	return
}

func (u *User) Get(db *sql.DB, nickname string) (err error) {
	row := db.QueryRow(
		"SELECT u.about, u.email, u.fullname, u.nickname FROM forum_db.user u WHERE u.nickname=$1",
		nickname,
		)
	err = row.Scan(&u.About, &u.Email, &u.Fullname, &u.Nickname)
	if err != nil {
		fmt.Println("user model Get 1 ",err)
		return err
	}
	return
}

func (u *User) Update(db *sql.DB) (rows int, err error) {
	rs, err := db.Exec(
		"UPDATE forum_db.user u SET about=$1, email=$2, fullname=$3 WHERE u.nickname=$4",
		u.About, u.Email, u.Fullname, u.Nickname,
	)
	if err != nil {
		fmt.Println("user model Update 1 ",err)
		return 0, err
	}
	row, err := rs.RowsAffected()
	if err != nil {
		fmt.Println("user model Update 2 ",err)
		return 0, err
	}
	rows = int(row)
	return
}
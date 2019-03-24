package models

import (
	"database/sql"
	"fmt"
)

type Error struct {
	Message string `json:"message"`
}

type Empty struct {
	Null string `json:"-"`
}


type InfoDB struct {
	Forum int64 `json:"forum"`
	Post int64 `json:"post"`
	Thread int64 `json:"thread"`
	User int64 `json:"user"`
}


func DropAll(db *sql.DB) {
	_, err := db.Exec(`
DELETE FROM forum_db.post; 
DELETE FROM forum_db.vote;
DELETE FROM forum_db.thread;
DELETE FROM forum_db.forum;
DELETE FROM forum_db.user;
`)
	if err != nil {
		fmt.Println("Fail delete service", err)
	}
}

func (i *InfoDB)Get(db *sql.DB) {

	var count int64
	rows, err := db.Query("SELECT COUNT(*) FROM forum_db.thread")
	defer rows.Close()
	if err != nil {
		fmt.Println("service Get 1 ",err)
		i.Thread = 0
	} else {

		for rows.Next() {
			err:= rows.Scan(&count)
			if err != nil {
				fmt.Println("service Get 2 ",err)
				i.Thread = 0
				return
			}
		}
		i.Thread = count

	}


	count = 0
	rows, err = db.Query("SELECT COUNT(*) FROM forum_db.post")
	defer rows.Close()
	if err != nil {
		fmt.Println("service Get 1 ",err)
		i.Post = 0
	} else {

		for rows.Next() {
			err:= rows.Scan(&count)
			if err != nil {
				fmt.Println("service Get 2 ",err)
				i.Post = 0
				return
			}
		}
		i.Post = count

	}


	count = 0
	rows, err = db.Query("SELECT COUNT(*) FROM forum_db.forum")
	defer rows.Close()
	if err != nil {
		fmt.Println("service Get 1 ",err)
		i.Forum = 0
		return
	} else {

		for rows.Next() {
			err:= rows.Scan(&count)
			if err != nil {
				fmt.Println("service Get 2 ",err)
				i.Forum = 0
				return
			}
		}
		i.Forum = count

	}



	count = 0
	rows, err = db.Query("SELECT COUNT(*) FROM forum_db.user")
	defer rows.Close()
	if err != nil {
		fmt.Println("service Get 1 ",err)
		i.User = 0
		return
	} else {

		for rows.Next() {
			err:= rows.Scan(&count)
			if err != nil {
				fmt.Println("service Get 2 ",err)
				i.User = 0
				return
			}
		}
		i.User = count

	}

}
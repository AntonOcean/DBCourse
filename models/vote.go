package models

import (
	"database/sql"
	"fmt"
)

type Vote struct {
	Nickname string `json:"nickname"`
	Voice int32 `json:"voice"`
}

func (v Vote) Create(db *sql.DB, thread_id int32) (err error) {
	//  user_id citext references forum_db.user(nickname) not null,
	//  thread_id integer references forum_db.thread(id) not null,
	//  value integer default 1,
	//  primary key (user_id, thread_id)
	//);
	_, err = db.Exec(
		"INSERT INTO forum_db.vote(user_id, thread_id, value) VALUES ($1, $2, $3)",
		v.Nickname, thread_id, v.Voice,
	)


	if err != nil {
		fmt.Println("vote model Create 1 ",err)
		_, err = db.Exec(
			"UPDATE forum_db.vote SET value = $1 WHERE user_id=$2 AND thread_id=$3",
			v.Voice, v.Nickname, thread_id,
		)
		if err != nil {
			fmt.Println("vote model Create 2 ",err)
		}
		return
	}

	return
}

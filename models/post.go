package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Post struct {
	Author string `json:"author"`
	Created string `json:"created"`
	Forum string `json:"forum"`
	Id int64 `json:"id"`
	IsEdited bool `json:"isEdited"`
	Message string `json:"message"`
	Parent int64 `json:"parent"`
	Thread int32 `json:"thread"`
}

type PostWithParent struct {
	Author string `json:"author"`
	Created string `json:"created"`
	Forum string `json:"forum"`
	Id int64 `json:"id"`
	IsEdited bool `json:"isEdited"`
	Message string `json:"message"`
	Parent int64 `json:"parent"`
	Thread int32 `json:"thread"`
}

type PostWithOutParent struct {
	Author string `json:"author"`
	Created string `json:"created"`
	Forum string `json:"forum"`
	Id int64 `json:"id"`
	IsEdited bool `json:"isEdited"`
	Message string `json:"message"`
	Thread int32 `json:"thread"`
}

type PostFull struct {
	Author User `json:"author"`
	Forum Forum `json:"forum"`
	Post Post `json:"post"`
	Thread Thread `json:"thread"`
}

type PostUpdate struct {
	Message string `json:"message"`
}


func (p *Post) CheckParentExists(db *sql.DB) (exist bool,err error) {
	var rows *sql.Rows
	//
	//if p.Parent == 0 {
	//	rows, err = db.Query("SELECT NOT EXISTS (SELECT * FROM forum_db.post p WHERE p.thread_id=$1 AND p.id is NULL)", // p.parent_id is NULL
	//		p.Thread)
	//} else {
	rows, err = db.Query("SELECT EXISTS (SELECT * FROM forum_db.post p WHERE p.thread_id=$1 AND p.id=$2)",
		p.Thread, p.Parent)
	//}

	defer rows.Close()
	if err != nil {
		fmt.Println("post model CheckParentExists 1 ",err)
		exist = false
		return
	}

	for rows.Next() {
		err = rows.Scan(&exist)
		if err != nil {
			fmt.Println("post model CheckParentExists 2 ",err)
			exist = false
			return
		}
	}
	return
}

func (p Post) CreatePost(db *sql.DB, date time.Time) (err error) {
	if p.Created == "" {
		if p.Parent == 0 {
			_, err = db.Exec(
				"INSERT INTO forum_db.post(author_id, forum_id, message, thread_id, isEdited, created) VALUES ($1, $2, $3, $4, $5, $6)",
				p.Author, p.Forum, p.Message, p.Thread, p.IsEdited, date,
			)
		} else {
			_, err = db.Exec(
				"INSERT INTO forum_db.post(author_id, forum_id, message, parent_id, thread_id, isEdited, created) VALUES ($1, $2, $3, $4, $5, $6, $7)",
				p.Author, p.Forum, p.Message, p.Parent, p.Thread, p.IsEdited, date,
			)
			//fmt.Println("hi ", p.Created, p.Parent.Int64)
		}
	} else {
		if p.Parent == 0 {
			_, err = db.Exec(
				"INSERT INTO forum_db.post(author_id, forum_id, message, thread_id, created, isEdited) VALUES ($1, $2, $3, $4, $5, $6)",
				p.Author, p.Forum, p.Message, p.Thread, p.Created, p.IsEdited,
			)
		} else {
			_, err = db.Exec(
				"INSERT INTO forum_db.post(author_id, forum_id, message, parent_id, thread_id, created, isEdited) VALUES ($1, $2, $3, $4, $5, $6, $7)",
				p.Author, p.Forum, p.Message, p.Parent, p.Thread, p.Created, p.IsEdited,
			)
		}
	}

	if err != nil {
		fmt.Println("post model CreatePost 1 ",err)
		return
	}

	return
}

func (p *Post) Get(db *sql.DB) (err error) {
	var row *sql.Row

	// titile and others
	if p.Parent == 0 {
		row = db.QueryRow(
			"SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id FROM forum_db.post p " +
				"WHERE p.thread_id=$1 AND p.parent_id is NULL AND p.author_id=$2 AND p.message=$3",
			p.Thread, p.Author, p.Message,
		)
	} else {
		row = db.QueryRow(
			"SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id FROM forum_db.post p " +
				"WHERE p.thread_id=$1 AND p.parent_id=$2 AND p.author_id=$3 AND p.message=$4",
			p.Thread, p.Parent, p.Author, p.Message,
		)
		//fmt.Println(p.Thread, p.Parent)
	}

	var temp sql.NullInt64

	err = row.Scan(&p.Id, &p.Author, &p.Created, &p.Forum, &p.IsEdited, &p.Message, &temp, &p.Thread)

	if err != nil {
		fmt.Println("post model Get 1 ",err)
		p.Parent = 0
		return
	}

	if !temp.Valid {
		p.Parent = 0
	} else {
		p.Parent = temp.Int64
	}

	return
}

func (p *Post) GetById(db *sql.DB, id int) (err error) {
	var row *sql.Row


	row = db.QueryRow(
		"SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id FROM forum_db.post p "+
			"WHERE p.id=$1",
		id,
	)

	var temp sql.NullInt64

	err = row.Scan(&p.Id, &p.Author, &p.Created, &p.Forum, &p.IsEdited, &p.Message, &temp, &p.Thread)

	if err != nil {
		fmt.Println("post model GetById 1 ", err)
		return
	}

	if !temp.Valid {
		p.Parent = 0
	} else {
		p.Parent = temp.Int64
	}


	return
}

func (p *Post) Update(db *sql.DB) (err error) {

	_, err = db.Exec(
		`
UPDATE forum_db.post SET message = $1, isEdited=true WHERE id = $2
`,
		p.Message, p.Id,
	)

	return
}
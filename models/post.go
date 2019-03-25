package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Post struct {
	AuthorID  string        `json:"author"`
	Created   NullTime      `json:"created"`
	ForumSlug string        `json:"forum"`
	Id        int64         `json:"id"`
	IsEdited  bool          `json:"isEdited"`
	Message   string        `json:"message"`
	ParentID  JsonNullInt64 `json:"parent"`
	ThreadID  int32         `json:"thread"`
}

type PostFull struct {
	Author User   `json:"author"`
	Forum  Forum  `json:"forum"`
	Post   Post   `json:"post"`
	Thread Thread `json:"thread"`
}

type PostUpdate struct {
	Message string `json:"message"`
}

func (p *Post) CheckParentExists(db *sql.DB) (exist bool) {
	var rows *sql.Rows
	//
	//if p.ParentID == 0 {
	//	rows, err = db.Query("SELECT NOT EXISTS (SELECT * FROM forum_db.post p WHERE p.thread_id=$1 AND p.id is NULL)", // p.parent_id is NULL
	//		p.ThreadID)
	//} else {
	//fmt.Println(p.ThreadID)
	rows, err := db.Query("SELECT EXISTS (SELECT * FROM forum_db.post p WHERE p.thread_id=$1 AND p.id=$2)",
		p.ThreadID, p.ParentID.Int64)
	//}

	defer rows.Close()

	if err != nil {
		log.Fatal("post model CheckParentExists 1 ", err)
		//fmt.Println("post model CheckParentExists 1 ",err)
		//exist = false
		//return
	}

	for rows.Next() {
		err = rows.Scan(&exist)
		if err != nil {
			log.Fatal("post model CheckParentExists 2 ", err)
			//fmt.Println("post model CheckParentExists 2 ",err)
			//exist = false
			//return
		}
	}
	return
}

func (p *Post) CheckAuthorExists(db *sql.DB) (exist bool) {
	var rows *sql.Rows

	rows, err := db.Query("SELECT EXISTS (SELECT * FROM forum_db.user u WHERE u.nickname=$1)",
		p.AuthorID)

	defer rows.Close()

	if err != nil {
		log.Fatal("post model CheckAuthorExists 1 ", err)

	}

	for rows.Next() {
		err = rows.Scan(&exist)
		if err != nil {
			log.Fatal("post model CheckAuthorExists 2 ", err)
		}
	}
	return
}

func CreatePostsBulk(db *sql.DB, data []Post) (posts []Post) {

	currentTime := time.Now()
	countPosts := len(data)

	txn, err := db.Begin()
	if err != nil {
		log.Fatal("post model CreatePostsBulk 1 ", err)
	}

	stmt, err := txn.Prepare(`COPY forum_db.post (author_id, created, forum_id, isEdited, message, parent_id, thread_id) FROM STDIN`)
	if err != nil {
		log.Fatal("post model CreatePostsBulk 2 ", err)
	}

	for _, post := range data {
		if !post.Created.Valid {
			post.Created.Time = currentTime
			post.Created.Valid = true
			_, err = stmt.Exec(post.AuthorID, post.Created, post.ForumSlug, post.IsEdited, post.Message, post.ParentID, post.ThreadID)
			if err != nil {
				log.Fatal("post model CreatePostsBulk 3.1 ", err)
			}
		} else {
			_, err = stmt.Exec(post.AuthorID, post.Created, post.ForumSlug, post.IsEdited, post.Message, post.ParentID, post.ThreadID)
			if err != nil {
				log.Fatal("post model CreatePostsBulk 3 ", err)
			}
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		log.Fatal("post model CreatePostsBulk 4 ", err)
	}

	err = stmt.Close()
	if err != nil {
		log.Fatal("post model CreatePostsBulk 5 ", err)
	}

	err = txn.Commit()
	if err != nil {
		log.Fatal("post model CreatePostsBulk 6 ", err)
	}

	rows, err := db.Query(`
SELECT * FROM (SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id
 			FROM forum_db.post p
		  ORDER BY p.id DESC
		  LIMIT $1) t
ORDER BY t.id;`, countPosts)

	defer rows.Close()

	if err != nil {
		log.Fatal("post model CreatePostsBulk 2 ", err)

	}

	for rows.Next() {
		var post Post

		err = rows.Scan(&post.Id, &post.AuthorID, &post.Created, &post.ForumSlug, &post.IsEdited, &post.Message, &post.ParentID, &post.ThreadID)
		if err != nil {
			log.Fatal("post model CreatePostsBulk 3 ", err)
		}
		posts = append(posts, post)
	}

	return
}

func (p *Post) Get(db *sql.DB) (err error) {
	var row *sql.Row

	// titile and others
	if !p.ParentID.Valid || p.ParentID.Int64 == 0 {
		row = db.QueryRow(
			"SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id FROM forum_db.post p "+
				"WHERE p.thread_id=$1 AND p.parent_id is NULL AND p.author_id=$2 AND p.message=$3",
			p.ThreadID, p.AuthorID, p.Message,
		)
	} else {
		row = db.QueryRow(
			"SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id FROM forum_db.post p "+
				"WHERE p.thread_id=$1 AND p.parent_id=$2 AND p.author_id=$3 AND p.message=$4",
			p.ThreadID, p.ParentID, p.AuthorID, p.Message,
		)
		//fmt.Println(p.ThreadID, p.ParentID)
	}

	//var temp sql.NullInt64

	err = row.Scan(&p.Id, &p.AuthorID, &p.Created, &p.ForumSlug, &p.IsEdited, &p.Message, &p.ParentID, &p.ThreadID)

	if err != nil {
		fmt.Println("post model Get 1 ", err)
		//p.ParentID = 0
		return
	}

	//if !temp.Valid {
	//	p.ParentID = 0
	//} else {
	//	p.ParentID = temp.Int64
	//}

	return
}

func (p *Post) GetById(db *sql.DB, id int) (err error) {
	var row *sql.Row

	row = db.QueryRow(
		"SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id FROM forum_db.post p "+
			"WHERE p.id=$1",
		id,
	)

	//var temp sql.NullInt64

	err = row.Scan(&p.Id, &p.AuthorID, &p.Created, &p.ForumSlug, &p.IsEdited, &p.Message, &p.ParentID, &p.ThreadID)

	if err != nil {
		fmt.Println("post model GetById 1 ", err)
		return
	}

	//if !temp.Valid {
	//	p.ParentID = 0
	//} else {
	//	p.ParentID = temp.Int64
	//}

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

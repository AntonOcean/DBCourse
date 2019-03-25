package models

import (
	"database/sql"
	"fmt"
)

type Thread struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Id      int32  `json:"id"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Votes   int32  `json:"votes"`
}

type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

func (t Thread) Create(db *sql.DB) (err error) {

	if t.Created == "" {
		if t.Slug == "" {
			_, err = db.Exec(
				"INSERT INTO forum_db.thread(author_id, forum_id, message, title) VALUES ($1, $2, $3, $4)",
				t.Author, t.Forum, t.Message, t.Title,
			)
		} else {
			_, err = db.Exec(
				"INSERT INTO forum_db.thread(author_id, forum_id, message, slug, title) VALUES ($1, $2, $3, $4, $5)",
				t.Author, t.Forum, t.Message, t.Slug, t.Title,
			)
		}
	} else {
		if t.Slug == "" {
			_, err = db.Exec(
				"INSERT INTO forum_db.thread(author_id, forum_id, message, title, created) VALUES ($1, $2, $3, $4, $5)",
				t.Author, t.Forum, t.Message, t.Title, t.Created,
			)
		} else {
			_, err = db.Exec(
				"INSERT INTO forum_db.thread(author_id, forum_id, message, slug, title, created) VALUES ($1, $2, $3, $4, $5, $6)",
				t.Author, t.Forum, t.Message, t.Slug, t.Title, t.Created,
			)
		}
	}

	if err != nil {
		fmt.Println("thread model Create 1 ", err)
		return
	}
	return
}

func (t *Thread) Get(db *sql.DB, slug string, id int32, title string) (err error) {
	var row *sql.Row
	if slug != "" {
		row = db.QueryRow(
			"SELECT t.id, t.author_id, t.created, t.forum_id, t.message, t.slug, t.title FROM forum_db.thread t WHERE t.slug=$1",
			slug,
		)
	} else if id != 0 {
		row = db.QueryRow(
			"SELECT t.id, t.author_id, t.created, t.forum_id, t.message, t.slug, t.title FROM forum_db.thread t WHERE t.id=$1",
			id,
		)
	} else {
		row = db.QueryRow(
			"SELECT t.id, t.author_id, t.created, t.forum_id, t.message, t.slug, t.title FROM forum_db.thread t WHERE t.title=$1",
			title,
		)
	}

	err = row.Scan(&t.Id, &t.Author, &t.Created, &t.Forum, &t.Message, &t.Slug, &t.Title)
	if err != nil {
		fmt.Println("thread model Get 1 ", err)
		return
	}
	return
}

func (t *Thread) Update(db *sql.DB) (err error) {
	_, err = db.Exec(
		"UPDATE forum_db.thread SET message=$1, title=$2 WHERE id=$3",
		t.Message, t.Title, t.Id,
	)
	if err != nil {
		fmt.Println("thread model Update 1 ", err)
		return
	}

	return
}

func (t *Thread) SetVotesCount(db *sql.DB) {
	var sum int32

	rows, err := db.Query("SELECT SUM(v.value) FROM forum_db.vote v WHERE v.thread_id=$1",
		t.Id)
	defer rows.Close()
	if err != nil {
		fmt.Println("thread model SetVotesCount 1 ", err)
		t.Votes = 0
		return
	}

	for rows.Next() {
		err := rows.Scan(&sum)
		if err != nil {
			fmt.Println("thread model SetVotesCount 2 ", err)
			t.Votes = 0
			return
		}
	}

	t.Votes = sum
}

func (t Thread) GetLastPost(db *sql.DB) (post Post, err error) {
	row := db.QueryRow(
		"SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id FROM forum_db.post p WHERE p.thread_id=$1 ORDER BY p.id DESC LIMIT 1",
		t.Id,
	)

	err = row.Scan(&post.Id, &post.AuthorID, &post.Created, &post.ForumSlug, &post.IsEdited, &post.Message, &post.ParentID, &post.ThreadID)
	if err != nil {
		fmt.Println("thread model GetLastPost 1 ", err)
	}
	return
}

func (t Thread) GetListPost(db *sql.DB, since int64, sort string, desc string, limit int) (posts []Post, err error) {
	var rows *sql.Rows

	switch sort {
	case "", "flat":

		if desc == "true" {
			if since == 0 {
				rows, err = db.Query(`
SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id FROM forum_db.post p
WHERE p.thread_id = $1 
ORDER BY p.created DESC, p.id DESC
LIMIT $2
`,
					t.Id, limit)
			} else {
				rows, err = db.Query(`
SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id FROM forum_db.post p
WHERE p.thread_id = $1 AND p.id < $2
ORDER BY p.created DESC, p.id DESC
LIMIT $3
`,
					t.Id, since, limit)
			}
		} else {
			if since == 0 {
				rows, err = db.Query(`
SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id FROM forum_db.post p
WHERE p.thread_id = $1 
ORDER BY p.created, p.id
LIMIT $2
`,
					t.Id, limit)
			} else {
				rows, err = db.Query(`
SELECT p.id, p.author_id, p.created, p.forum_id, p.isEdited, p.message, p.parent_id, p.thread_id FROM forum_db.post p
WHERE p.thread_id = $1 AND p.id > $2
ORDER BY p.created, p.id
LIMIT $3
`,
					t.Id, since, limit)
			}
		}

	case "tree":

		if desc == "true" {
			if since == 0 {
				rows, err = db.Query(`
SELECT child.id, child.author_id, child.created, child.forum_id, child.isEdited, child.message, child.parent_id, child.thread_id FROM forum_db.post parent
JOIN forum_db.post child ON child.parent_id = parent.id OR (child.id = parent.id AND child.parent_id is NULL)
WHERE parent.thread_id = $1
ORDER BY child.created DESC, parent.id DESC, child.id DESC 
LIMIT $2
`,
					t.Id, limit)
			} else {
				rows, err = db.Query(`
SELECT child.id, child.author_id, child.created, child.forum_id, child.isEdited, child.message, child.parent_id, child.thread_id FROM forum_db.post parent
JOIN forum_db.post child ON child.parent_id = parent.id OR (child.id = parent.id AND child.parent_id is NULL)
WHERE parent.thread_id = $1 AND parent.id < $2
ORDER BY parent.created DESC, parent.id DESC, child.id DESC
LIMIT $3
`,
					t.Id, since, limit)
			}
		} else {
			if since == 0 {
				rows, err = db.Query(`
SELECT child.id, child.author_id, child.created, child.forum_id, child.isEdited, child.message, child.parent_id, child.thread_id FROM forum_db.post parent
JOIN forum_db.post child ON child.parent_id = parent.id OR (child.id = parent.id AND child.parent_id is NULL)
WHERE parent.thread_id = $1
ORDER BY parent.created, parent.id, child.id
LIMIT $2
`,
					t.Id, limit)
			} else {
				rows, err = db.Query(`
SELECT child.id, child.author_id, child.created, child.forum_id, child.isEdited, child.message, child.parent_id, child.thread_id FROM forum_db.post parent
JOIN forum_db.post child ON child.parent_id = parent.id OR (child.id = parent.id AND child.parent_id is NULL)
WHERE parent.thread_id = $1 AND parent.id > $2
ORDER BY parent.created, parent.id, child.id
LIMIT $3
`,
					t.Id, since, limit)
			}
		}

	case "parent_tree":

		if desc == "true" {
			if since == 0 {
				rows, err = db.Query(`
SELECT child.id, child.author_id, child.created, child.forum_id, child.isEdited, child.message, child.parent_id, child.thread_id FROM forum_db.post parent
JOIN forum_db.post child ON child.parent_id = parent.id OR (child.id = parent.id AND child.parent_id is NULL)
WHERE parent.thread_id = $1 AND parent.id IN (
    SELECT parent.id FROM forum_db.post parent
WHERE parent.thread_id = $1
ORDER BY parent.created DESC, parent.id DESC
LIMIT $2
  )
ORDER BY parent.created DESC, parent.id DESC, child.id
`,
					t.Id, limit)
			} else {
				rows, err = db.Query(`
SELECT child.id, child.author_id, child.created, child.forum_id, child.isEdited, child.message, child.parent_id, child.thread_id FROM forum_db.post parent
JOIN forum_db.post child ON child.parent_id = parent.id OR (child.id = parent.id AND child.parent_id is NULL)
WHERE parent.thread_id = $1 AND parent.id IN (
    SELECT parent.id FROM forum_db.post parent
WHERE parent.thread_id = $1 AND parent.id < $2
ORDER BY parent.created DESC, parent.id DESC
LIMIT $3
  )
ORDER BY parent.created DESC, parent.id DESC, child.id
`,
					t.Id, since, limit)
			}
		} else {
			if since == 0 {
				rows, err = db.Query(`
SELECT child.id, child.author_id, child.created, child.forum_id, child.isEdited, child.message, child.parent_id, child.thread_id FROM forum_db.post parent
JOIN forum_db.post child ON child.parent_id = parent.id OR (child.id = parent.id AND child.parent_id is NULL)
WHERE parent.thread_id = $1 AND parent.id IN (
    SELECT parent.id FROM forum_db.post parent
WHERE parent.thread_id = $1
ORDER BY parent.created, parent.id
LIMIT $2
  )
ORDER BY parent.created, parent.id, child.id
`,
					t.Id, limit)
			} else {
				rows, err = db.Query(`
SELECT child.id, child.author_id, child.created, child.forum_id, child.isEdited, child.message, child.parent_id, child.thread_id FROM forum_db.post parent
JOIN forum_db.post child ON child.parent_id = parent.id OR (child.id = parent.id AND child.parent_id is NULL)
WHERE parent.thread_id = $1 AND parent.id IN (
    SELECT parent.id FROM forum_db.post parent
WHERE parent.thread_id = $1 AND parent.id > $2
ORDER BY parent.created, parent.id
LIMIT $3
  )
ORDER BY parent.created, parent.id, child.id
`,
					t.Id, since, limit)
			}
		}

	}

	defer rows.Close()

	if err != nil {
		fmt.Println("forum model GetListPost 1 ", err)
		return
	}

	for rows.Next() {
		var post Post
		//var temp sql.NullInt64
		err = rows.Scan(&post.Id, &post.AuthorID, &post.Created, &post.ForumSlug, &post.IsEdited, &post.Message, &post.ParentID, &post.ThreadID)
		if err != nil {
			fmt.Println("forum model GetListPost 2 ", err)
		}

		//if !temp.Valid {
		//	post.ParentID = 0
		//} else {
		//	post.ParentID = temp.Int64
		//}

		posts = append(posts, post)
	}

	return
}

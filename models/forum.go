package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Forum struct {
	Posts int64 `json:"posts"`
	Slug string `json:"slug"`
	Threads int32 `json:"threads"`
	Title string `json:"title"`
	User string `json:"user"`
}

func (f Forum) Validate() (err error) {
	if f.Slug == "" || f.Title == "" || f.User == "" {
		return errors.New("oops E")
	}
	return
}

func (f Forum) Create(db *sql.DB) (err error) {
	_, err = db.Exec(
		"INSERT INTO forum_db.forum(slug, author_id, title) VALUES ($1, $2, $3)",
		f.Slug, f.User, f.Title,
	)
	if err != nil {
		fmt.Println("forum model Create 1 ",err)
		return
	}
	return
}

func (f *Forum) Get(db *sql.DB, slug string) (err error) {
	row := db.QueryRow(
		"SELECT f.slug, f.title, f.author_id FROM forum_db.forum f WHERE f.slug=$1",
		slug,
	)

	err = row.Scan(&f.Slug, &f.Title, &f.User)
	if err != nil {
		fmt.Println("forum model Get 1 ",err)
		return
	}

	return
}

func (f *Forum) SetPostsCount(db *sql.DB) {
	var count int64

	rows, err := db.Query("SELECT COUNT(*) FROM forum_db.post p WHERE p.forum_id=$1",
		f.Slug)
	defer rows.Close()
	if err != nil {
		fmt.Println("forum model SetPostsCount 1 ",err)
		f.Posts = 0
		return
	}

	for rows.Next() {
		err:= rows.Scan(&count)
		if err != nil {
			fmt.Println("forum model SetPostsCount 2 ",err)
			f.Posts = 0
			return
		}
	}
	f.Posts = count
}

func (f *Forum) SetThreadsCount(db *sql.DB) {
	var count int32

	rows, err := db.Query("SELECT COUNT(*) FROM forum_db.thread p WHERE p.forum_id=$1",
		f.Slug)
	defer rows.Close()
	if err != nil {
		fmt.Println("forum model SetThreadsCount 1 ",err)
		f.Threads = 0
		return
	}

	for rows.Next() {
		err:= rows.Scan(&count)
		if err != nil {
			fmt.Println("forum model SetThreadsCount 2 ",err)
			f.Threads = 0
			return
		}
	}
	f.Threads = count
}

func (f *Forum) GetThreadList(db *sql.DB, limit int, since string, desc string) (threads []Thread, err error){

	var rows *sql.Rows
	if desc == "true" {
		if since == "" {
			since = "9999-12-11 23:59:59.997"
		}
		rows, err = db.Query("" +
			"SELECT t.author_id, t.created, t.forum_id, t.id, t.message, t.slug, t.title FROM forum_db.thread t WHERE t.forum_id=$1 " +
			"AND t.created <= $2 " +
			"ORDER BY t.created DESC " +
			"LIMIT $3 ",
			f.Slug, since, limit)
	} else if desc == "false" {
		if since == "" {
			since = "1900-01-01T00:00:00.000Z"
			//since = "9999-12-11 23:59:59.997"
		}
		rows, err = db.Query("" +
			"SELECT t.author_id, t.created, t.forum_id, t.id, t.message, t.slug, t.title FROM forum_db.thread t WHERE t.forum_id=$1 " +
			"AND t.created >= $2 " +
			"ORDER BY t.created " +
			"LIMIT $3 ",
			f.Slug, since, limit)
	} else {
		if since == "" {
			since = "1900-01-01T00:00:00.000Z"
			//since = "9999-12-11 23:59:59.997"
		}
		rows, err = db.Query("" +
			"SELECT t.author_id, t.created, t.forum_id, t.id, t.message, t.slug, t.title FROM forum_db.thread t WHERE t.forum_id=$1 " +
			"AND t.created >= $2 " +
			"ORDER BY t.created " +
			"LIMIT $3 ",
			f.Slug, since, limit)
	}

	defer rows.Close()

	if err != nil {
		fmt.Println("forum model GetThreadList 1 ",err)
		return
	}

	for rows.Next() {
	   var thread Thread
	   err = rows.Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title)
	   if err != nil {
		  fmt.Println("forum model GetThreadList 2 ",err)
	   }
	   thread.SetVotesCount(db)
	   threads = append(threads, thread)
    }
	return
}



func (f *Forum) GetUserList(db *sql.DB, limit int, since string, desc string) (users []User, err error){
	//create table forum_db.user (
	//  nickname citext primary key,
	//  fullname varchar(500) not null,
	//  email citext unique not null,
	//  about text
	//);

	// то что меньше, sice  если desc

	var rows *sql.Rows
	if desc == "true" {
		if since == "" {
			rows, err = db.Query(`
SELECT * FROM (SELECT u.nickname, u.fullname, u.email, u.about FROM forum_db.user u
JOIN forum_db.thread t ON t.forum_id=$1 AND t.author_id=u.nickname
GROUP BY u.nickname
UNION
SELECT u.nickname, u.fullname, u.email, u.about FROM forum_db.user u
JOIN forum_db.post p ON p.forum_id=$1 AND p.author_id=u.nickname
GROUP BY u.nickname) t
ORDER BY lower(t.nickname COLLATE "C") DESC
LIMIT $2
`,
				f.Slug, limit)
		} else {
			rows, err = db.Query(`
SELECT * FROM (SELECT u.nickname, u.fullname, u.email, u.about FROM forum_db.user u
JOIN forum_db.thread t ON t.forum_id=$1 AND t.author_id=u.nickname
GROUP BY u.nickname
UNION
SELECT u.nickname, u.fullname, u.email, u.about FROM forum_db.user u
JOIN forum_db.post p ON p.forum_id=$1 AND p.author_id=u.nickname
GROUP BY u.nickname) t
WHERE t.nickname COLLATE "C" < $2 COLLATE "C"
ORDER BY lower(t.nickname COLLATE "C") DESC
LIMIT $3
`,
				f.Slug, since, limit)
		}
	} else if desc == "false" {
		if since == "" {
			rows, err = db.Query(`
SELECT * FROM (SELECT u.nickname, u.fullname, u.email, u.about FROM forum_db.user u
JOIN forum_db.thread t ON t.forum_id=$1 AND t.author_id=u.nickname
GROUP BY u.nickname
UNION
SELECT u.nickname, u.fullname, u.email, u.about FROM forum_db.user u
JOIN forum_db.post p ON p.forum_id=$1 AND p.author_id=u.nickname
GROUP BY u.nickname) t
ORDER BY lower(t.nickname COLLATE "C")
LIMIT $2
`,
				f.Slug, limit)
		} else {
			rows, err = db.Query(`
SELECT * FROM (SELECT u.nickname, u.fullname, u.email, u.about FROM forum_db.user u
JOIN forum_db.thread t ON t.forum_id=$1 AND t.author_id=u.nickname
GROUP BY u.nickname
UNION
SELECT u.nickname, u.fullname, u.email, u.about FROM forum_db.user u
JOIN forum_db.post p ON p.forum_id=$1 AND p.author_id=u.nickname
GROUP BY u.nickname) t
WHERE t.nickname COLLATE "C" > $2 COLLATE "C"
ORDER BY lower(t.nickname COLLATE "C")
LIMIT $3
`,
				f.Slug, since, limit)
		}
	}

	defer rows.Close()

	if err != nil {
		fmt.Println("forum model GetThreadList 1 ",err)
		return
	}

	for rows.Next() {
		var user User
		err = rows.Scan(&user.Nickname, &user.Fullname, &user.Email, &user.About)
		if err != nil {
			fmt.Println("forum model GetThreadList 2 ",err)
		}
		users = append(users, user)
	}
	return
}
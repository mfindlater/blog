package blog

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type BlogContext struct {
	Db sql.DB
}

func (b *BlogContext) Connect() {

	db, err := sql.Open("sqlite3", "blog.db")

	if err != nil {
		log.Fatal(err)
	}

	b.Db = *db
}

func (b *BlogContext) GetPosts() []Post {

	posts := make([]Post, 0)
	rows, err := b.Db.Query("select title,body,created,updated,status from BlogPost")

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		post := Post{}
		var posted string
		var updated string
		err := rows.Scan(&post.Title, &post.Body, &posted, &updated, &post.Status)

		if err != nil {
			log.Fatal(err)
		}

		post.Posted = posted
		post.Updated = updated

		posts = append([]Post{post}, posts...)
	}

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	return posts
}

func (b *BlogContext) SavePost(post Post) {

	_, err := b.Db.Exec("INSERT INTO BlogPost(title,body,created,updated,status) VALUES(?,?,?,?,?)", post.Title, post.Body, post.Posted, post.Updated, post.Status)

	if err != nil {
		log.Fatal(err)
	}
}

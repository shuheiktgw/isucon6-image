package main

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

type Entry struct {
	ID          int
	AuthorID    int
	Keyword     string
	Description string
	Link        sql.NullString
	UpdatedAt   time.Time
	CreatedAt   time.Time

	Html  string
	Stars []*Star
}

type Star struct {
	ID        int       `json:"id"`
	Keyword   string    `json:"keyword"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	host := os.Getenv("ISUDA_DB_HOST")
	if host == "" {
		host = "localhost"
	}
	portstr := os.Getenv("ISUDA_DB_PORT")
	if portstr == "" {
		portstr = "3306"
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		log.Fatalf("Failed to read DB port number from an environment variable ISUDA_DB_PORT.\nError: %s", err.Error())
	}
	user := "isucon"
	password := "isucon"
	dbname := os.Getenv("ISUDA_DB_NAME")
	if dbname == "" {
		dbname = "isuda"
	}

	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?loc=Local&parseTime=true",
		user, password, host, port, dbname,
	))
	if err != nil {
		log.Fatalf("Failed to connect to DB: %s.", err.Error())
	}
	db.Exec("SET SESSION sql_mode='TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY'")
	db.Exec("SET NAMES utf8mb4")

	rows, err := db.Query(`
		SELECT id, author_id, keyword, description, updated_at, created_at FROM entry
	`)

	if err != nil {
		panic(err)
	}

	entries := make([]*Entry, 0, 500)
	for rows.Next() {
		e := Entry{}
		err := rows.Scan(&e.ID, &e.AuthorID, &e.Keyword, &e.Description, &e.UpdatedAt, &e.CreatedAt)
		if err != nil {
			panic(err)
		}
		entries = append(entries, &e)
	}
	rows.Close()
	for _, entry := range entries {
		u, err := url.Parse("http://172.28.128.7/keyword/" + (&url.URL{Path: entry.Keyword}).String())
		if err != nil {
			panic(err)
		}
		link := fmt.Sprintf("<a href=\"%s\">%s</a>", u, html.EscapeString(entry.Keyword))
		_, err = db.Exec("UPDATE entry SET link = ? WHERE id = ?", link, entry.ID)
		if err != nil {
			panic(err)
		}
	}



}
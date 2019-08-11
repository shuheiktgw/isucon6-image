package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"html"
	"log"
	"os"
	"strconv"
	"strings"
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
		SELECT keyword, link, description FROM entry ORDER BY CHARACTER_LENGTH(keyword) DESC
	`)

	if err != nil {
		panic(err)
	}

	keywordEntries := make([]*Entry, 0, 500)
	for rows.Next() {
		e := Entry{}
		err := rows.Scan(&e.Keyword, &e.Link, &e.Description)
		if err != nil {
			panic(err)
		}
		keywordEntries = append(keywordEntries, &e)
	}
	rows.Close()

	keywords := make([][]string, 0, 500)
	for _, entry := range keywordEntries {
		keywords = append(keywords, []string{entry.Keyword, entry.Link.String})
	}

	args := make([]string, 0, 500)
	sha2link := make(map[string]string)
	for _, keyword := range keywords {
		hash := "isuda_" + fmt.Sprintf("%x", sha1.Sum([]byte(keyword[0])))
		sha2link[hash] = keyword[1]
		args = append(args, keyword[0], hash)
	}

	replacer := strings.NewReplacer(args...)

	for _, ke := range keywordEntries {
		content := html.EscapeString(replacer.Replace(ke.Description))

		for sha, link := range sha2link {
			content = strings.Replace(content, sha, link, -1)
		}

		content = strings.Replace(content, "\n", "<br />\n", -1)

		_, err := db.Exec(`
		INSERT INTO cached_content (keyword, content)
		VALUES (?, ?)
		ON DUPLICATE KEY UPDATE
		keyword = ?, content = ?
	`, ke.Keyword, content, ke.Keyword, content)
		if err != nil {
			panic(err)
		}
	}

}
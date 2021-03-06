package main

import (
	"context"
	"database/sql"
	"time"
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

type CachedContent struct {
	Keyword string
	Content string
}

type User struct {
	ID        int
	Name      string
	Salt      string
	Password  string
	CreatedAt time.Time
}

type Star struct {
	ID        int       `json:"id"`
	Keyword   string    `json:"keyword"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

type EntryWithCtx struct {
	Context context.Context
	Entry   Entry
}

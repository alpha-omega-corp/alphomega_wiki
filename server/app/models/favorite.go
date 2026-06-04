package models

import (
	"time"

	"github.com/uptrace/bun"
)

// Favorite links a user (by their AOAuth subject) to a page they've starred.
// The (user_sub, page_id) pair is unique so a page is favorited at most once
// per user.
type Favorite struct {
	bun.BaseModel `bun:"table:favorites,alias:f"`

	Id        int64     `bun:",pk,autoincrement" json:"id"`
	UserSub   string    `bun:"user_sub,notnull,unique:uq_user_page" json:"user_sub"`
	PageID    int64     `bun:"page_id,notnull,unique:uq_user_page" json:"page_id"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}
package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Page struct {
	bun.BaseModel `bun:"table:pages,alias:p"`

	Id        int64     `bun:",pk,autoincrement" json:"id"`
	Slug      string    `bun:"slug,notnull" json:"slug"`
	Title     string    `bun:"title,notnull" json:"title"`
	Content   string    `bun:"content" json:"content"`
	Excerpt   string    `bun:"excerpt" json:"excerpt,omitempty"`
	Tags      []string  `bun:"tags,array" json:"tags"`
	AuthorID  int64     `bun:"author_id" json:"author_id"`
	Published bool      `bun:"published,notnull,default:false" json:"published"`
	ParentID  *int64    `bun:"parent_id" json:"parent_id,omitempty"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	Pages []*Page `bun:"rel:has-many,join:id=parent_id" json:"pages,omitempty"`
}

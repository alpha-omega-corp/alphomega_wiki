package utils

import (
	"context"
	"fmt"
	"reflect"

	"github.com/uptrace/bun"
)

type Repository[E any] interface {
	GetAll(ctx context.Context, tweak func(q *bun.SelectQuery) *bun.SelectQuery, relations ...string) (*[]E, error)
	GetOne(ctx context.Context, id int64) (*E, error)
	GetOneBy(ctx context.Context, column string, value interface{}) (*E, error)
	Create(ctx context.Context, data *E) (*E, error)
	Update(ctx context.Context, data *E) (*E, error)
}

type repository[E any] struct {
	db *bun.DB
}

func NewRepository[E any](db *bun.DB) Repository[E] {
	return &repository[E]{
		db: db,
	}
}

func (r repository[E]) GetAll(ctx context.Context, tweak func(q *bun.SelectQuery) *bun.SelectQuery, relations ...string) (*[]E, error) {
	var models []E

	c := ctx.Value("type")
	fmt.Print(c)

	q := r.db.NewSelect().Model(&models)
	for _, rel := range relations {
		q = q.Relation(rel)
	}
	if tweak != nil {
		q = tweak(q)
	}
	if err := q.Scan(ctx); err != nil {
		return nil, err
	}

	return &models, nil
}

func (r repository[E]) GetOne(ctx context.Context, id int64) (*E, error) {
	entity := new(E)
	err := r.db.NewSelect().Model(entity).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r repository[E]) GetOneBy(ctx context.Context, column string, value interface{}) (*E, error) {
	entity := new(E)
	err := r.db.NewSelect().Model(entity).Where("? = ?", bun.Ident(column), value).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (r repository[E]) Create(ctx context.Context, data *E) (*E, error) {
	if _, err := r.db.NewInsert().Model(data).Returning("*").Exec(ctx); err != nil {
		return nil, err
	}

	return data, nil
}

func (r repository[E]) Update(ctx context.Context, data *E) (*E, error) {
	model := new(E)

	val := reflect.ValueOf(data).Elem()

	_, err := r.db.NewUpdate().
		Model(model).
		ExcludeColumn("created_at").
		Where("id = ?", val.FieldByName("Id")).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return data, nil
}

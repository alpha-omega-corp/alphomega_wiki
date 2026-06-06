package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"wiki/app/controllers"
	"wiki/app/middleware"
	"wiki/app/models"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/bunrouterotel"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"github.com/urfave/cli/v3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	appCli := &cli.Command{
		Usage: "application cli",
		Commands: []*cli.Command{
			createCommand("server", "serve", runServer),
			createCommand("db", "migrate", runMigrations),
			createCommand("db", "sync", runSync),
		},
	}

	if err := appCli.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("app start error: %v\n", err)
	}
}

func newDB() *bun.DB {
	conn := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN(os.Getenv("DSN")),
	))

	db := bun.NewDB(conn, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(true),
		bundebug.WithVerbose(true),
	))
	return db
}

// runMigrations (re)creates the schema and reloads the seed rows from
// cmd/migration.yml. Run this before `serve` — serving no longer resets the DB.
func runMigrations(ctx context.Context, _ *cli.Command) {
	db := newDB()
	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)

	migrator := migrate.NewMigrator(db, migrate.NewMigrations())
	if err := migrator.Init(ctx); err != nil {
		log.Fatalf("migrator init: %v", err)
	}

	if err := db.ResetModel(ctx, schemaModels...); err != nil {
		log.Fatalf("reset model: %v", err)
	}

	fixture := dbfixture.New(db)
	if err := fixture.Load(ctx, os.DirFS("cmd"), "migration.yml"); err != nil {
		log.Fatalf("load fixture: %v", err)
	}

	// Advance the id sequence past the seeded rows so future inserts don't
	// collide with existing fixture rows on pages_pkey.
	if _, err := db.NewRaw(
		`SELECT setval(pg_get_serial_sequence('pages', 'id'), COALESCE((SELECT MAX(id) FROM pages), 1))`,
	).Exec(ctx); err != nil {
		log.Fatalf("reset sequence: %v", err)
	}

	fmt.Println("migrations complete")
}

var schemaModels = []interface{}{
	(*models.Page)(nil),
	(*models.Favorite)(nil),
}

// runSync brings the database schema up to date non-destructively: it creates
// any missing tables and adds any model columns that are absent from existing
// tables. It never drops tables, columns, or rows, so no data is lost (unlike
// `migrate`, which resets everything).
func runSync(ctx context.Context, _ *cli.Command) {
	db := newDB()
	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)

	for _, model := range schemaModels {
		// Create the table only if it isn't there yet — an existing table and
		// its rows are left untouched.
		if _, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			log.Fatalf("create table: %v", err)
		}

		table := db.Table(reflect.TypeOf(model).Elem())

		// Columns the live table already has.
		var existing []string
		if err := db.NewSelect().
			ColumnExpr("column_name").
			TableExpr("information_schema.columns").
			Where("table_schema = current_schema()").
			Where("table_name = ?", table.Name).
			Scan(ctx, &existing); err != nil {
			log.Fatalf("inspect %s columns: %v", table.Name, err)
		}
		have := make(map[string]bool, len(existing))
		for _, c := range existing {
			have[c] = true
		}

		// Add any model column missing from the table. New columns are added
		// nullable (CreateTableSQLType carries the type, not the NOT NULL
		// constraint) so the ALTER can't fail on a table that already has rows —
		// existing rows get NULL, or the model default when one is defined.
		for _, f := range table.Fields {
			if have[f.Name] {
				continue
			}

			q := db.NewAddColumn().Model(model).IfNotExists()
			if f.SQLDefault != "" {
				q.ColumnExpr("? ? DEFAULT ?", bun.Ident(f.Name), bun.Safe(f.CreateTableSQLType), bun.Safe(f.SQLDefault))
			} else {
				q.ColumnExpr("? ?", bun.Ident(f.Name), bun.Safe(f.CreateTableSQLType))
			}
			if _, err := q.Exec(ctx); err != nil {
				log.Fatalf("add column %s.%s: %v", table.Name, f.Name, err)
			}
			fmt.Printf("added column %s.%s\n", table.Name, f.Name)
		}
	}

	fmt.Println("schema sync complete")
}

func runServer(_ context.Context, _ *cli.Command) {
	db := newDB()

	router := bunrouter.New(
		bunrouter.WithMiddleware(reqlog.NewMiddleware(
			reqlog.WithEnabled(true),
			reqlog.WithVerbose(true),
		)),

		bunrouter.Use(bunrouterotel.NewMiddleware(
			bunrouterotel.WithClientIP(),
		)))

	pageHandler := controllers.NewPageHandler(db)
	favoriteHandler := controllers.NewFavoriteController(db)

	issuer := os.Getenv("AUTH_ISSUER")
	jwksURL := os.Getenv("AUTH_JWKS_URL")
	auth := middleware.NewJWTVerifier(issuer, jwksURL).Middleware

	router.GET("/pages", pageHandler.GetAll)
	router.GET("/page/:id", pageHandler.GetOne)
	router.POST("/page", auth(pageHandler.Create))
	router.PUT("/page/:id", auth(pageHandler.Update))
	router.DELETE("/page/:id", auth(pageHandler.Delete))
	router.GET("/favorites", auth(favoriteHandler.List))
	router.POST("/favorite", auth(favoriteHandler.Add))
	router.DELETE("/favorite/:id", auth(favoriteHandler.Remove))

	httpServer := &http.Server{
		Addr:         os.Getenv("HOST"),
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
			log.Fatalf("ListenAndServe failed: %s", err)
		}
	}()

	fmt.Printf("listening on http://%s\n", httpServer.Addr)

	ch := make(chan os.Signal, 3)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	<-ch
}

func createCommand(category string, name string, action func(ctx context.Context, cmd *cli.Command)) *cli.Command {
	return &cli.Command{
		Name:     name,
		Category: category,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			action(ctx, cmd)

			return nil
		},
	}
}

package main

import (
	"context"
	"database/sql"
	"dwiki/app/controllers"
	"dwiki/app/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/migrate"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/bunrouterotel"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func main() {
	ctx := context.Background()
	conn := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN("postgres://wiki:wiki@localhost:5432/wiki?sslmode=disable"),
	))

	db := bun.NewDB(conn, pgdialect.New())
	db.WithQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(true),
		bundebug.WithVerbose(true),
	))

	migrator := migrate.NewMigrator(db, migrate.NewMigrations())
	if err := migrator.Init(ctx); err != nil {
		panic(err)
	}

	if err := db.ResetModel(ctx, []interface{}{
		(*models.Page)(nil),
	}...); err != nil {
		panic(err)
	}

	fixture := dbfixture.New(db)
	if err := fixture.Load(ctx, os.DirFS("cmd"), "migration.yml"); err != nil {
		fmt.Printf("load fixture error: %v\n", err)
		panic(err)
	}

	router := bunrouter.New(
		bunrouter.WithMiddleware(reqlog.NewMiddleware(
			reqlog.WithEnabled(true),
			reqlog.WithVerbose(true),
		)),

		bunrouter.Use(bunrouterotel.NewMiddleware(
			bunrouterotel.WithClientIP(),
		)))

	articleController := controllers.NewArticleController(db)
	navigationController := controllers.NewNavigationController(db)

	router.POST("/page", articleController.Create)
	router.GET("/page/:slug", articleController.GetOne)
	router.GET("/pages", articleController.GetAll)
	router.GET("/navigation", navigationController.GetAll)

	httpServer := &http.Server{
		Addr:         "localhost:1000",
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

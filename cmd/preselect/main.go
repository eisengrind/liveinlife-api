package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/51st-state/api/pkg/apis/preselect"
	"github.com/51st-state/api/pkg/apis/preselect/cockroachdb"
	"github.com/51st-state/api/pkg/encode"

	"github.com/51st-state/api/pkg/api"
	"github.com/playnet-public/flagenv"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

var (
	httpAddr   = flagenv.String("http-addr", ":8080", "the http address of this service")
	dbHost     = flagenv.String("db-host", "localhost", "the host of the database")
	dbPort     = flagenv.Int("db-port", 1234, "the port of the database")
	dbUsername = flagenv.String("db-username", "user", "the username of the database")
	dbPassword = flagenv.String("db-password", "1234", "the password of the database")
	dbName     = flagenv.String("db-name", "preselect", "the name of the database")
)

func main() {
	flagenv.Parse()

	logger, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Fatal(err.Error())
	}

	logger.Info("connecting to database")
	db, err := makeCockroachDBDatabase()
	if err != nil {
		logger.Fatal("failed creating database", zap.Error(err))
	}

	logger.Info("creating database schema")
	if err := cockroachdb.CreateSchema(context.Background(), db); err != nil {
		logger.Fatal("failed creating database scheme")
	}

	m := preselect.NewManager(cockroachdb.NewRepository(db))

	a := api.New(*httpAddr, logger)
	a.Get(
		"/character/preselect/count",
		preselect.MakeGetLeftPreSelectionsEndpoint(logger, m, encode.NewJSONEncoder()),
	)
	a.Get(
		"/character/preselect/next",
		preselect.MakeGetNextPreSelectionsEndpoint(logger, m, encode.NewJSONEncoder()),
	)
	a.Post(
		"/character/preselect",
		preselect.MakeCreatePreSelectionsEndpoint(logger, m, encode.NewJSONEncoder()),
	)
	a.Patch(
		"/character/preselect",
		preselect.MakeSetPreSelectionsEndpoint(logger, m, encode.NewJSONEncoder()),
	)

	if err := a.Serve(); err != nil {
		logger.Fatal("http server failed listening", zap.Error(err))
	}
}

func makeCockroachDBDatabase() (*sql.DB, error) {
	return sql.Open("postgres", fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		*dbUsername,
		*dbPassword,
		*dbHost,
		*dbPort,
		*dbName,
	))
}

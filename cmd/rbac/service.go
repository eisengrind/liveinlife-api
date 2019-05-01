package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/51st-state/api/pkg/rbac"

	"github.com/51st-state/api/pkg/rbac/cockroachdb"
	pb "github.com/51st-state/api/pkg/rbac/proto"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/playnet-public/flagenv"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

var (
	grpcAddr   = flagenv.String("grpc-addr", ":1234", "the grpc addr to host the grpc server on")
	dbHost     = flagenv.String("db-host", "localhost", "the host of the database")
	dbPort     = flagenv.Int("db-port", 1234, "the port of the database")
	dbUsername = flagenv.String("db-username", "user", "the username of the database")
	dbPassword = flagenv.String("db-password", "1234", "the password of the database")
	dbName     = flagenv.String("db-name", "preselect", "the name of the database")
)

func main() {
	flagenv.Parse()

	l, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Fatal(err.Error())
	}

	l.Info("connecting to database")
	db, err := makeCockroachDBDatabase()
	if err != nil {
		l.Fatal(err.Error())
	}

	l.Info("creating database schema")
	if err := cockroachdb.CreateSchema(context.Background(), db); err != nil {
		l.Fatal(err.Error())
	}

	l.Info(fmt.Sprintf("creating grpc listener on %s", *grpcAddr))
	grpcListener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		l.Fatal(err.Error())
	}

	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			grpcZap.StreamServerInterceptor(l),
		)),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcZap.UnaryServerInterceptor(l),
		)),
	)
	pb.RegisterControlServer(
		grpcServer,
		rbac.NewGRPCServer(
			rbac.NewControl(
				cockroachdb.NewRepository(
					db,
				),
			),
		),
	)

	if err := grpcServer.Serve(grpcListener); err != nil {
		l.Fatal(err.Error())
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

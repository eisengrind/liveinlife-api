package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/51st-state/api/pkg/encode"

	"github.com/51st-state/api/pkg/api"
	"github.com/51st-state/api/pkg/apis/role"

	"github.com/51st-state/api/pkg/rbac"
	"google.golang.org/grpc"

	"github.com/51st-state/api/pkg/keys"

	"github.com/51st-state/api/pkg/apis/role/cockroachdb"
	"github.com/playnet-public/flagenv"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

var (
	httpAddr        = flagenv.String("http-addr", ":8080", "the http address of the service")
	dbHost          = flagenv.String("db-host", "localhost", "the host of the database")
	dbPort          = flagenv.Int("db-port", 1234, "the port of the database")
	dbUsername      = flagenv.String("db-username", "user", "the username of the database")
	dbPassword      = flagenv.String("db-password", "1234", "the password of the database")
	dbName          = flagenv.String("db-name", "preselect", "the name of the database")
	publicKeyPath   = flagenv.String("public-key-path", "/secrets/public.pem", "the public key to validate jwt token")
	rbacGRPCAddress = flagenv.String("rbac-grpc-addr", "rbac-service:2345", "the grpc address to the rbac control")
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

	if err := cockroachdb.CreateSchema(context.Background(), db); err != nil {
		l.Fatal(err.Error())
	}

	publicKey, err := keys.GetPublicKey(*publicKeyPath)
	if err != nil {
		l.Fatal(err.Error())
	}

	l.Info("creating rbac grpc connection")
	rbacCtrl, rbacConn, err := makeRBACControl()
	if err != nil {
		l.Fatal(err.Error())
	}
	defer rbacConn.Close()

	m := role.NewManager(cockroachdb.NewRepository(db), rbacCtrl)

	a := api.New(*httpAddr, l)

	a.Get("/roles/{id}", role.MakeGetEndpoint(l, m, encode.NewJSONEncoder(), *publicKey, rbacCtrl))
	a.Patch("/roles/{id}", role.MakeSetEndpoint(l, m, encode.NewJSONEncoder(), *publicKey, rbacCtrl))
	a.Delete("/roles/{id}", role.MakeDeleteEndpoint(l, m, encode.NewJSONEncoder(), *publicKey, rbacCtrl))
	a.Post("/roles/{id}", role.MakeCreateEndpoint(l, m, encode.NewJSONEncoder(), *publicKey, rbacCtrl))

	if err := a.Serve(); err != nil {
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

func makeGRPCConn(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(time.Second*10),
	)
}

func makeRBACControl() (rbac.Control, *grpc.ClientConn, error) {
	conn, err := makeGRPCConn(*rbacGRPCAddress)
	if err != nil {
		return nil, nil, err
	}

	return rbac.NewGRPCClient(conn), conn, nil
}

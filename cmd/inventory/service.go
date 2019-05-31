package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/51st-state/api/pkg/api"
	"github.com/51st-state/api/pkg/encode"
	"github.com/51st-state/api/pkg/keys"
	"github.com/51st-state/api/pkg/rbac"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/51st-state/api/pkg/apis/inventory"
	"github.com/51st-state/api/pkg/apis/inventory/cockroachdb"
	pb "github.com/51st-state/api/pkg/apis/inventory/proto"

	"github.com/playnet-public/flagenv"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

var (
	httpAddr        = flagenv.String("http-addr", ":8080", "the http address of the service")
	grpcAddr        = flagenv.String("grpc-addr", ":2345", "the grpc addr of the service")
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

	m := inventory.NewManager(
		cockroachdb.NewRepository(db),
	)

	a := api.New(*httpAddr, l)
	a.Get("/inventory/{guid}", inventory.MakeGetEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))
	a.Patch("/inventory/{guid}/items/add", inventory.MakeAddItemEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))
	a.Patch("/inventory/{guid}/items/remove", inventory.MakeRemoveItemEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))
	a.Delete("/inventory{guid}", inventory.MakeDeleteEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))
	a.Post("/inventory", inventory.MakeCreateEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))

	go serveGrpc(l, m)

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

func serveGrpc(l *zap.Logger, m inventory.Manager) {
	l.Info("preparing grpc server")
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			grpcZap.StreamServerInterceptor(l),
		)),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcZap.UnaryServerInterceptor(l),
		)),
	)
	pb.RegisterManagerServer(s, inventory.NewGRPCServer(m))
	reflection.Register(s)

	listener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		l.Fatal(err.Error())
	}

	l.Info("starting grpc server")
	if err := s.Serve(listener); err != nil {
		l.Fatal(err.Error())
	}
}

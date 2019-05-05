package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"time"

	pubsubNSQ "github.com/51st-state/api/pkg/pubsub/nsq"

	"github.com/nsqio/go-nsq"

	"google.golang.org/grpc/reflection"

	"github.com/51st-state/api/pkg/encode"

	"github.com/51st-state/api/pkg/api"
	"github.com/51st-state/api/pkg/apis/user"
	"github.com/51st-state/api/pkg/event"

	"github.com/51st-state/api/pkg/rbac"
	"google.golang.org/grpc"

	"github.com/51st-state/api/pkg/keys"

	"github.com/51st-state/api/pkg/apis/user/cockroachdb"
	"github.com/51st-state/api/pkg/apis/user/mysql"
	pb "github.com/51st-state/api/pkg/apis/user/proto"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/playnet-public/flagenv"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var (
	httpAddr        = flagenv.String("http-addr", ":8080", "the http addr of the service")
	grpcAddr        = flagenv.String("grpc-addr", ":2345", "the grpc address of this service")
	nsqdAddr        = flagenv.String("nsqd-addr", "nsqd:4151", "the address of the nsq lookupd servers")
	publicKeyPath   = flagenv.String("public-key-path", "/secrets/public.key", "the public key to validate jwt token")
	rbacGRPCAddress = flagenv.String("rbac-grpc-addr", "rbac-service:2345", "the grpc address to the rbac control")

	dbHost        = flagenv.String("db-host", "localhost", "the host of the database")
	dbPort        = flagenv.Int("db-port", 1234, "the port of the database")
	dbUsername    = flagenv.String("db-username", "user", "the username of the database")
	dbPassword    = flagenv.String("db-password", "1234", "the password of the database")
	dbName        = flagenv.String("db-name", "preselect", "the name of the database")
	wcfDBHost     = flagenv.String("wcf-db-host", "localhost", "the host of the database of the wcf framework")
	wcfDBPort     = flagenv.Int("wcf-db-port", 1234, "the port of the wcf database")
	wcfDBUsername = flagenv.String("wcf-db-username", "user", "the username of the wcf database to login")
	wcfDBPassword = flagenv.String("wcf-db-password", "1234", "the password of the wcf database user")
	wcfDBName     = flagenv.String("wcf-db-name", "name", "the name of the wcf database")
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

	wcfDB, err := makeWCFMysqlDatabase()
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

	eventProd, err := makeNSQEventProducer()

	m := user.NewManager(
		cockroachdb.NewRepository(db),
		mysql.NewWCFRepository(wcfDB),
		eventProd,
		rbacCtrl,
	)

	a := api.New(*httpAddr, l)
	a.Get("/users/{uuid}", user.MakeGetEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))
	a.Get("/users/hash/{uuid}", user.MakeGetByGameSerialHashEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))
	a.Post("/users", user.MakeCreateEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))
	a.Delete("/users/{uuid}", user.MakeDeleteEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))
	a.Patch("/users/{uuid}", user.MakeUpdateEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))
	a.Get("/users/{uuid}/roles", user.MakeGetRolesEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))
	a.Patch("/users/{uuid}/roles", user.MakeSetRolesEndpoint(l, m, encode.NewJSONEncoder(), rbacCtrl, *publicKey))

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

func makeWCFMysqlDatabase() (*sql.DB, error) {
	return sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		*wcfDBUsername,
		*wcfDBPassword,
		*wcfDBHost,
		*wcfDBPort,
		*wcfDBName,
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

func makeNSQEventProducer() (*event.Producer, error) {
	p, err := nsq.NewProducer(*nsqdAddr, nsq.NewConfig())
	if err != nil {
		return nil, err
	}

	return event.NewProducer(pubsubNSQ.NewProducer(p, "events")), nil
}

func serveGrpc(l *zap.Logger, m *user.Manager) {
	l.Info("preparing grpc server")
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpcMiddleware.ChainStreamServer(
			grpcZap.StreamServerInterceptor(l),
		)),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcZap.UnaryServerInterceptor(l),
		)),
	)
	pb.RegisterManagerServer(s, user.NewGRPCServer(m))
	reflection.Register(s)

	l.Info("preparing grpc server")
	listener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		l.Fatal(err.Error())
	}

	if err := s.Serve(listener); err != nil {
		l.Fatal(err.Error())
	}
}

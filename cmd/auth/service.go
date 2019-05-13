package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/51st-state/api/pkg/recaptcha"

	"github.com/51st-state/api/pkg/encode"

	"github.com/51st-state/api/pkg/api"
	"github.com/51st-state/api/pkg/apis/auth"
	"github.com/51st-state/api/pkg/apis/user"

	"google.golang.org/grpc"

	"github.com/51st-state/api/pkg/keys"

	"github.com/51st-state/api/pkg/apis/auth/cockroachdb"
	"github.com/playnet-public/flagenv"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
)

var (
	httpAddr       = flagenv.String("http-addr", ":8080", "the http addr of the service")
	publicKeyPath  = flagenv.String("public-key-path", "/secrets/public.pem", "the public key to validate jwt token")
	privateKeyPath = flagenv.String("private-key-path", "/secrets/private.pem", "the private key to sign valid access token")
	grpcUserAddr   = flagenv.String("user-addr", "user:2345", "the grpc address to the user microservice")

	recaptchaPrivateKey = flagenv.String("recaptcha-private-key", "", "the private key for authenticating with a recaptcha")

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

	if err := cockroachdb.CreateSchema(context.Background(), db); err != nil {
		l.Fatal(err.Error())
	}

	publicKey, err := keys.GetPublicKey(*publicKeyPath)
	if err != nil {
		l.Fatal(err.Error())
	}

	privateKey, err := keys.GetPrivateKey(*privateKeyPath)
	if err != nil {
		l.Fatal(err.Error())
	}

	l.Info("creating rbac grpc connection")
	userMgr, userMgrConn, err := makeUserManager()
	if err != nil {
		l.Fatal(err.Error())
	}
	defer userMgrConn.Close()

	m := auth.NewManager(
		privateKey,
		cockroachdb.NewRepository(db),
		userMgr,
		&recaptcha.Verifier{*recaptchaPrivateKey},
	)

	a := api.New(*httpAddr, l)
	a.Post("/auth/login", auth.MakeLoginEndpoint(l, m, encode.NewJSONEncoder()))
	a.Post("/auth/login/recaptcha", auth.MakeRecaptchaLoginEndpoint(l, m, encode.NewJSONEncoder()))
	a.Post("/auth/refresh", auth.MakeRefreshTokenEndpoint(l, m, encode.NewJSONEncoder(), *publicKey))

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

func makeUserManager() (user.Manager, *grpc.ClientConn, error) {
	conn, err := makeGRPCConn(*grpcUserAddr)
	if err != nil {
		return nil, nil, err
	}

	return user.NewGRPCClient(conn), conn, nil
}

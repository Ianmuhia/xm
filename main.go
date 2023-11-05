package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/url"
	"os"
	"slices"
	"time"

	pb "xmserver/gen"

	"github.com/BurntSushi/toml"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	config := flag.String("config", "config.toml", "config file location")
	conf, err := Load(*config)
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", conf.Server["grpc"].Address, conf.Server["grpc"].Port))
	if err != nil {
		log.Fatal(err)
	}

	pg, err := NewPostgreSQL(logger, conf)
	if err != nil {
		log.Fatal(err)
	}

	usersService := NewUsersService(pg)
	companyService := NewCompanyService(pg)
	tokenService, err := NewTokenService(conf.Secrets.Jwt)
	if err != nil {
		log.Fatal(err)
	}
	inter := NewAuthInterceptor(tokenService)

	handler := NewCompanyHandler(usersService, companyService, tokenService)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(inter.Unary()),
	)
	pb.RegisterCompanyServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)
	slog.Info("grpc server running")
	log.Fatal(grpcServer.Serve(listener))
}

// NewPostgreSQL instantiates the PostgreSQL database using configuration defined in environment variables.
func NewPostgreSQL(l *slog.Logger, cfg Config) (*pgxpool.Pool, error) {
	dsn := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.Database.User, cfg.Database.Password),
		Host:     fmt.Sprintf("%s:%v", cfg.Database.Server, cfg.Database.Port),
		Path:     cfg.Database.Database,
		RawQuery: "sslmode=disable",
	}

	poolConfig, err := pgxpool.ParseConfig(dsn.String())
	if err != nil {
		return nil, err
	}
	poolConfig.HealthCheckPeriod = 10 * time.Second
	pool, errx := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if errx != nil {
		return nil, errx
	}

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	migrationString := `
	CREATE SCHEMA IF NOT EXISTS companies;
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE company_type AS ENUM (
    'Corporations',
    'NonProfit',
    'Cooperative',
    'SoleProprietorship'
);

CREATE TABLE company (
    id UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    name VARCHAR(15) NOT NULL UNIQUE,
    description VARCHAR(3000),
    employees INT NOT NULL,
    registered BOOLEAN NOT NULL,
    type company_type NOT NULL
);

CREATE TABLE users (
    name VARCHAR(255) NOT NULL
)
	
	`
	if cfg.Database.Automigrate {
		_, err := pool.Exec(context.Background(), migrationString)
		if err != nil {
			return nil, err
		}
		l.Info("migrations done")
	}

	l.Info("Database connected")
	return pool, nil
}

type Config struct {
	Database struct {
		Server      string
		Port        uint16
		Database    string
		User        string
		Password    string
		Timeout     int
		Automigrate bool
		Log         bool
	}
	Minio struct {
		Address    string
		Port       string
		BucketName string
		User       string
		Password   string
		Timeout    int
	}
	Secrets struct {
		Jwt string
	}

	Server map[string]struct {
		Address string
		Port    string
		Timeout int
		Prod    bool
	}
}

func Load(loc string) (Config, error) {
	var config Config
	if _, err := toml.DecodeFile(loc, &config); err != nil {
		return config, err
	}

	return config, nil
}

type AuthInterceptor struct {
	jwtManager TokenService
}

func NewAuthInterceptor(jwtManager TokenService) *AuthInterceptor {
	return &AuthInterceptor{jwtManager: jwtManager}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)
		authMethod := []string{"/CompanyService/Login", "/CompanyService/Register"}

		if !slices.Contains(authMethod, info.FullMethod) {
			md, _ := metadata.FromIncomingContext(ctx)
			token := md.Get("Authorization")
			if token == nil {
				return 0, status.Errorf(codes.PermissionDenied, "user isn't authorized")
			}
			_, err := interceptor.jwtManager.Verify(token[0])
			if err != nil {
				return nil, status.Error(codes.PermissionDenied, "no permission to access this RPC")
			}
		}

		return handler(ctx, req)
	}
}

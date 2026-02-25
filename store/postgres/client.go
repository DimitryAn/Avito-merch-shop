package postgres

import (
	"context"
	"fmt"
	"log"
	"root/store/postgres/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	DbPool *pgxpool.Pool
	Ctx    context.Context
}

func NewClient(ctx context.Context, cnf *config.Config) (*Client, error) {
	var dbpool *pgxpool.Pool
	var err error
	dsn := makeDsn(cnf)

	for cnf.MaxAttemptsConn > 0 {
		time.Sleep(1 * time.Second)
		dbpool, err = pgxpool.New(context.Background(), dsn)

		if err != nil {
			fmt.Printf("Unable to create connection pool to postgres: %v\n", err)
			cnf.MaxAttemptsConn--
		} else {
			log.Print("Success make dbpool to postgresSql")
			break
		}
	}

	if dbpool == nil || err != nil {
		log.Fatal("Unable to create connection pool to postgres")
	}

	if err := dbpool.Ping(ctx); err != nil {
		log.Fatal("Unable to ping connection pool to postgres")
	}
	log.Printf("Success connect to postgres sql, name - %v \n", cnf.DatabaseName)
	return &Client{dbpool, ctx}, nil
}

func makeDsn(cnf *config.Config) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cnf.DatabaseUser,
		cnf.DatabasePassword, cnf.DatabaseHost, cnf.DatabasePort, cnf.DatabaseName)
}

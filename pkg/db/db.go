package db

import (
	"context"
	"es/config"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"sync"
)

var once sync.Once

func CreateDbConn(config *config.ESConfig) *pgxpool.Pool {
	var err error
	var DbConnMain *pgxpool.Pool
	once.Do(func() {
		DbConnMain, err = pgxpool.Connect(context.Background(), fmt.Sprintf("postgresql://%s/%s?user=%s&password=%s&pool_max_conns=100", config.DbConfig.DbAddr, config.DbConfig.DbName, config.DbConfig.DbUser, config.DbConfig.DbPass))
		if err != nil {
			log.Fatalf("error in connecting to db %v", err)
		}
	})

	return DbConnMain
}

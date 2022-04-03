package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joineroff/social-network/backend/internal/config"
)

type MysqlDB struct {
	*sqlx.DB
}

func NewMsSQLDatabase(cfg *config.Config) *MysqlDB {
	conStr := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s",
		cfg.Mysql.User,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Database,
	)

	db, err := sqlx.Open("mysql", conStr)
	if err != nil {
		log.Fatalf("failed to connect mysql db: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping mysql db: %v", err)
	}

	return &MysqlDB{db}
}

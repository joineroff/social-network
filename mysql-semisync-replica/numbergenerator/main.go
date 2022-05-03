package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	db := newConnection(
		"0.0.0.0:6612",
		"appuser",
		"appPassw0rd",
		"testdb",
	)

	for i := 0; i < 100000; i++ {
		now := time.Now()

		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("failed to start transaction: %v", err)
		}

		defer tx.Rollback()

		r, err := tx.Exec("INSERT INTO numbers (createAt) VALUES (?)", now)
		if err != nil {
			log.Fatalf("failed to query table list: %v", err)
		}

		id, err := r.LastInsertId()
		if err != nil {
			log.Fatalf("failed to get last insertID: %v", err)
		}

		err = tx.Commit()
		if err != nil {
			_ = tx.Rollback()
			log.Fatalf("failed to commit transaction: %v", err)
		}

		fmt.Println(id)
	}
}

func newConnection(
	address string,
	user string,
	password string,
	database string,
) *sqlx.DB {
	conStr := fmt.Sprintf(
		"%s:%s@(%s)/%s",
		user,
		password,
		address,
		database,
	)

	db, err := sqlx.Open("mysql", conStr)
	if err != nil {
		log.Fatalf("failed to connect to %s db: %v", address, err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping to %s db: %v", address, err)
	}

	return db
}

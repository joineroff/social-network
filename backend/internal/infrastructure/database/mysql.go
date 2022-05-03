package database

import (
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joineroff/social-network/backend/internal/config"
)

const (
	defaultMaxIdleConnections = 5
	defaultMaxOpenConnections = 20
	defaultMaxLifeTimeSeconds = 600
)

type MysqlDBOption func(*sqlx.DB)

type MysqlDB struct {
	mu               sync.Mutex
	master           *sqlx.DB
	replicas         []*sqlx.DB
	nextReplicaIndex int
	lastReplicaIndex int
}

func (db *MysqlDB) Master() *sqlx.DB {
	return db.master
}

func (db *MysqlDB) Slave() *sqlx.DB {
	db.mu.Lock()
	curReplicaIndex := db.nextReplicaIndex

	db.nextReplicaIndex++

	if db.nextReplicaIndex > db.lastReplicaIndex {
		db.nextReplicaIndex = 0
	}
	db.mu.Unlock()

	return db.replicas[curReplicaIndex]
}

func (db *MysqlDB) Close() error {
	if err := db.master.Close(); err != nil {
		return err
	}

	for _, replica := range db.replicas {
		if err := replica.Close(); err != nil {
			return err
		}
	}

	return nil
}

func NewMsSQLDatabase(cfg *config.Config, opts ...MysqlDBOption) *MysqlDB {
	masterConn := newConnection(
		cfg.Infrastructure.Mysql.Address,
		cfg.Infrastructure.Mysql.User,
		cfg.Infrastructure.Mysql.Password,
		cfg.Infrastructure.Mysql.Database,
		opts...,
	)
	replicas := make([]*sqlx.DB, 0, 1)

	for _, replicaAddress := range cfg.Infrastructure.Mysql.Replicas {
		replicaConn := newConnection(
			replicaAddress,
			cfg.Infrastructure.Mysql.User,
			cfg.Infrastructure.Mysql.Password,
			cfg.Infrastructure.Mysql.Database,
			opts...,
		)

		replicas = append(replicas, replicaConn)
	}

	if len(replicas) == 0 {
		replicas = append(replicas, masterConn)
	}

	return &MysqlDB{
		master:           masterConn,
		replicas:         replicas,
		lastReplicaIndex: len(replicas) - 1,
	}
}

func newConnection(
	address string,
	user string,
	password string,
	database string,
	opts ...MysqlDBOption,
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

	db.SetMaxIdleConns(defaultMaxIdleConnections)
	db.SetMaxOpenConns(defaultMaxOpenConnections)
	db.SetConnMaxLifetime(defaultMaxLifeTimeSeconds * time.Second)

	for _, o := range opts {
		o(db)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping to %s db: %v", address, err)
	}

	return db
}

func WithMaxIdleConnections(connections int) MysqlDBOption {
	return func(db *sqlx.DB) {
		db.SetMaxIdleConns(connections)
	}
}

func WithMaxLifetime(t time.Duration) MysqlDBOption {
	return func(db *sqlx.DB) {
		db.SetConnMaxLifetime(t)
	}
}

func WithMaxOpenConnections(connections int) MysqlDBOption {
	return func(db *sqlx.DB) {
		db.SetMaxOpenConns(connections)
	}
}

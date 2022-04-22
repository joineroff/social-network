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
	mu        sync.Mutex
	master    *sqlx.DB
	slaves    []*sqlx.DB
	nextSlave int
	lastSlave int
}

func (db *MysqlDB) Master() *sqlx.DB {
	return db.master
}

func (db *MysqlDB) Slave() *sqlx.DB {
	db.mu.Lock()
	curSlave := db.nextSlave

	db.nextSlave++

	if db.nextSlave > db.lastSlave {
		db.nextSlave = 0
	}
	db.mu.Unlock()

	return db.slaves[curSlave]
}

func (db *MysqlDB) Close() error {
	for _, slave := range db.slaves {
		if err := slave.Close(); err != nil {
			return err
		}
	}

	if err := db.master.Close(); err != nil {
		return err
	}

	return nil
}

func NewMsSQLDatabase(cfg *config.Config, opts ...MysqlDBOption) *MysqlDB {
	masterDB := newConnection(cfg.Mysql.Master, opts...)
	slaveDBs := make([]*sqlx.DB, 0, 1)

	for _, slave := range cfg.Mysql.Slaves {
		slaveDB := newConnection(slave, opts...)

		slaveDBs = append(slaveDBs, slaveDB)
	}

	if len(slaveDBs) == 0 {
		slaveDBs = append(slaveDBs, masterDB)
	}

	return &MysqlDB{
		master:    masterDB,
		slaves:    slaveDBs,
		lastSlave: len(slaveDBs) - 1,
	}
}

func newConnection(cfg config.MysqlConfig, opts ...MysqlDBOption) *sqlx.DB {
	conStr := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sqlx.Open("mysql", conStr)
	if err != nil {
		log.Fatalf("failed to connect to %s:%s db: %v", cfg.Host, cfg.Port, err)
	}

	db.SetMaxIdleConns(defaultMaxIdleConnections)
	db.SetMaxOpenConns(defaultMaxOpenConnections)
	db.SetConnMaxLifetime(defaultMaxLifeTimeSeconds * time.Second)

	for _, o := range opts {
		o(db)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping to %s:%s db: %v", cfg.Host, cfg.Port, err)
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

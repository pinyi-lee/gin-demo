package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/vgarvardt/go-pg-adapter/pgx4adapter"
)

var (
	postgresManager *PostgresManager
)

func GetPostgres() *PostgresManager {
	return postgresManager
}

type PostgresManager struct {
	conn    *pgxpool.Pool
	config  PostgresConfig
	context context.Context
}

type PostgresConfig struct {
	Username                string
	Password                string
	Host                    string
	Port                    string
	DatabaseName            string
	MinConnSize             int32
	MaxConnSize             int32
	MaxConnIdleTimeBySecond time.Duration
	MaxConnLifetimeBySecond time.Duration
}

func (manager *PostgresManager) Setup(config PostgresConfig) (err error) {

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.DatabaseName)

	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		fmt.Printf("postgres parse config fail, %+v\n", err)
		return err
	}

	background := context.Background()
	pgxConfig.MinConns = config.MinConnSize
	pgxConfig.MaxConns = config.MaxConnSize
	pgxConfig.MaxConnIdleTime = config.MaxConnIdleTimeBySecond * time.Second
	pgxConfig.MaxConnLifetime = config.MaxConnLifetimeBySecond * time.Second

	pool, err := pgxpool.ConnectConfig(background, pgxConfig)
	if err != nil {
		fmt.Printf("postgres connect fail, %+v\n", err)
		return err
	}

	err = pool.Ping(background)
	if err != nil {
		fmt.Printf("postgres ping fail, %+v\n", err)
		return err
	}

	postgresManager = &PostgresManager{
		conn:    pool,
		config:  config,
		context: background,
	}

	return nil
}

func (manager *PostgresManager) Close() {

}

func (manager *PostgresManager) GetAdapter() *pgx4adapter.Pool {
	return pgx4adapter.NewPool(manager.conn)
}

func (manager *PostgresManager) Query(sql string, args ...interface{}) (pgx.Rows, error) {
	return manager.conn.Query(manager.context, sql, args...)
}

func (manager *PostgresManager) Exec(sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return manager.conn.Exec(manager.context, sql, arguments...)
}

func (manager *PostgresManager) TxQuery(tx *pgx.Tx, sql string, args ...interface{}) (pgx.Rows, error) {
	return (*tx).Query(manager.context, sql, args...)
}

func (manager *PostgresManager) TxExec(tx *pgx.Tx, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	return (*tx).Exec(manager.context, sql, arguments...)
}

func (manager *PostgresManager) DoInTx(fn func(*pgx.Tx) error) error {

	tx, err := manager.conn.Begin(manager.context)
	if err != nil {
		return err
	}
	defer tx.Rollback(manager.context)

	err = fn(&tx)
	if err != nil {
		return err
	}

	return tx.Commit(manager.context)
}

func (manager *PostgresManager) DoInTxInterface(fn func(*pgx.Tx) (interface{}, error)) (interface{}, error) {

	tx, err := manager.conn.Begin(manager.context)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(manager.context)

	in, fnErr := fn(&tx)
	if fnErr != nil {
		return in, fnErr
	}

	return in, tx.Commit(manager.context)
}

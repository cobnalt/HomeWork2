package database

import (
	"context"
	"fmt"
	"time"

	"../config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Database interface
type Database interface {
	GetCurrencies(ctx context.Context) (result []Currency, err error)
	CreateCurrency(ctx context.Context, c Currency) error
	DeleteCurrency(ctx context.Context, id int) error
	UpdateCurrency(ctx context.Context, c Currency, id int) error
	Close()
}

// DB connect
type DB struct {
	conn *sqlx.DB
}

// New DB
func New(cfg config.DatabaseConfig) (*DB, error) {
	conn, err := sqlx.Connect("postgres", cfg.ConnectionString)
	if err != nil {
		return nil, err
	}
	return &DB{
		conn: conn,
	}, nil
}

// Currency Structure
type Currency struct {
	tableName  struct{} `sql:"currencyPair"`
	ID         int
	Currency1  string
	Currency2  string
	Rate       float64
	Lastupdate time.Time
}

// GetCurrencies func
func (d *DB) GetCurrencies(ctx context.Context) (result []Currency, err error) {
	q := "SELECT id, currency1, currency2, rate, lastupdate FROM currencyPair;"
	if err = d.conn.SelectContext(ctx, &result, q); err != nil {
		return nil, err
	}
	return result, err
}

// CreateCurrency func
func (d *DB) CreateCurrency(ctx context.Context, c Currency) error {
	q := "INSERT INTO currencyPair (currency1, currency2, rate) VALUES ($1, $2, $3);"
	_, err := d.conn.ExecContext(ctx, q, c.Currency1, c.Currency2, c.Rate)
	return err
}

// DeleteCurrency func
func (d *DB) DeleteCurrency(ctx context.Context, id int) error {
	q := "DELETE FROM currencyPair WHERE id = $1;"
	_, err := d.conn.ExecContext(ctx, q, id)
	return err
}

// UpdateCurrency func
func (d *DB) UpdateCurrency(ctx context.Context, c Currency, id int) error {
	q := "UPDATE currencyPair SET currency1 = $1, currency2 = $2, rate = $3 WHERE id = $4;"
	_, err := d.conn.ExecContext(ctx, q, c.Currency1, c.Currency2, c.Rate, id)
	return err
}

// Close DB
func (d *DB) Close() {
	d.conn.Close()
}

// GetDBCurrencies func
func (d *DB) GetDBCurrencies() ([]Currency, error) {
	result := []Currency{}
	q := "SELECT id, currency1, currency2, rate FROM currencyPair;"
	rows, err := d.conn.Query(q)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		res := Currency{}
		err := rows.Scan(&res.ID, &res.Currency1, &res.Currency2, &res.Rate)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, res)
	}
	return result, err
}

// UpdateDBCurrencies func
func (d *DB) UpdateDBCurrencies(c Currency, val float64) error {
	q := "UPDATE currencyPair SET currency1 = $1, currency2 = $2, rate = $3, lastupdate = $4 WHERE id = $5;"
	_, err := d.conn.Exec(q, c.Currency1, c.Currency2, val, time.Now(), c.ID)
	return err
}

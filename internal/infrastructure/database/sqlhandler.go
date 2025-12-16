package databaseInfra

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"Aicon-assignment/internal/infrastructure/config"
	"Aicon-assignment/internal/interfaces/database"
)

type MySqlHandler struct {
	Conn *sql.DB
}

func NewSqlHandler() database.SqlHandler {
	dsn := config.GetDSN()
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("❌ Failed to connect to database: %v", err))
	}

	// DB接続が確立できているかを確認
	if err := conn.Ping(); err != nil {
		panic(fmt.Sprintf("❌ Failed to ping database: %v", err))
	}

	fmt.Println("✅ Successfully connected to the database!")

	sqlBytes, err := os.ReadFile("sql/init.sql")
	if err != nil {
		fmt.Printf("❌ Failed to read init.sql: %v\n", err)
	} else {
		if _, err := conn.Exec(string(sqlBytes)); err != nil {
			fmt.Printf("❌ Failed to execute init.sql: %v\n", err)
		} else {
			fmt.Println("✅ Successfully initialized database from init.sql")
		}
	}

	return &MySqlHandler{Conn: conn}
}

func (h *MySqlHandler) Execute(ctx context.Context, statement string, args ...interface{}) (database.Result, error) {
	result, err := h.Conn.ExecContext(ctx, statement, args...)
	if err != nil {
		return nil, err
	}
	return &mysqlResult{result: result}, nil
}

func (h *MySqlHandler) Query(ctx context.Context, statement string, args ...interface{}) (database.Rows, error) {
	rows, err := h.Conn.QueryContext(ctx, statement, args...)
	if err != nil {
		return nil, err
	}
	return &mysqlRows{rows: rows}, nil
}

func (h *MySqlHandler) QueryRow(ctx context.Context, statement string, args ...interface{}) database.Row {
	row := h.Conn.QueryRowContext(ctx, statement, args...)
	return &mysqlRow{row: row}
}

func (h *MySqlHandler) Close() error {
	if h.Conn != nil {
		return h.Conn.Close()
	}
	return nil
}

type mysqlResult struct {
	result sql.Result
}

func (r *mysqlResult) LastInsertId() (int64, error) {
	return r.result.LastInsertId()
}

func (r *mysqlResult) RowsAffected() (int64, error) {
	return r.result.RowsAffected()
}

type mysqlRows struct {
	rows *sql.Rows
}

func (r *mysqlRows) Next() bool {
	return r.rows.Next()
}

func (r *mysqlRows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r *mysqlRows) Close() error {
	return r.rows.Close()
}

func (r *mysqlRows) Err() error {
	return r.rows.Err()
}

type mysqlRow struct {
	row *sql.Row
}

func (r *mysqlRow) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}

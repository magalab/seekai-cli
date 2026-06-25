package client

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Client struct {
	db *sql.DB
}

type QueryResult struct {
	Columns []string
	Rows    []map[string]any
}

func Open(cfg Config) (*Client, error) {
	q := url.Values{}
	q.Set("charset", "utf8mb4")
	q.Set("parseTime", "true")
	q.Set("timeout", "8s")
	q.Set("readTimeout", "0")
	q.Set("writeTimeout", "30s")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, q.Encode())
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return &Client{db: db}, nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) Query(statement string, args ...any) ([]map[string]any, error) {
	result, err := c.QueryResult(statement, args...)
	if err != nil {
		return nil, err
	}
	return result.Rows, nil
}

func (c *Client) QueryResult(statement string, args ...any) (QueryResult, error) {
	rows, err := c.db.Query(statement, args...)
	if err != nil {
		return QueryResult{}, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (c *Client) Exec(statement string, args ...any) error {
	_, err := c.db.Exec(statement, args...)
	return err
}

func scanRows(rows *sql.Rows) (QueryResult, error) {
	columns, err := rows.Columns()
	if err != nil {
		return QueryResult{}, err
	}
	out := []map[string]any{}
	for rows.Next() {
		values := make([]any, len(columns))
		ptrs := make([]any, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return QueryResult{}, err
		}
		row := map[string]any{}
		for i, column := range columns {
			row[column] = normalizeValue(values[i])
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return QueryResult{}, err
	}
	return QueryResult{Columns: columns, Rows: out}, nil
}

func normalizeValue(value any) any {
	switch v := value.(type) {
	case []byte:
		return string(v)
	case time.Time:
		return v.Format(time.RFC3339)
	default:
		return v
	}
}

package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/denisenkom/go-mssqldb"
)

type SQLStr struct {
	url *url.URL
	db  *sql.DB
}


func MakeSQL(host, port, username, password string) (*SQLStr, error) {
	server := &SQLStr{}
	server.url = &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword(username, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
	}
	return server, server.connect()
}


func (server *SQLStr) connect() error {
	var err error
	server.db, err = sql.Open("sqlserver", server.url.String())
	if err != nil {
		return err
	}
	return server.db.PingContext(context.Background())
}

func (s *SQLStr) DB() *sql.DB {
	return s.db
}


func (s *SQLStr) GetTokenTiny() (string, error) {
	var token string
	err := s.db.QueryRow("SELECT TOP 1 Token FROM ConfigTiny ORDER BY Id DESC").Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

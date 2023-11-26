package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	_ "modernc.org/sqlite"
)

type SqliteRepo struct {
	logger *slog.Logger
	*sql.DB
}

func (s *SqliteRepo) Insert(data string) error {
	s.logger.Info("inserting hash in to content store", "hash", data)
	_, err := s.DB.Exec("INSERT INTO content (hash) VALUES (?);", data)
	if err != nil {
		return fmt.Errorf("error inserting into db %v", err)
	}

	s.logger.Info("hash inserted successfully")
	return nil
}

func (s *SqliteRepo) Remove(data string) error {
	//TODO implement me
	panic("implement me")
}

func (s *SqliteRepo) Get(data string) (string, error) {
	var queryData string

	s.logger.Info("getting hash from db", "hash", data)
	query := s.DB.QueryRow("SELECT hash FROM content WHERE hash = ?;", data)

	if queryErr := query.Scan(&queryData); queryErr != nil {
		if errors.Is(queryErr, sql.ErrNoRows) {
			return "", nil
		}
		return "", queryErr
	}

	return queryData, nil
}

func NewSqliteRepo(dsnUri string, logger *slog.Logger) *SqliteRepo {
	sr := new(SqliteRepo)
	sr.logger = logger

	db, err := sql.Open("sqlite", dsnUri)
	if err != nil {
		sr.logger.Error("could not connect to sqlite db", "error", err)
		panic(err)
	}
	sr.DB = db
	sr.logger.Info("connected to sqlite db", "db", db)

	return sr
}

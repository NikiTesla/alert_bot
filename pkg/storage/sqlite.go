package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"alert_bot/pkg/model"

	_ "github.com/mattn/go-sqlite3"
)

const sqliteStorageFilename = "db.sql"

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(filename string) (*SQLiteStorage, error) {
	db, err := connectSQLite(filename)
	if err != nil {
		return nil, fmt.Errorf("connectSQLite: %w", err)
	}
	return &SQLiteStorage{
		db: db,
	}, nil
}

func connectSQLite(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping: %w", err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS subscribers (
		chat_id INTEGER PRIMARY KEY,
		status TEXT
	)`)
	if err != nil {
		return nil, fmt.Errorf("create table: %w", err)
	}
	return db, nil
}

func (s *SQLiteStorage) Subscribe(chatId int64) error {
	_, err := s.db.Exec("INSERT INTO subscribers (chat_id) VALUES (?)", chatId)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStorage) Unsubscribe(chatId int64) error {
	_, err := s.db.Exec("DELETE FROM subscribers WHERE chat_id = ?", chatId)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStorage) GetSubscribersUids() ([]int64, error) {
	rows, err := s.db.Query("SELECT chat_id FROM subscribers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatIds []int64
	for rows.Next() {
		var chatId int64
		if err := rows.Scan(&chatId); err != nil {
			return nil, err
		}
		chatIds = append(chatIds, chatId)
	}
	return chatIds, nil
}

func (s *SQLiteStorage) SetStatus(chatId int64, status model.Status) error {
	_, err := s.db.Exec("UPDATE subscribers SET status = ? WHERE chat_id = ?", status, chatId)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStorage) GetStatus(chatId int64) (model.Status, error) {
	var status model.Status
	err := s.db.QueryRow("SELECT status FROM subscribers WHERE chat_id = ?", chatId).Scan(&status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return status, nil
		}
		return status, err
	}
	return status, nil
}

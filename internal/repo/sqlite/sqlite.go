package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type SqliteRepo struct {
	db *sql.DB
}

func New(path string) (*SqliteRepo, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Создание таблиц (пример)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (id INTEGER PRIMARY KEY, message TEXT,tag TEXT)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return &SqliteRepo{db: db}, nil
}

func (s *SqliteRepo) GetMessages(tag string) ([]string, error) {
	query := "SELECT message FROM messages WHERE tag = ?"

	rows, err := s.db.Query(query, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []string

	for rows.Next() {
		var message string
		err := rows.Scan(&message) // Считываем значения в переменные
		if err != nil {
			log.Print("Ошибка при чтении строки: %v", err)
		}
		res = append(res, message)
	}
	return res, nil
}

func (s *SqliteRepo) SaveMessage(tag string, msg string) error {
	query := "INSERT INTO messages (tag,message) VALUES (?,?)"
	_, err := s.db.Exec(query, tag, msg)
	if err != nil {
		return fmt.Errorf("failed to add message: %v", err)
	}
	return nil
}

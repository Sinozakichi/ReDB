package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Card 表示 cards 表中的一條記錄
// type Card struct {
// 	ID        int
// 	Name      string
// 	Level     int
// 	Attribute string
// 	Race      string
// 	Attack    int
// 	Defense   int
// 	Effect    string
// }

// Database 是資料庫操作的介面
// type Database interface {
// 	Close() error
// 	CreateTable() error
// 	InsertCard(card Card) error
// 	GetAllCards() ([]Card, error)
// }

// SQLiteDatabase 是 SQLite 的實現
type SQLiteDatabase struct {
	db *sql.DB
}

// NewSQLiteDatabase 初始化 SQLite 資料庫
func NewSQLiteDatabase(filepath string) (*SQLiteDatabase, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	return &SQLiteDatabase{db: db}, nil
}

// Close 關閉資料庫連線
func (s *SQLiteDatabase) Close() error {
	return s.db.Close()
}

// CreateTable 建立表格
func (s *SQLiteDatabase) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS cards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		level INTEGER,
		attribute TEXT,
		race TEXT,
		attack INTEGER,
		defense INTEGER,
		effect TEXT
	);`
	_, err := s.db.Exec(query)
	return err
}

// InsertCard 插入一條卡片資料
func (s *SQLiteDatabase) InsertCard(card Card) error {
	query := `INSERT INTO cards (name, level, attribute, race, attack, defense, effect) VALUES (?, ?, ?, ?, ?, ?, ?);`
	_, err := s.db.Exec(query, card.Name, card.Level, card.Attribute, card.Race, card.Attack, card.Defense, card.Effect)
	return err
}

// GetAllCards 查詢所有卡片資料
func (s *SQLiteDatabase) GetAllCards() ([]Card, error) {
	query := `SELECT * FROM cards;`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var card Card
		err = rows.Scan(&card.ID, &card.Name, &card.Level, &card.Attribute, &card.Race, &card.Attack, &card.Defense, &card.Effect)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

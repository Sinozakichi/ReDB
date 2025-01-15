package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLDatabase 是 MySQL 的實現
type MySQLDatabase struct {
	db *sql.DB
}

// NewMySQLDatabase 初始化 MySQL 資料庫
func NewMySQLDatabase(user, password, host, dbname string) (*MySQLDatabase, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return &MySQLDatabase{db: db}, nil
}

// Close 關閉資料庫連線
func (m *MySQLDatabase) Close() error {
	return m.db.Close()
}

// CreateTable 建立表格
func (m *MySQLDatabase) CreateTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS cards (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		level INT,
		attribute VARCHAR(50),
		race VARCHAR(50),
		attack INT,
		defense INT,
		effect TEXT
	);
	`
	_, err := m.db.Exec(query)
	return err
}

// InsertCard 插入一條卡片資料
func (m *MySQLDatabase) InsertCard(card Card) error {
	query := `INSERT INTO cards (name, level, attribute, race, attack, defense, effect) VALUES (?, ?, ?, ?, ?, ?, ?);`
	_, err := m.db.Exec(query, card.Name, card.Level, card.Attribute, card.Race, card.Attack, card.Defense, card.Effect)
	return err
}

// GetAllCards 查詢所有卡片資料
func (m *MySQLDatabase) GetAllCards() ([]Card, error) {
	query := `SELECT * FROM cards;`
	rows, err := m.db.Query(query)
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

// UpdateCard 更新一條卡片資料
func (m *MySQLDatabase) UpdateCard(card Card) error {
	query := `UPDATE cards SET name = ?, level = ?, attribute = ?, race = ?, attack = ?, defense = ?, effect = ? WHERE id = ?;`
	_, err := m.db.Exec(query, card.Name, card.Level, card.Attribute, card.Race, card.Attack, card.Defense, card.Effect, card.ID)
	return err
}

// DeleteCard 刪除一條卡片資料
func (m *MySQLDatabase) DeleteCard(id int) error {
	query := `DELETE FROM cards WHERE id = ?;`
	_, err := m.db.Exec(query, id)
	return err
}

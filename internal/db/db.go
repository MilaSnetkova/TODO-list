package db

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/MilaSnetkova/TODO-list/internal/config"
)


func ConnectDB(cfg *config.Config) (*sqlx.DB, error) {
	dbPath := cfg.DBFile
	if dbPath == "" {
		appPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("error getting executable path: %w", err)
		}
		dbPath = filepath.Join(filepath.Dir(appPath), "scheduler.db")
	}

	// Проверка существования базы данных
	_, err := os.Stat(dbPath)
	install := os.IsNotExist(err)

	// Подключение к базе данных
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Если база данных не существует, создаём таблицу и индекс
	if install {
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date TEXT NOT NULL,
				title TEXT NOT NULL,
				comment TEXT,
				repeat TEXT CHECK(length(repeat) <= 128)
			);
			CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
		`)
		if err != nil {
			return nil, fmt.Errorf("error creating database: %w", err)
		}
		fmt.Println("Scheduler table created successfully")
	}

	return db, nil
}
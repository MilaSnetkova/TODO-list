package repository

import (
	"database/sql"
	"time"
	"fmt"
	 
	"github.com/jmoiron/sqlx"
	"github.com/MilaSnetkova/TODO-list/internal/constants"
	"github.com/MilaSnetkova/TODO-list/internal/models"
    
)

// TaskRepository определяет интерфейс для работы с задачами
type TaskRepository interface {
	Create(task *models.Task) (int64, error)
	SearchTasks(filter Filter, id string) ([]models.Task, error)
	UpdateTask(task *models.Task) error
	Delete(id string) error

}

// Filter используется для фильтрации задач
type Filter struct {
	ID     []string
	Search string
	Date   string
}

type TaskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) TaskRepository {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(task *models.Task) (int64, error) {
	res, err := r.db.Exec(
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?,?,?,?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Получение списка задач с фильтрацией
func (r *TaskRepo) SearchTasks(filter Filter, id string) ([]models.Task, error) {
	var tasks []models.Task

	// Начальное условие
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE 1=1"
	var params []interface{}

	// Если передан ID, то ищем по ID
	if id != "" {
		query += " AND id = ?"
		params = append(params, id)
	}

	// Выполняем фильтрацию по дате или заголовку/комментарию
	if filter.Search != "" {
		parsedDate, err := time.Parse("02.01.2006", filter.Search)
		if err == nil {
			query += " AND date = ?"
			params = append(params, parsedDate.Format(constants.DateFormat))
		} else {
			query += " AND (LOWER(title) LIKE LOWER(?) OR LOWER(comment) LIKE LOWER(?))"
			search := "%" + filter.Search + "%"
			params = append(params, search, search)
		}
	}

	// Добавляем сортировку и лимит
	query += " ORDER BY date LIMIT ?"
	params = append(params, constants.Limit)

	
	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	// Если ищем по ID, то возвращаем ошибку, если задача не найдена
	if id != "" && len(tasks) == 0 {
		return nil, fmt.Errorf("task not found")
	}

	return tasks, nil
}

// Обновление задачи
func (r *TaskRepo) UpdateTask(task *models.Task) error {
	res, err := r.db.Exec(
		"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// Удаление задачи
func (r *TaskRepo) Delete(id string) error {
	res, err := r.db.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows 
	}

	return nil
}
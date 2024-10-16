package repository

import (
	"database/sql"
	"time"
    "fmt"

	"github.com/jmoiron/sqlx"
	"github.com/MilaSnetkova/TODO-list/internal/constants"
	"github.com/MilaSnetkova/TODO-list/internal/models"
    "github.com/MilaSnetkova/TODO-list/internal/repeat"
)

// TaskRepository определяет интерфейс для работы с задачами
type TaskRepository interface {
	Create(task *models.Task) (int64, error)
	SearchTasks(filter Filter) ([]models.Task, error)
	GetTaskByID(id string) (*models.Task, error)
	UpdateTask(task *models.Task) error
	Delete(id string) error
	UpdateTaskDate(id string, newDate string) error
	DoneTask(id string, completedAt time.Time) error
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
func (r *TaskRepo) SearchTasks(filter Filter) ([]models.Task, error) {
	var tasks []models.Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE 1=1"
	var params []interface{}

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

	query += " ORDER BY date LIMIT ?"
	params = append(params, constants.Lim)

	rows, err := r.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Получение информации о задаче по ID
func (r *TaskRepo) GetTaskByID(id string) (*models.Task, error) {
	var task models.Task
	err := r.db.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(
		&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil 
		}
		return nil, err
	}

	return &task, nil
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


// Обновление даты задачи
func (r *TaskRepo) UpdateTaskDate(id string, newDate string) error {
	_, err := r.db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", newDate, id)
	return err
}

// Задача выполнена 
func (r *TaskRepo) DoneTask(id string, now time.Time) error {
	
	task, err := r.GetTaskByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch task: %v", err)
	}
	if task == nil {
		return fmt.Errorf("task not found")
	}

	// Если задача одноразовая — удаляем её
	if task.Repeat == "" {
		if err := r.Delete(id); err != nil {
			return fmt.Errorf("failed to delete task: %v", err)
		}
	} else {
		// Если задача повторяющаяся — рассчитываем следующую дату
		nextDate, err := repeat.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("cannot calculate next date: %v", err)
		}
		// Обновляем дату выполнения задачи
		if err := r.UpdateTaskDate(id, nextDate); err != nil {
			return fmt.Errorf("failed to update task date: %v", err)
		}
	}

	return nil
} 
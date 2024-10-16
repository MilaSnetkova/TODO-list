package service

import (
	"fmt"
	"time"
	"log"

	"github.com/MilaSnetkova/TODO-list/internal/constants"
	"github.com/MilaSnetkova/TODO-list/internal/models"
	"github.com/MilaSnetkova/TODO-list/internal/repository"
	"github.com/MilaSnetkova/TODO-list/internal/repeat"
)

type TaskService interface {
	AddTask(task *models.Task) (int64, error)
	GetTasks(search string) ([]models.Task, error)
	UpdateTask(task *models.Task) error
	DeleteTask(id string) error
	TaskDone(id string, now time.Time) error
}

type TaskServiceImpl struct {
	Repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskServiceImpl {
	return &TaskServiceImpl{
		Repo: repo,
	}
}

// Добавление новой задачи
func (s *TaskServiceImpl) AddTask(task *models.Task) (int64, error) {
	now := time.Now()
	var taskDate time.Time

	if task.Date == "" || task.Date == now.Format(constants.DateFormat) {
		taskDate = now
		task.Date = now.Format(constants.DateFormat)
	} else {
		var err error
		taskDate, err = time.Parse(constants.DateFormat, task.Date)
		if err != nil {
			log.Printf("Invalid date format: %v", err)
			return 0, fmt.Errorf("wrong date format")
		}
	}

	if taskDate.Before(now) {
		if task.Repeat == "" || task.Repeat == "d 1" {
			task.Date = now.Format(constants.DateFormat)
		} else {
			nextDate, err := repeat.NextDate(now, taskDate.Format(constants.DateFormat), task.Repeat)
			if err != nil {
				return 0, fmt.Errorf("cannot calculate the next date: %v", err)
			}
			task.Date = nextDate
		}
	}

	id, err := s.Repo.Create(task)
	if err != nil {
		return 0, fmt.Errorf("failed to save task: %v", err)
	}

	return id, nil
}

// Получение задач с фильтрацией
func (s *TaskServiceImpl) GetTasks(search string) ([]models.Task, error) {
	filter := repository.Filter{
		Search: search,
	}

	tasks, err := s.Repo.SearchTasks(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tasks: %v", err)
	}

	return tasks, nil
}


// Обновление задачи
func (s *TaskServiceImpl) UpdateTask(task *models.Task) error {
	if task.ID == "" || task.Title == "" {
		return fmt.Errorf("ID or title field is empty")
	}

	now := time.Now()
	var taskDate time.Time

	if task.Date == "" || task.Date == now.Format(constants.DateFormat) {
		taskDate = now
		task.Date = now.Format(constants.DateFormat)
	} else {
		var err error
		taskDate, err = time.Parse(constants.DateFormat, task.Date)
		if err != nil {
			return fmt.Errorf("wrong date format")
		}
	}

	if taskDate.Before(now) {
		if task.Repeat == "" || task.Repeat == "d 1" {
			task.Date = now.Format(constants.DateFormat)
		} else {
			nextDate, err := repeat.NextDate(now, taskDate.Format(constants.DateFormat), task.Repeat)
			if err != nil {
				return fmt.Errorf("cannot calculate the next date: %v", err)
			}
			task.Date = nextDate
		}
	}
	err := s.Repo.UpdateTask(task)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	return nil
}

// Удаление задачи
func (s *TaskServiceImpl) DeleteTask(id string) error {
	if id == "" {
		return fmt.Errorf("missing task ID")
	}

	err := s.Repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}

// Выполнение задачи
func (s *TaskServiceImpl) TaskDone(id string, now time.Time) error {
    filter := repository.Filter{ID: []string{id}} // Используем фильтр для поиска по ID
    tasks, err := s.Repo.SearchTasks(filter)
    if err != nil {
        return fmt.Errorf("failed to fetch task: %v", err)
    }
    if len(tasks) == 0 {
        return fmt.Errorf("task not found")
    }

    task := tasks[0] // Получаем первую задачу (если она найдена)

    // Если задача не имеет повторений, удалить её
    if task.Repeat == "" {
        if err := s.Repo.Delete(task.ID); err != nil {
            return fmt.Errorf("failed to delete task: %v", err)
        }
    } else {
        // Если задача имеет повторения, вычислить следующую дату
        nextDate, err := repeat.NextDate(now, task.Date, task.Repeat)
        if err != nil {
            return fmt.Errorf("cannot calculate next date: %v", err)
        }

        // Обновляем дату задачи 
        task.Date = nextDate
        if err := s.Repo.UpdateTask(&task); err != nil {
            return fmt.Errorf("failed to update task date: %v", err)
        }
    }

    return nil
}
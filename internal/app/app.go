package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/MilaSnetkova/TODO-list/internal/auth"
	"github.com/MilaSnetkova/TODO-list/internal/config"
	"github.com/MilaSnetkova/TODO-list/internal/db"
	"github.com/MilaSnetkova/TODO-list/internal/handlers"
	"github.com/MilaSnetkova/TODO-list/internal/repository"
	"github.com/MilaSnetkova/TODO-list/internal/service"
)

type App struct {
	Router *chi.Mux
	server *http.Server
	db     *sqlx.DB
}

func NewRouter(version string, serverAddress string, taskService service.TaskSer) *chi.Mux {
	r := chi.NewRouter()

	// Статические файлы из директории web
	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/*", http.StripPrefix("/", fileServer))

	taskHandler := handlers.NewTaskHandler(taskService)

	// Маршрут для обработки правил повторения
	r.Get("/api/nextdate", handlers.NextDateHandler)

	// Маршрут для аутентификации 
	r.Post("/api/signin", auth.SignInHandler)

	r.Group(func(r chi.Router) {
	r.Use(auth.AuthMiddleware)

	r.Get("/api/tasks", taskHandler.GetTasksHandler)
	r.Mount("/api/task", http.HandlerFunc(taskHandler.TaskHandler)) 
	})

	return r
}

func New() (*App, error) {
	cfg, err := config.MustLoad()
	if err != nil {
		return nil, err
	}

	database, err := db.ConnectDB(cfg)
	if err != nil {
		return nil, err
	}

	taskRepo := repository.NewTaskRepo(database)
	taskSer := service.NewTaskService(taskRepo)

	router := NewRouter(cfg.Version, cfg.ServerAddress, taskSer)
	if router == nil {
		return nil, fmt.Errorf("failed to create router")
	}

	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	a := &App{
		Router: router,
		server: server,
		db:     database,
	}

	return a, nil
}

func (a *App) Run() error {
	fmt.Printf("Starting app: Listening on %s\n", a.server.Addr)

	// Запуск сервера
	err := a.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("App stopped")
	}

	return err
}
package database

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/SevvyP/tasks_v1/pkg/model"
	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type TaskDatabase interface {
	GetTasks() (*[]model.Task, error)
	GetTaskByID(id string) (*model.Task, error)
	GetTasksByUserID(userID string) (*[]model.Task, error)
	CreateTask(task model.Task) error
	UpdateTask(task model.Task) error
	DeleteTask(task model.Task) error
}

type PostgresDatabase struct {
	db *sql.DB
}

func NewDatabase(config *PostgresConfig) (*PostgresDatabase, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.Username, url.QueryEscape(config.Password), config.Host, config.Port, config.Database)
	if config.Host == "localhost" {
		connStr += "?sslmode=disable"
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &PostgresDatabase{
		db: db,
	}, nil
}

func (d *PostgresDatabase) GetTasks() (*[]model.Task, error) {
	rows, err := d.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %v", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Body, &task.Completed, &task.Parent, &task.Reminder)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, task)
	}

	return &tasks, nil
}

func (d *PostgresDatabase) GetTaskByID(id string) (*model.Task, error) {
	row := d.db.QueryRow("SELECT * FROM tasks WHERE id = $1", id)

	var task model.Task
	err := row.Scan(&task.ID, &task.UserID, &task.Body, &task.Completed, &task.Parent, &task.Reminder)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	return &task, nil
}

func (d *PostgresDatabase) GetTasksByUserID(userID string) (*[]model.Task, error) {
	rows, err := d.db.Query("SELECT * FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %v", err)
	}
	defer rows.Close()
	tasks := []model.Task{}
	for rows.Next() {
		var task model.Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Body, &task.Completed, &task.Parent, &task.Reminder)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, task)
	}
	return &tasks, nil
}

func (d *PostgresDatabase) CreateTask(task model.Task) error {
	_, err := d.db.Exec("INSERT INTO tasks (id, user_id, body, completed, parent, reminder) VALUES ($1, $2, $3, $4, $5, $6)",
		task.ID, task.UserID, task.Body, task.Completed, task.Parent, task.Reminder)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	return nil
}

func (d *PostgresDatabase) UpdateTask(updatedTask model.Task) error {
	_, err := d.db.Exec("UPDATE tasks SET user_id = $1, body = $2, completed = $3, parent = $4, reminder = $5 WHERE id = $6",
		updatedTask.UserID, updatedTask.Body, updatedTask.Completed, updatedTask.Parent, updatedTask.Reminder, updatedTask.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	return nil
}

func (d *PostgresDatabase) DeleteTask(taskToDelete model.Task) error {
	_, err := d.db.Exec("DELETE FROM tasks WHERE id = $1", taskToDelete.ID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}

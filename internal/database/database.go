package database

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type Task struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Body      string  `json:"body"`
	Completed bool    `json:"completed"`
	Parent    *string `json:"parent"`
	Reminder  *string `json:"reminder"`
}

type TaskDatabase interface {
	GetTasks() (*[]Task, error)
	GetTaskByID(id string) (*Task, error)
	CreateTask(task Task) error
	UpdateTask(task Task) error
	DeleteTask(task Task) error
}

type PostgresDatabase struct {
	db *sql.DB
}

func NewDatabase(config *PostgresConfig) (*PostgresDatabase, error) {
	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.Username, url.QueryEscape(config.Password), config.Host, config.Port, config.Database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return &PostgresDatabase{
		db: db,
	}, nil
}

func (d *PostgresDatabase) GetTasks() (*[]Task, error) {
	rows, err := d.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.UserID, &task.Body, &task.Completed, &task.Parent, &task.Reminder)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %v", err)
		}
		tasks = append(tasks, task)
	}

	return &tasks, nil
}

func (d *PostgresDatabase) GetTaskByID(id string) (*Task, error) {
	row := d.db.QueryRow("SELECT * FROM tasks WHERE id = $1", id)

	var task Task
	err := row.Scan(&task.ID, &task.UserID, &task.Body, &task.Completed, &task.Parent, &task.Reminder)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task: %v", err)
	}

	return &task, nil
}

func (d *PostgresDatabase) CreateTask(task Task) error {
	_, err := d.db.Exec("INSERT INTO tasks (id, user_id, body, completed, parent, reminder) VALUES ($1, $2, $3, $4, $5, $6)",
		task.ID, task.UserID, task.Body, task.Completed, task.Parent, task.Reminder)
	if err != nil {
		return fmt.Errorf("failed to create task: %v", err)
	}

	return nil
}

func (d *PostgresDatabase) UpdateTask(updatedTask Task) error {
	_, err := d.db.Exec("UPDATE tasks SET user_id = $1, body = $2, completed = $3, parent = $4, reminder = $5 WHERE id = $6",
		updatedTask.UserID, updatedTask.Body, updatedTask.Completed, updatedTask.Parent, updatedTask.Reminder, updatedTask.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %v", err)
	}

	return nil
}

func (d *PostgresDatabase) DeleteTask(taskToDelete Task) error {
	_, err := d.db.Exec("DELETE FROM tasks WHERE id = $1", taskToDelete.ID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}

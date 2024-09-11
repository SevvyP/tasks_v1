package database

import (
	"github.com/SevvyP/tasks_v1/pkg/model"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetTasks() (*[]model.Task, error) {
	args := m.Called()
	return args.Get(0).(*[]model.Task), args.Error(1)
}

func (m *MockDatabase) GetTaskByID(id string) (*model.Task, error) {
	args := m.Called(id)
	return args.Get(0).(*model.Task), args.Error(1)
}

func (m *MockDatabase) GetTasksByUserID(userID string) (*[]model.Task, error) {
	args := m.Called(userID)
	return args.Get(0).(*[]model.Task), args.Error(1)
}

func (m *MockDatabase) CreateTask(task model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockDatabase) UpdateTask(task model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockDatabase) DeleteTask(task model.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

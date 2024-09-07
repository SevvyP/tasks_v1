package database

import (
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) GetTasks() (*[]Task, error) {
	args := m.Called()
	return args.Get(0).(*[]Task), args.Error(1)
}

func (m *MockDatabase) GetTaskByID(id string) (*Task, error) {
	args := m.Called(id)
	return args.Get(0).(*Task), args.Error(1)
}

func (m *MockDatabase) CreateTask(task Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockDatabase) UpdateTask(task Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockDatabase) DeleteTask(task Task) error {
	args := m.Called(task)
	return args.Error(0)
}

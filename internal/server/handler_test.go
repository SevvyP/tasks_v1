package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SevvyP/tasks_v1/internal/database"
)

func TestGetTasksHandler(t *testing.T) {
	tests := []struct {
		name           string
		dbResponse     *[]database.Task
		dbError        error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "GetTasks_Success",
			dbResponse: &[]database.Task{
				{ID: "1", Items: []database.Item{{Name: "1", Price: 1}}},
				{ID: "2", Items: []database.Item{{Name: "2", Price: 2}}},
			},
			dbError:        nil,
			expectedStatus: http.StatusOK,
			expectedBody: []database.Task{
				{ID: "1", Items: []database.Item{{Name: "1", Price: 1}}},
				{ID: "2", Items: []database.Item{{Name: "2", Price: 2}}},
			},
		},
		{
			name:           "GetTasks_Error",
			dbResponse:     nil,
			dbError:        fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(database.MockDatabase)
			mockDB.On("GetTasks").Return(tt.dbResponse, tt.dbError)
			resolver := &Resolver{Database: mockDB}

			req, err := http.NewRequest("GET", "/tasks", nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(resolver.GetTasks)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != nil {
				var responseBody []database.Task
				err = json.NewDecoder(rr.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetTaskByIdHandler(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		dbResponse     *database.Task
		dbError        error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "GetTasks_byID_Success",
			id:             "1",
			dbResponse:     &database.Task{ID: "1", Items: []database.Item{{Name: "1", Price: 1}}},
			dbError:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   database.Task{ID: "1", Items: []database.Item{{Name: "1", Price: 1}}},
		},
		{
			name:           "GetTasks_byID_NotFound",
			id:             "1",
			dbResponse:     nil,
			dbError:        nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   nil,
		},
		{
			name:           "GetTasks_byID_Error",
			id:             "1",
			dbResponse:     nil,
			dbError:        fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(database.MockDatabase)
			mockDB.On("GetTaskByID", tt.id).Return(tt.dbResponse, tt.dbError)
			resolver := &Resolver{Database: mockDB}

			req, err := http.NewRequest("GET", "/tasks", nil)
			q := req.URL.Query()
			q.Add("id", tt.id)
			req.URL.RawQuery = q.Encode()

			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(resolver.GetTasks)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != nil {
				var responseBody database.Task
				err = json.NewDecoder(rr.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
			mockDB.AssertExpectations(t)
		})
	}
}

func TestCreateTaskHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           database.Task
		dbResponse     error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "CreateTask_Success",
			body:           database.Task{ID: "1", Items: []database.Item{{Name: "1", Price: 1}}},
			dbResponse:     nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   "Task created successfully",
		},
		{
			name:           "CreateTask_Error",
			body:           database.Task{ID: "1", Items: []database.Item{{Name: "1", Price: 1}}},
			dbResponse:     fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(database.MockDatabase)
			mockDB.On("CreateTask", tt.body).Return(tt.dbResponse)
			resolver := &Resolver{Database: mockDB}

			bodyBytes, _ := json.Marshal(tt.body)
			req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(resolver.CreateTask)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != nil {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestUpdateTaskHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           database.Task
		dbResponse     error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "UpdateTask_Success",
			body:           database.Task{ID: "1", Items: []database.Item{{Name: "1", Price: 1}}},
			dbResponse:     nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Task updated successfully",
		},
		{
			name:           "UpdateTask_Error",
			body:           database.Task{ID: "1", Items: []database.Item{{Name: "1", Price: 1}}},
			dbResponse:     fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(database.MockDatabase)
			mockDB.On("UpdateTask", tt.body).Return(tt.dbResponse)
			resolver := &Resolver{Database: mockDB}

			bodyBytes, _ := json.Marshal(tt.body)
			req, err := http.NewRequest("PUT", "/tasks", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(resolver.UpdateTask)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != nil {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestDeleteTaskHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           database.Task
		dbResponse     error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "DeleteTask_Success",
			body:           database.Task{ID: "1", Items: []database.Item{{Name: "1", Price: 1}}},
			dbResponse:     nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Task deleted successfully",
		},
		{
			name:           "DeleteTask_Error",
			body:           database.Task{ID: "1", Items: []database.Item{{Name: "1", Price: 1}}},
			dbResponse:     fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(database.MockDatabase)
			mockDB.On("DeleteTask", tt.body).Return(tt.dbResponse)
			resolver := &Resolver{Database: mockDB}

			bodyBytes, _ := json.Marshal(tt.body)
			req, err := http.NewRequest("DELETE", "/tasks", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(resolver.DeleteTask)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != nil {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}

			mockDB.AssertExpectations(t)
		})
	}
}

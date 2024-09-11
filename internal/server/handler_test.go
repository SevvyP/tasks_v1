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
	"github.com/SevvyP/tasks_v1/pkg/model"
)

func TestGetTasksHandler(t *testing.T) {
	tests := []struct {
		name           string
		dbResponse     *[]model.Task
		dbError        error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "GetTasks_Success",
			dbResponse: &[]model.Task{
				{ID: "1", Body: "Task 1", Completed: false},
				{ID: "2", Body: "Task 2", Completed: true},
			},
			dbError:        nil,
			expectedStatus: http.StatusOK,
			expectedBody: []model.Task{
				{ID: "1", Body: "Task 1", Completed: false},
				{ID: "2", Body: "Task 2", Completed: true},
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
				var responseBody []model.Task
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
		dbResponse     *model.Task
		dbError        error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "GetTasks_byID_Success",
			id:             "1",
			dbResponse:     &model.Task{ID: "1", Body: "Task 1", Completed: false},
			dbError:        nil,
			expectedStatus: http.StatusOK,
			expectedBody:   model.Task{ID: "1", Body: "Task 1", Completed: false},
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
				var responseBody model.Task
				err = json.NewDecoder(rr.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetTasksByUserIDHandler(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		dbResponse     *[]model.Task
		dbError        error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "GetTasks_byUserID_Success",
			userID: "1",
			dbResponse: &[]model.Task{
				{ID: "1", Body: "Task 1", Completed: false},
				{ID: "2", Body: "Task 2", Completed: true},
			},
			dbError:        nil,
			expectedStatus: http.StatusOK,
			expectedBody: []model.Task{
				{ID: "1", Body: "Task 1", Completed: false},
				{ID: "2", Body: "Task 2", Completed: true},
			},
		},
		{
			name:           "GetTasks_byUserID_Error",
			userID:         "1",
			dbResponse:     nil,
			dbError:        fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(database.MockDatabase)
			mockDB.On("GetTasksByUserID", tt.userID).Return(tt.dbResponse, tt.dbError)
			resolver := &Resolver{Database: mockDB}

			req, err := http.NewRequest("GET", "/tasks", nil)
			q := req.URL.Query()
			q.Add("user_id", tt.userID)
			req.URL.RawQuery = q.Encode()

			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(resolver.GetTasks)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != nil {
				var responseBody []model.Task
				err = json.NewDecoder(rr.Body).Decode(&responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
			mockDB.AssertExpectations(t)
		})
	}
}

func TestGetTasksByUserAndIDError(t *testing.T) {
	resolver := &Resolver{}

	req, err := http.NewRequest("GET", "/tasks", nil)
	q := req.URL.Query()
	q.Add("user_id", "1")
	q.Add("id", "1")
	req.URL.RawQuery = q.Encode()

	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(resolver.GetTasks)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "Cannot query by both id and user\n", rr.Body.String())
}

func TestCreateTaskHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           model.Task
		dbResponse     error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "CreateTask_Success",
			body:           model.Task{ID: "1", Body: "Task 1", Completed: false},
			dbResponse:     nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   "Task created successfully",
		},
		{
			name:           "CreateTask_Error",
			body:           model.Task{ID: "1", Body: "Task 1", Completed: false},
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
		body           model.Task
		dbResponse     error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "UpdateTask_Success",
			body:           model.Task{ID: "1", Body: "Task 1", Completed: false},
			dbResponse:     nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Task updated successfully",
		},
		{
			name:           "UpdateTask_Error",
			body:           model.Task{ID: "1", Body: "Task 1", Completed: false},
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
		body           model.Task
		dbResponse     error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "DeleteTask_Success",
			body:           model.Task{ID: "1", Body: "Task 1", Completed: false},
			dbResponse:     nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "Task deleted successfully",
		},
		{
			name:           "DeleteTask_Error",
			body:           model.Task{ID: "1", Body: "Task 1", Completed: false},
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

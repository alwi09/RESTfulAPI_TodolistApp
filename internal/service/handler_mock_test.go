package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"todolist_gin_gorm/internal/model/dto"
	"todolist_gin_gorm/internal/model/entity"
	"todolist_gin_gorm/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRunMock(t *testing.T) {
	t.Run("TestGetAllSuccess", TestGetAllSuccess)
	t.Run("TestGetAllInternalServerError", TestGetAllInternalServerError)
	t.Run("TestGetAllEmpty", TestGetAllEmpty)
}

func TestGetAllSuccess(t *testing.T) {
	expectedTodo := []entity.Todos{
		{Id: 1, Title: "Todo 1", Description: "Description 1"},
		{Id: 2, Title: "Todo 2", Description: "Description 2"},
	}

	repoMock := mocks.NewRepository(t)
	repoMock.On("GetAll").Return(expectedTodo, nil)

	handler := NewHandlerImpl(repoMock)

	req, err := http.NewRequest("GET", "/find_all_todolist", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/find_all_todolist", handler.GetAllHandlerTodolist)
	router.ServeHTTP(recorder, req)

	var respon dto.TodolistResponseGetAll
	err = json.Unmarshal(recorder.Body.Bytes(), &respon)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "get all todolist successfully", respon.Message)
	assert.Equal(t, http.StatusOK, respon.Status)
	assert.Equal(t, len(expectedTodo), respon.More)
	assert.Equal(t, expectedTodo, respon.Data)
}

func TestGetAllInternalServerError(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	repoMock.On("GetAll").Return(nil, errors.New("internal server error"))

	handler := NewHandlerImpl(repoMock)

	req, err := http.NewRequest("GET", "/find_all_todolist", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/find_all_todolist", handler.GetAllHandlerTodolist)
	router.ServeHTTP(recorder, req)

	var respon dto.ErrorResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &respon)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "internal server error", respon.Message)
	assert.Equal(t, http.StatusInternalServerError, respon.Status)
}

func TestGetAllEmpty(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	repoMock.On("GetAll").Return([]entity.Todos{}, nil)

	handler := NewHandlerImpl(repoMock)

	req, err := http.NewRequest("GET", "/find_all_todolist", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/find_all_todolist", handler.GetAllHandlerTodolist)
	router.ServeHTTP(recorder, req)

	var respon dto.TodolistResponseGetAll
	err = json.Unmarshal(recorder.Body.Bytes(), &respon)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "get all todolist successfully", respon.Message)
	assert.Equal(t, http.StatusOK, respon.Status)
	assert.Equal(t, 0, respon.More)
	assert.Equal(t, []entity.Todos{}, respon.Data)
}

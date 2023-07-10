package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"todolist_gin_gorm/internal/model/dto"
	"todolist_gin_gorm/internal/model/entity"
	"todolist_gin_gorm/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunMock(t *testing.T) {
	t.Run("TestGetAllTodolistsSuccess", TestGetAllTodolistsSuccess)
	t.Run("TestGetAllTodolistsInternalServerError", TestGetAllTodolistsInternalServerError)
	t.Run("TestGetAllTodolistsEmpty", TestGetAllTodolistsEmpty)
}

func TestGetAllTodolistsSuccess(t *testing.T) {
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

func TestGetAllTodolistsInternalServerError(t *testing.T) {
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

func TestGetAllTodolistsEmpty(t *testing.T) {
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

func TestCreateTodolistSuccess(t *testing.T) {
	repoMock := mocks.NewRepository(t)

	newTodo := &entity.Todos{
		Title:       "sholat",
		Description: "sholat tahajud",
		Status:      false,
	}

	repoMock.On("Create", "sholat", "sholat tahajud").Return(newTodo, nil)

	handler := NewHandlerImpl(repoMock)

	point := "/create_todolist"
	router := gin.New()
	router.POST(point, handler.CreateHandlerTodolist)

	reqBody := bytes.NewBufferString(`{"title": "sholat", "description": "sholat tahajud"}`)
	req, err := http.NewRequest(http.MethodPost, point, reqBody)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	responBody, err := io.ReadAll(recorder.Body)
	if err != nil {
		log.Printf("response body: %s\n", responBody)
	}

	require.NoError(t, err)

	var result dto.TodolistResponseCreate
	if err := json.Unmarshal(responBody, &result); err != nil {
		log.Printf("failed to unmarshal JSON response body: %v", err)
	}

	assert.Equal(t, http.StatusCreated, result.Status)
	assert.Equal(t, "create todolist successfully", result.Message)
	assert.Equal(t, *newTodo, result.Data)
}

func TestCreateTodolistInvalidValidation(t *testing.T) {
	repoMock := mocks.NewRepository(t)
	handler := NewHandlerImpl(repoMock)

	expectedErrors := errors.New("invalid input validation")

	point := "/create_todolist"

	reqBody := bytes.NewBufferString(`{"title": "", "description": ""}`)
	req, err := http.NewRequest(http.MethodPost, point, reqBody)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	ctx, router := gin.CreateTestContext(recorder)
	router.POST(point, func(context *gin.Context) {
		handler.CreateHandlerTodolist(context)
	})

	ctx.Request = req
	router.ServeHTTP(recorder, req)

	responBody, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)

	var ErrorResponse dto.ErrorResponse
	err = json.Unmarshal(responBody, &ErrorResponse)
	require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, ErrorResponse.Status)
	assert.Equal(t, expectedErrors.Error(), ErrorResponse.Message)
}

func TestCreateTodolistInternalServerError(t *testing.T) {
	repoMock := mocks.NewRepository(t)

	handler := NewHandlerImpl(repoMock)

	expectedErrors := errors.New("internal server error")
	point := "/create_todolist"

	repoMock.On("Create", "Sholat", "Sholat Tahajud").Return(nil, expectedErrors)

	reqBody := bytes.NewBufferString(`{"title": "Sholat", "description": "Sholat Tahajud"}`)
	req, err := http.NewRequest(http.MethodPost, point, reqBody)
	require.NoError(t, err)

	recorder := httptest.NewRecorder()
	ctx, router := gin.CreateTestContext(recorder)
	router.POST(point, handler.CreateHandlerTodolist)

	ctx.Request = req
	router.ServeHTTP(recorder, req)

	responBody, err := io.ReadAll(recorder.Body)
	require.NoError(t, err)

	var ErrorResponse dto.ErrorResponse
	err = json.Unmarshal(responBody, &ErrorResponse)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, ErrorResponse.Status)
	assert.Equal(t, expectedErrors.Error(), ErrorResponse.Message)
}

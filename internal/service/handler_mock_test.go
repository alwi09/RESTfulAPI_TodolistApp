package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"todolist_gin_gorm/internal/model/entity"
	"todolist_gin_gorm/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGetAllSuccess(t *testing.T) {
	repoMock := &mocks.Repository{}
	handler := NewHandlerImpl(repoMock)

	expectedTodo := []entity.Todos{
		{Id: 1, Title: "Todo 1", Description: "Description 1"},
		{Id: 2, Title: "Todo 2", Description: "Description 2"},
	}

	repoMock.On("GetAll").Return(expectedTodo, nil)

	req, err := http.NewRequest("GET", "/todos", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	expectedResponse := []byte(`[{"id": 1, "title": "Todo 1", "description": "Description 1"}, {"id": 2, "title": "Todo 2", "description": "Description 2"}]`)
	assert.Equal(t, expectedResponse, recorder.Body.Bytes())

	repoMock.AssertCalled(t, "GetAll")
}

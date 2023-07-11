package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"todolist_gin_gorm/internal/model/dto"
	"todolist_gin_gorm/internal/model/entity"
	"todolist_gin_gorm/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRunTableDriven(t *testing.T) {
	t.Run("TestTableDrivenGetAllTodolist", TestTableDrivenGetAllTodolist)
}

func TestTableDrivenGetAllTodolist(t *testing.T) {

	testCase := []struct {
		name               string
		expectedStatusCode int
		mockTodolist       []entity.Todos
		mockErr            error
		expextedResponse   dto.TodolistResponseGetAll
	}{
		{
			name:               "Success",
			expectedStatusCode: http.StatusOK,
			mockTodolist: []entity.Todos{
				{Id: 1, Title: "title 1", Description: "description 1", Status: false},
				{Id: 1, Title: "title 1", Description: "description 1", Status: false},
			},
			mockErr: nil,
			expextedResponse: dto.TodolistResponseGetAll{
				Status:  http.StatusOK,
				Message: "get all todolist successfully",
				More:    2,
				Data: []entity.Todos{
					{Id: 1, Title: "title 1", Description: "description 1", Status: false},
					{Id: 1, Title: "title 1", Description: "description 1", Status: false},
				},
			},
		},
		{
			name:               "internal server error",
			expectedStatusCode: http.StatusInternalServerError,
			mockTodolist:       []entity.Todos{},
			mockErr:            errors.New("internal server error"),
			expextedResponse: dto.TodolistResponseGetAll{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
				More:    0,
				Data:    []entity.Todos(nil),
			},
		},
		{
			name:               "empty",
			expectedStatusCode: http.StatusOK,
			mockTodolist:       []entity.Todos{},
			mockErr:            nil,
			expextedResponse: dto.TodolistResponseGetAll{
				Status:  http.StatusOK,
				Message: "get all todolist successfully",
				More:    0,
				Data:    []entity.Todos{},
			},
		},
	}

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			repo.On("GetAll").Return(test.mockTodolist, test.mockErr)

			handler := NewHandlerImpl(repo)

			req, err := http.NewRequest("GET", "/find_all_todolist", nil)
			if err != nil {
				t.Fatal(err)
			}

			recorder := httptest.NewRecorder()
			router := gin.Default()
			router.GET("/find_all_todolist", handler.GetAllHandlerTodolist)
			router.ServeHTTP(recorder, req)

			assert.Equal(t, test.expectedStatusCode, recorder.Code)

			var respon dto.TodolistResponseGetAll
			err = json.Unmarshal(recorder.Body.Bytes(), &respon)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Printf("Type: %T\n", test.expextedResponse.Data)
			fmt.Printf("Type: %T\n", respon.Data)

			assert.Equal(t, test.expectedStatusCode, respon.Status)
			assert.Equal(t, test.expextedResponse.Message, respon.Message)
			assert.Equal(t, test.expextedResponse.More, respon.More)
			assert.Equal(t, test.expextedResponse.Data, respon.Data)
		})
	}
}

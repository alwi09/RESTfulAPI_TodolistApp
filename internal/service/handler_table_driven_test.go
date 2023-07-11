package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func TestRunTableDriven(t *testing.T) {
	t.Run("TestTableDrivenGetAllTodolist", TestTableDrivenGetAllTodolist)
	t.Run("TestTableDrivenCreateTodolist", TestTableDrivenCreateTodolist)
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

func TestTableDrivenCreateTodolist(t *testing.T) {

	testCase := []struct {
		name           string
		bodyRequest    string
		mock           func(repository *mocks.Repository)
		expectedStatus int
		expectedData   entity.Todos
		expectedError  string
	}{
		{
			name:        "success",
			bodyRequest: `{"title": "sholat", "description": "sholat tahajud"}`,
			mock: func(mock *mocks.Repository) {
				newTodo := &entity.Todos{
					Title:       "sholat",
					Description: "sholat tahajud",
					Status:      false,
				}
				mock.On("Create", "sholat", "sholat tahajud").Return(newTodo, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedData: entity.Todos{
				Title:       "sholat",
				Description: "sholat tahajud",
				Status:      false,
			},
			expectedError: "",
		},
		{
			name:           "invalid validation",
			bodyRequest:    `{"title": "", "description": ""}`,
			mock:           func(mock *mocks.Repository) {},
			expectedStatus: http.StatusBadRequest,
			expectedData:   entity.Todos{},
			expectedError:  "invalid input validation",
		},
		{
			name:        "internal server error",
			bodyRequest: `{"title": "sholat", "description": "sholat tahajud"}`,
			mock: func(mock *mocks.Repository) {
				expectedErr := errors.New("internal server error")
				mock.On("Create", "sholat", "sholat tahajud").Return(nil, expectedErr)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedData:   entity.Todos{},
			expectedError:  "internal server error",
		},
	}

	for _, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := mocks.NewRepository(t)
			test.mock(mockRepo)
			handler := NewHandlerImpl(mockRepo)

			point := "/create_todolist"

			reqBody := bytes.NewBufferString(test.bodyRequest)
			req, err := http.NewRequest(http.MethodPost, point, reqBody)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			ctx, router := gin.CreateTestContext(recorder)
			router.POST(point, handler.CreateHandlerTodolist)

			ctx.Request = req
			router.ServeHTTP(recorder, req)

			responBody, err := io.ReadAll(recorder.Body)
			require.NoError(t, err)

			if test.expectedError != "" {
				var errResponse dto.ErrorResponse
				err = json.Unmarshal(responBody, &errResponse)
				require.NoError(t, err)
				assert.Equal(t, test.expectedError, errResponse.Message)
			} else {
				var result dto.TodolistResponseCreate
				err = json.Unmarshal(responBody, &result)
				require.NoError(t, err)
				assert.Equal(t, test.expectedData, result.Data)
			}

			assert.Equal(t, test.expectedStatus, recorder.Code)
		})
	}
}

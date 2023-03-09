package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"todolist_gin_gorm/cmd/router"
	"todolist_gin_gorm/internal/database"
	"todolist_gin_gorm/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open("alwi09:alwiirfani091199@tcp(localhost:3306)/todolist_test"))
	if err != nil {
		logrus.Panic(err)
	}

	logrus.Info("connect to database")
	return db, nil
}

func SetupTestRouter(db *gorm.DB) *gin.Engine {
	TodolistRepository := database.NewTodoRepository(db)
	TodolistService := service.NewHandlerImpl(TodolistRepository)
	TodolistRouteBuilder := router.NewRouteBuilder(TodolistService)
	TodolistRouteInit := TodolistRouteBuilder.RouteInit()

	return TodolistRouteInit
}

func TruncateTodolist(DB *gorm.DB) {
	DB.Exec("TRUNCATE todos")
}

func TestTodolisService(t *testing.T) {
	t.Run("TestCreateTodolistSuccess", TestCreateTodolistSuccess)
	t.Run("TestCreateTodolistFailedValidation", TestCreateTodolistFailedValidation)
	t.Run("TestUpdateTodolistSuccess", TestUpdateTodolistSuccess)
	t.Run("TestUpdateTodolistFailed", TestUpdateTodolistFailed)
	t.Run("TestGetIDTodolistSuccess", TestGetIDTodolistSuccess)
	t.Run("TestGetIDTodolistFailed", TestGetIDTodolistFailedNotFound)
	t.Run("TestDeleteTodolistSuccess", TestDeleteTodolistSuccess)
	t.Run("TestDeleteTodolistFailed", TestDeleteTodolistFailedNotFound)
	t.Run("TestListTodoSuccess", TestGetAllTodolistSuccess)
	t.Run("TestUnauthorizod", TestUnauthorizod)

}

func TestCreateTodolistSuccess(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		logrus.Fatal(err)
	}

	TruncateTodolist(db)
	router := SetupTestRouter(db)

	requestBody := strings.NewReader(`{"title": "ibadah", "description": "sholat shubuh", "status": false}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:1234/api/create_todolist", requestBody)
	request.Header.Add("Authorization", "secret_lock")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(http.StatusBadRequest, "bad request", err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(responseBody)

	assert.Equal(t, http.StatusCreated, int(responseBody["status"].(float64)))
	assert.Equal(t, "create todolist successfully", responseBody["message"])
	assert.Equal(t, "ibadah", responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, "sholat shubuh", responseBody["data"].(map[string]interface{})["description"])
	assert.Equal(t, false, responseBody["data"].(map[string]interface{})["status"])
}

func TestCreateTodolistFailedValidation(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		logrus.Fatal(err)
	}

	TruncateTodolist(db)
	router := SetupTestRouter(db)

	requestBody := strings.NewReader(`{"title": "", "description": "", "status":}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:1234/api/create_todolist", requestBody)
	request.Header.Add("Authorization", "secret_lock")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(http.StatusBadRequest, "bad request", err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(responseBody)

	assert.Equal(t, http.StatusBadRequest, int(responseBody["status"].(float64)))
}

func TestUpdateTodolistSuccess(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		logrus.Fatal(err)
	}

	TruncateTodolist(db)
	router := SetupTestRouter(db)

	tx := db.Begin()
	todolistRepository := database.NewTodoRepository(db)
	todolist, err := todolistRepository.Create("ibadah", "sholat tahajud")
	if err != nil {
		logrus.Error(err)
	}

	tx.Commit()

	requestBody := strings.NewReader(`{"title": "ibadah", "description": "sholat shubuh", "status": true}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:1234/api/update_todolist/"+strconv.Itoa(int(todolist.Id)), requestBody)
	request.Header.Add("Authorization", "secret_lock")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(http.StatusBadRequest, "bad request", err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(responseBody)

	assert.Equal(t, http.StatusOK, int(responseBody["status"].(float64)))
	assert.Equal(t, "update todolist successfully", responseBody["message"])
	// assert.Equal(t, todolist.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "ibadah", responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, "sholat shubuh", responseBody["data"].(map[string]interface{})["description"])
	assert.Equal(t, true, responseBody["data"].(map[string]interface{})["status"])

}

func TestUpdateTodolistFailed(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		logrus.Fatal(err)
	}

	TruncateTodolist(db)
	router := SetupTestRouter(db)

	tx := db.Begin()
	todolistRepository := database.NewTodoRepository(db)
	todolist, err := todolistRepository.Create("ibadah", "sholat tahajud")
	if err != nil {
		logrus.Error(err)
	}

	tx.Commit()

	requestBody := strings.NewReader(`{"title": "", "description": "", "status":}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:1234/api/update_todolist/"+strconv.Itoa(int(todolist.Id)), requestBody)
	request.Header.Add("Authorization", "secret_lock")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(http.StatusBadRequest, "bad request", err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(responseBody)

	assert.Equal(t, http.StatusBadRequest, int(responseBody["status"].(float64)))
}

func TestGetIDTodolistSuccess(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		logrus.Fatal(err)
	}

	TruncateTodolist(db)
	router := SetupTestRouter(db)

	tx := db.Begin()
	todolistRepository := database.NewTodoRepository(db)
	todolist, err := todolistRepository.Create("ibadah", "sholat tahajud")
	if err != nil {
		logrus.Error(err)
	}

	tx.Commit()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:1234/api/find_by_id_todolist/1", nil)
	request.Header.Add("Authorization", "secret_lock")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(http.StatusBadRequest, "bad request", err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(responseBody)

	assert.Equal(t, http.StatusOK, int(responseBody["status"].(float64)))
	assert.Equal(t, "get todolist by id successfully", responseBody["message"])
	// assert.Equal(t, todolist.Id, int64(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, todolist.Title, responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, todolist.Description, responseBody["data"].(map[string]interface{})["description"])
	assert.Equal(t, todolist.Status, responseBody["data"].(map[string]interface{})["status"])

}

func TestGetIDTodolistFailedNotFound(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		logrus.Fatal(err)
	}

	TruncateTodolist(db)
	router := SetupTestRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:1234/api/find_by_id_todolist/404", nil)
	request.Header.Add("Authorization", "secret_lock")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(http.StatusBadRequest, "bad request", err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(responseBody)

	assert.Equal(t, http.StatusNotFound, int(responseBody["status"].(float64)))
	assert.Equal(t, "todolist by id not found", responseBody["message"])
}

func TestDeleteTodolistSuccess(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		logrus.Fatal(err)
	}

	TruncateTodolist(db)
	router := SetupTestRouter(db)

	tx := db.Begin()
	todolistRepository := database.NewTodoRepository(db)
	todolist, err := todolistRepository.Create("ibadah", "sholat tahajud")
	if err != nil {
		logrus.Error(err)
	}

	tx.Commit()

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:1234/api/delete_todolist/"+strconv.Itoa(int(todolist.Id)), nil)
	request.Header.Add("Authorization", "secret_lock")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(http.StatusBadRequest, "bad request", err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(responseBody)

	assert.Equal(t, http.StatusOK, int(responseBody["status"].(float64)))
	assert.Equal(t, "delete todolist successfully", responseBody["message"])
}

func TestDeleteTodolistFailedNotFound(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		logrus.Fatal(err)
	}

	TruncateTodolist(db)
	router := SetupTestRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:1234/api/delete_todolist/404", nil)
	request.Header.Add("Authorization", "secret_lock")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusNotFound, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(http.StatusBadRequest, "bad request", err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(responseBody)

	assert.Equal(t, http.StatusNotFound, int(responseBody["status"].(float64)))
	assert.Equal(t, "id not found", responseBody["message"])

}

func TestGetAllTodolistSuccess(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		logrus.Fatal(err)
	}

	TruncateTodolist(db)
	router := SetupTestRouter(db)

	tx := db.Begin()
	todolistRepository := database.NewTodoRepository(db)
	todolist1, _ := todolistRepository.Create("ibadah", "sholat tahajud")
	todolist2, _ := todolistRepository.Create("ibadah", "sholat shubuh")

	tx.Commit()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:1234/api/find_all_todolist", nil)
	request.Header.Add("Authorization", "secret_lock")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(http.StatusBadRequest, "bad request", err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(responseBody)

	assert.Equal(t, http.StatusOK, int(responseBody["status"].(float64)))
	assert.Equal(t, "get all todolist successfully", responseBody["message"])

	var todolists = responseBody["data"].([]interface{})

	todolistResponse1 := todolists[0].(map[string]interface{})
	todolistResponse2 := todolists[1].(map[string]interface{})

	assert.Equal(t, todolist1.Id, int64(todolistResponse1["id"].(float64)))
	assert.Equal(t, todolist1.Title, todolistResponse1["title"])
	assert.Equal(t, todolist1.Description, todolistResponse1["description"])

	assert.Equal(t, todolist2.Id, int64(todolistResponse2["id"].(float64)))
	assert.Equal(t, todolist2.Title, todolistResponse2["title"])
	assert.Equal(t, todolist2.Description, todolistResponse2["description"])
}

func TestUnauthorizod(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		logrus.Fatal(err)
	}

	TruncateTodolist(db)
	router := SetupTestRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:1234/api/find_all_todolist", nil)
	request.Header.Add("Authorization", "")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		logrus.Error(http.StatusBadRequest, "bad request", err)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		logrus.Error(err)
	}
	fmt.Println(responseBody)

	assert.Equal(t, http.StatusUnauthorized, int(responseBody["status"].(float64)))
	assert.Equal(t, "Unauthorized", responseBody["message"])
}

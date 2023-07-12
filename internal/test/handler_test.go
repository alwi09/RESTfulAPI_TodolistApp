package test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open("alwi09:alwiirfani091199@tcp(localhost:3306)/todolist_test"))
	if err != nil {
		panic("cannot connect to database")
	}

	logrus.Info("connect to database successfully")
	return db, nil
}

func setupRouter(DB *gorm.DB) *gin.Engine {

	// initialize repositories
	todolistRepository := database.NewTodoRepository(DB)

	// initialize service
	todolistHandler := service.NewHandlerImpl(todolistRepository)

	// initialize router
	routeBuilder := router.NewRouteBuilder(todolistHandler)
	routerInit := routeBuilder.RouteInit()

	return routerInit
}

func truncateTodolist(DB *gorm.DB) {
	DB.Exec("TRUNCATE todos")
}

func TestRunHandler(t *testing.T) {
	t.Run("TestCreateTodolistSuccess", TestCreateTodolistSuccess)
	t.Run("TestCreateTodolistFailedBadRequest", TestCreateTodolistFailedBadRequest)
}

func TestCreateTodolistSuccess(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		log.Print(err)
	}

	truncateTodolist(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"title": "sholat", "description": "sholat tahajud"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:1234/api/create_todolist", requestBody)
	request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBvbmRva3Byb2dyYW1tZXJAZ21haWwuY29tIiwiZXhwIjoxNjg5MjYxMTQ2fQ.yNr_tCLqZIqXtW0PO3N1kDdwT3-IEWWKoCD6nCPUWY8")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusCreated, int(responseBody["status"].(float64)))
	assert.Equal(t, "create todolist successfully", responseBody["message"])
	assert.Equal(t, "sholat", responseBody["data"].(map[string]interface{})["title"])
	assert.Equal(t, "sholat tahajud", responseBody["data"].(map[string]interface{})["description"])
}

func TestCreateTodolistFailedBadRequest(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		log.Print(err)
	}

	truncateTodolist(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"title": "", "description": ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:1234/api/create_todolist", requestBody)
	request.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InBvbmRva3Byb2dyYW1tZXJAZ21haWwuY29tIiwiZXhwIjoxNjg5MjYxMTQ2fQ.yNr_tCLqZIqXtW0PO3N1kDdwT3-IEWWKoCD6nCPUWY8")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	json.Unmarshal(body, &responseBody)

	assert.Equal(t, http.StatusBadRequest, int(responseBody["status"].(float64)))
	assert.Equal(t, "invalid input validation", responseBody["message"])
}

func TestUpdateTodolistSuccess(t *testing.T) {

}

func TestUpdateTodolistFailed(t *testing.T) {

}

func TestGetTodolistByIDSuccess(t *testing.T) {

}

func TestGetTodolistByIDFailed(t *testing.T) {

}

func TestDeleteTodolistSuccess(t *testing.T) {

}

func TestDeleteTodolistFailed(t *testing.T) {

}

func TestListTodolistsSuccess(t *testing.T) {

}

func TestUnauthorized(t *testing.T) {

}
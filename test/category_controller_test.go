package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gorestfulapi/app"
	"gorestfulapi/controller"
	"gorestfulapi/helper"
	"gorestfulapi/middleware"
	"gorestfulapi/model/domain"
	"gorestfulapi/repository"
	"gorestfulapi/service"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-sql-driver/mysql"
)

func setupTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/go_mysql_test")

	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setupRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)
	
	router := app.NewRouter(categoryController)

	return middleware.NewAuthMiddleware(router)
}

func truncateCategory(db *sql.DB) {
	db.Exec("TRUNCATE categories")
}

func TestCreateCategorySuccess(t *testing.T) {
	db := setupTestDB()

	truncateCategory(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name": "Gadget"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "AUTH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
}

func TestCreateCategoryFailed(t *testing.T) {
	db := setupTestDB()

	truncateCategory(db)
	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name": ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "AUTH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
}

func TestUpdateCategorySuccess(t *testing.T) {
	db := setupTestDB()

	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Fashion",
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name": "Gadget"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/" + strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "AUTH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]any)["id"].(float64)))
}

func TestUpdateCategoryFailed(t *testing.T) {
	db := setupTestDB()

	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Fashion",
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name": ""}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/" + strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "AUTH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t,  400, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)

	assert.Equal(t, 400, int(responseBody["code"].(float64)))
}

func TestGetCategorySuccess(t *testing.T) {
	db := setupTestDB()

	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Fashion",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/" + strconv.Itoa(category.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "AUTH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]any)["id"].(float64)))
}

func TestGetCategoryFailed(t *testing.T) {
	db := setupTestDB()

	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	_ = categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Fashion",
	})
	tx.Commit()

	router := setupRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/" + "8", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "AUTH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	// assert.Equal(t, category.Id, int(responseBody["data"].(map[string]any)["id"].(float64)))
}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setupTestDB()

	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Fashion",
	})
	tx.Commit()

	router := setupRouter(db)

	requestBody := strings.NewReader(`{"name": "Gadget"}`)
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/" + strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "AUTH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	// assert.Equal(t, category.Id, int(responseBody["data"].(map[string]any)["id"].(float64)))
}

func TestDeleteCategoryFailed(t *testing.T) {
	db := setupTestDB()

	truncateCategory(db)

	// tx, _ := db.Begin()
	// categoryRepository := repository.NewCategoryRepository()
	// category := categoryRepository.Save(context.Background(), tx, domain.Category{
	// 	Name: "Fashion",
	// })
	// tx.Commit()

	router := setupRouter(db)

	// requestBody := strings.NewReader(`{"name": "Gadget"}`)
	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/" + "404", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "AUTH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 404, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)

	assert.Equal(t, 404, int(responseBody["code"].(float64)))
	// assert.Equal(t, category.Id, int(responseBody["data"].(map[string]any)["id"].(float64)))
}

func TestListCategoriesSuccess(t *testing.T) {
	db := setupTestDB()

	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Fashion",
	})
	tx.Commit()

	router := setupRouter(db)

	// requestBody := strings.NewReader(`{"name": "Gadget"}`)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "AUTH")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 200, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)
	fmt.Println(category)

	var categories = responseBody["data"].([]interface{})

	assert.Equal(t, 200, int(responseBody["code"].(float64)))
	assert.Equal(t, category.Id, int(categories[0].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, category.Name, categories[0].(map[string]interface{})["name"])
}

func TestUnauthorized(t *testing.T) {
	db := setupTestDB()

	truncateCategory(db)

	tx, _ := db.Begin()
	categoryRepository := repository.NewCategoryRepository()
	category := categoryRepository.Save(context.Background(), tx, domain.Category{
		Name: "Fashion",
	})
	tx.Commit()

	router := setupRouter(db)

	// requestBody := strings.NewReader(`{"name": "Gadget"}`)
	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "WRONG")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, 401, response.StatusCode)

	body, _ := io.ReadAll(response.Body)

	var responseBody map[string]any
	json.Unmarshal(body, &responseBody)
	fmt.Println(responseBody)
	fmt.Println(category)

	// var categories = responseBody["data"].([]interface{})

	assert.Equal(t, 401, int(responseBody["code"].(float64)))
	// assert.Equal(t, category.Id, int(categories[0].(map[string]interface{})["id"].(float64)))
	// assert.Equal(t, category.Name, categories[0].(map[string]interface{})["name"])
}
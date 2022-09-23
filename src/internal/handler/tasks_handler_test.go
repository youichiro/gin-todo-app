package handler

import (
	"example/web-service-gin/internal/client"
	"fmt"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}

func setup() {
	fmt.Println("setup")
}

func teardown() {
	fmt.Println("teardown")
}

func TestTaskHandlerIndex(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	client.DB = mockDB

	rows := mock.NewRows([]string{"id", "title", "done"})
	rows.AddRow(1, "dummy task1", false)
	rows.AddRow(2, "dummy task2", true)
	query := regexp.QuoteMeta(`SELECT "tasks".* FROM "tasks"`)
	mock.ExpectQuery(query).WillReturnRows(rows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/tasks", nil)

	TaskHander{}.Index(c)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, 200, w.Code)
	expected := `[
		{"id": 1, "title": "dummy task1", "done": false, "created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"},
		{"id": 2, "title": "dummy task2", "done": true, "created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}
	]`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestTaskHandlerShow(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()
	client.DB = mockDB

	rows := mock.NewRows([]string{"id", "title", "done"}).AddRow(3, "dummy task3", false)
	query := regexp.QuoteMeta(`select * from "tasks" where "id"=$1`)
	mock.ExpectQuery(query).WillReturnRows(rows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/tasks/3", nil)

	TaskHander{}.Show(c)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, 200, w.Code)
	expected := `{"id": 3, "title": "dummy task3", "done": false, "created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}`
	assert.JSONEq(t, expected, w.Body.String())
}

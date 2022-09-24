package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"example/web-service-gin/internal/client"
	"example/web-service-gin/internal/models"
	"fmt"
	"io"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func setup() {
	fmt.Println("setup")
}

func teardown() {
	fmt.Println("teardown")
}

func InitMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	client.DB = mockDB
	return mockDB, mock
}

func TestTaskHandlerIndex(t *testing.T) {
	tests := []struct {
		name         string
		expectStatus int
		expectBody   string
		expectError  bool
	}{
		{
			name:         "正常系",
			expectStatus: 200,
			expectBody: `[
				{"id": 1, "title": "dummy task1", "done": false, "created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"},
				{"id": 2, "title": "dummy task2", "done": true, "created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}
			]`,
			expectError: false,
		},
		{
			name:         "異常系",
			expectStatus: 404,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockDB, mock := InitMockDB(t)
			defer mockDB.Close()
			rows := mock.NewRows([]string{"id", "title", "done"})
			rows.AddRow(1, "dummy task1", false)
			rows.AddRow(2, "dummy task2", true)
			query := regexp.QuoteMeta(`SELECT "tasks".* FROM "tasks"`)
			if tt.expectError {
				mock.ExpectQuery(query).WillReturnError(fmt.Errorf("error"))
			} else {
				mock.ExpectQuery(query).WillReturnRows(rows)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/tasks", nil)

			TaskHander{}.Index(c)

			assert.Equal(t, tt.expectStatus, w.Code)
			if tt.expectError == false {
				assert.JSONEq(t, tt.expectBody, w.Body.String())
			}
		})
	}
}

func TestTaskHandlerShow(t *testing.T) {
	tests := []struct {
		name         string
		expectStatus int
		expectBody   string
		expectError  bool
	}{
		{
			name:         "正常系",
			expectStatus: 200,
			expectBody:   `{"id": 3, "title": "dummy task3", "done": false, "created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}`,
			expectError:  false,
		},
		{
			name:         "異常系",
			expectStatus: 404,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockDB, mock := InitMockDB(t)
			defer mockDB.Close()
			rows := mock.NewRows([]string{"id", "title", "done"}).AddRow(3, "dummy task3", false)
			query := regexp.QuoteMeta(`select * from "tasks" where "id"=$1`)
			if tt.expectError {
				mock.ExpectQuery(query).WillReturnError(fmt.Errorf("error"))
			} else {
				mock.ExpectQuery(query).WillReturnRows(rows)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/tasks/3", nil)

			TaskHander{}.Show(c)

			assert.Equal(t, tt.expectStatus, w.Code)
			if tt.expectError == false {
				assert.JSONEq(t, tt.expectBody, w.Body.String())
			}
		})
	}
}

func TestTaskHandlerCreate(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		mockDB, mock := InitMockDB(t)
		defer mockDB.Close()
		rows := mock.NewRows([]string{"id", "done"}).AddRow(0, false)
		query := regexp.QuoteMeta(`INSERT INTO "tasks" ("title","created_at","updated_at") VALUES ($1,$2,$3) RETURNING "id","done"`)
		mock.ExpectQuery(query).WillReturnRows(rows)

		d := &createParams{Title: "dummy insert task"}
		jsonString, err := json.Marshal(d)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonString))
		c.Request.Header.Set("Content-Type", "application/json")

		TaskHander{}.Create(c)

		assert.Equal(t, 201, w.Code)

		var task models.Task
		body, _ := io.ReadAll(w.Body)
		err = json.Unmarshal(body, &task)
		assert.NoError(t, err)
		opt := cmpopts.IgnoreFields(task, "CreatedAt", "UpdatedAt")
		expectBody := models.Task{ID: 0, Title: "dummy insert task", Done: false}
		assert.Empty(t, cmp.Diff(expectBody, task, opt))
	})

	t.Run("異常系_INSERTに失敗した場合", func(t *testing.T) {
		mockDB, mock := InitMockDB(t)
		defer mockDB.Close()
		query := regexp.QuoteMeta(`INSERT INTO "tasks" ("title","created_at","updated_at") VALUES ($1,$2,$3) RETURNING "id","done"`)
		mock.ExpectQuery(query).WillReturnError(fmt.Errorf("error"))

		d := &createParams{Title: "dummy insert task"}
		jsonString, err := json.Marshal(d)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonString))
		c.Request.Header.Set("Content-Type", "application/json")

		TaskHander{}.Create(c)

		assert.Equal(t, 500, w.Code)
	})

	t.Run("異常系_リクエストパラメーターが間違えている場合", func(t *testing.T) {
		mockDB, _ := InitMockDB(t)
		defer mockDB.Close()

		d := struct{ InvalidTitle string }{InvalidTitle: "invalid title"}
		jsonString, err := json.Marshal(d)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonString))
		c.Request.Header.Set("Content-Type", "application/json")

		TaskHander{}.Create(c)

		assert.Equal(t, 400, w.Code)
	})
}

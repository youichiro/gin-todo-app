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

var cmpOption cmp.Option

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func setup() {
	fmt.Println("setup")
	cmpOption = cmpopts.IgnoreFields(models.Task{}, "CreatedAt", "UpdatedAt")
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

func CreateTestContext(method string, path string, jsonString string) (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	if jsonString != "" {
		c.Request = httptest.NewRequest(method, path, bytes.NewBuffer([]byte(jsonString)))
	} else {
		c.Request = httptest.NewRequest(method, path, nil)
	}
	c.Request.Header.Set("Content-Type", "application/json")

	return w, c
}

func TestTaskHandlerIndex(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		mockDB, mock := InitMockDB(t)

		rows := mock.NewRows([]string{"id", "title", "done"})
		rows.AddRow(0, "dummy task1", false)
		rows.AddRow(1, "dummy task2", true)
		query := regexp.QuoteMeta(`SELECT "tasks".* FROM "tasks"`)
		mock.ExpectQuery(query).WillReturnRows(rows)

		w, c := CreateTestContext("GET", "/tasks", "")
		TaskHander{}.Index(c)

		assert.Equal(t, 200, w.Code)

		var tasks []models.Task
		body, _ := io.ReadAll(w.Body)
		err := json.Unmarshal(body, &tasks)
		assert.NoError(t, err)
		expectBodyFirst := models.Task{ID: 0, Title: "dummy task1", Done: false}
		expectBodySecond := models.Task{ID: 1, Title: "dummy task2", Done: true}
		assert.Equal(t, 2, len(tasks))
		assert.Empty(t, cmp.Diff(expectBodyFirst, tasks[0], cmpOption))
		assert.Empty(t, cmp.Diff(expectBodySecond, tasks[1], cmpOption))

		t.Cleanup(func() {
			mockDB.Close()
		})
	})

	t.Run("異常系_SELECTに失敗する場合", func(t *testing.T) {
		t.Parallel()
		mockDB, mock := InitMockDB(t)

		query := regexp.QuoteMeta(`SELECT "tasks".* FROM "tasks"`)
		mock.ExpectQuery(query).WillReturnError(fmt.Errorf("error"))

		w, c := CreateTestContext("GET", "/tasks", "")
		TaskHander{}.Index(c)

		assert.Equal(t, 404, w.Code)

		t.Cleanup(func() {
			mockDB.Close()
		})
	})
}

func TestTaskHandlerShow(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		mockDB, mock := InitMockDB(t)
		rows := mock.NewRows([]string{"id", "title", "done"}).AddRow(3, "dummy task3", false)
		query := regexp.QuoteMeta(`select * from "tasks" where "id"=$1`)
		mock.ExpectQuery(query).WillReturnRows(rows)

		w, c := CreateTestContext("GET", "/tasks/3", "")
		TaskHander{}.Show(c)

		assert.Equal(t, 200, w.Code)
		var task models.Task
		body, _ := io.ReadAll(w.Body)
		err := json.Unmarshal(body, &task)
		assert.NoError(t, err)
		opt := cmpopts.IgnoreFields(models.Task{}, "CreatedAt", "UpdatedAt")
		expectBody := models.Task{ID: 3, Title: "dummy task3", Done: false}
		assert.Empty(t, cmp.Diff(expectBody, task, opt))

		t.Cleanup(func() {
			mockDB.Close()
		})
	})

	t.Run("異常系_レコードが存在しない場合", func(t *testing.T) {
		t.Parallel()
		mockDB, mock := InitMockDB(t)
		query := regexp.QuoteMeta(`select * from "tasks" where "id"=$1`)
		mock.ExpectQuery(query).WillReturnError(fmt.Errorf("error"))

		w, c := CreateTestContext("GET", "/tasks/3", "")
		TaskHander{}.Show(c)

		assert.Equal(t, 404, w.Code)

		t.Cleanup(func() {
			mockDB.Close()
		})
	})
}

func TestTaskHandlerCreate(t *testing.T) {
	t.Parallel()
	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		mockDB, mock := InitMockDB(t)

		rows := mock.NewRows([]string{"id", "done"}).AddRow(0, false)
		query := regexp.QuoteMeta(`INSERT INTO "tasks" ("title","created_at","updated_at") VALUES ($1,$2,$3) RETURNING "id","done"`)
		mock.ExpectQuery(query).WillReturnRows(rows)

		d := &createParams{Title: "dummy insert task"}
		jsonString, err := json.Marshal(d)
		assert.NoError(t, err)

		w, c := CreateTestContext("POST", "/tasks", string(jsonString))
		TaskHander{}.Create(c)

		assert.Equal(t, 201, w.Code)

		var task models.Task
		body, _ := io.ReadAll(w.Body)
		err = json.Unmarshal(body, &task)
		assert.NoError(t, err)
		opt := cmpopts.IgnoreFields(models.Task{}, "CreatedAt", "UpdatedAt")
		expectBody := models.Task{ID: 0, Title: "dummy insert task", Done: false}
		assert.Empty(t, cmp.Diff(expectBody, task, opt))

		t.Cleanup(func() {
			mockDB.Close()
		})
	})

	t.Run("異常系_INSERTに失敗した場合", func(t *testing.T) {
		t.Parallel()
		mockDB, mock := InitMockDB(t)
		query := regexp.QuoteMeta(`INSERT INTO "tasks" ("title","created_at","updated_at") VALUES ($1,$2,$3) RETURNING "id","done"`)
		mock.ExpectQuery(query).WillReturnError(fmt.Errorf("error"))

		d := &createParams{Title: "dummy insert task"}
		jsonString, err := json.Marshal(d)
		assert.NoError(t, err)

		w, c := CreateTestContext("POST", "/tasks", string(jsonString))
		TaskHander{}.Create(c)

		assert.Equal(t, 500, w.Code)

		t.Cleanup(func() {
			mockDB.Close()
		})
	})

	t.Run("異常系_リクエストパラメーターが間違えている場合", func(t *testing.T) {
		t.Parallel()
		mockDB, _ := InitMockDB(t)

		d := struct{ InvalidTitle string }{InvalidTitle: "invalid title"}
		jsonString, err := json.Marshal(d)
		assert.NoError(t, err)

		w, c := CreateTestContext("POST", "/tasks", string(jsonString))
		TaskHander{}.Create(c)

		assert.Equal(t, 400, w.Code)

		t.Cleanup(func() {
			mockDB.Close()
		})
	})
}

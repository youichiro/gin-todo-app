package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/youichiro/go-todo-app/internal/client"
	"github.com/youichiro/go-todo-app/internal/models"
)

func TestMain(m *testing.M) {
	setup()
	defer teardown()
	m.Run()
}

func setup() {
	fmt.Println("setup")
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func teardown() {
	fmt.Println("teardown")
}

func TestHelloRoute(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	assert.NoError(t, err)

	r := SetupRouter(mockDB)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message": "hello world!"}`, w.Body.String())
}

func TestTasksCRUDRoute(t *testing.T) {
	var task models.Task
	var tasks []models.Task
	db := client.InitDB("test")
	s := httptest.NewServer(SetupRouter(db))
	httpClient := &http.Client{}

	// post /tasks
	res, err := http.Post(s.URL+"/tasks", "application/json", bytes.NewBuffer([]byte(`{"title": "new task"}`)))
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 201, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &task)
	assert.NoError(t, err)
	assert.Equal(t, "new task", task.Title)

	// get /tasks
	res, err = http.Get(s.URL + "/tasks")
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)
	body, _ = io.ReadAll(res.Body)
	err = json.Unmarshal(body, &tasks)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, "new task", tasks[0].Title)

	// put /tasks/:id
	req, err := http.NewRequest("PUT", s.URL+"/tasks/"+strconv.Itoa(tasks[0].ID), bytes.NewBuffer([]byte(`{"title": "update task"}`)))
	assert.NoError(t, err)
	res, err = httpClient.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)

	// get /tasks/:id
	res, err = http.Get(s.URL + "/tasks/" + strconv.Itoa(tasks[0].ID))
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)
	body, _ = io.ReadAll(res.Body)
	err = json.Unmarshal(body, &task)
	assert.NoError(t, err)
	assert.Equal(t, "update task", task.Title)

	// delete /tasks/:id
	req, err = http.NewRequest("DELETE", s.URL+"/tasks/"+strconv.Itoa(task.ID), nil)
	assert.NoError(t, err)
	res, err = httpClient.Do(req)
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 204, res.StatusCode)

	t.Cleanup(func() {
		db.Close()
	})
}

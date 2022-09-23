package router

import (
	"encoding/json"
	"example/web-service-gin/internal/client"
	"example/web-service-gin/internal/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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

func TestRootRoute(t *testing.T) {
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message": "hello world!"}`, w.Body.String())
}

func TestTasksRoute(t *testing.T) {
	var tasks []models.Task
	client.Connect("development")
	defer client.DB.Close()

	s := httptest.NewServer(SetupRouter())
	res, err := http.Get(s.URL + "/tasks")
	assert.NoError(t, err)
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &tasks)
	assert.NoError(t, err)
	assert.Equal(t, "sample task1", tasks[0].Title)
}

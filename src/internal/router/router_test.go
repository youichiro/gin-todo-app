package router

import (
	"encoding/json"
	"example/web-service-gin/internal/client"
	"example/web-service-gin/internal/models"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Task models.Task

func TestRootRoute(t *testing.T) {
	r := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"message": "hello world!"}`, w.Body.String())
}

func TestTasksRoute(t *testing.T) {
	var tasks []Task
	client.Connect("development")
	defer client.DB.Close()

	s := httptest.NewServer(SetupRouter())
	res, err := http.Get(s.URL + "/tasks")
	if err != nil {
		t.Errorf("http get err should be nil: %v", err)
	}
	defer res.Body.Close()

	assert.Equal(t, 200, res.StatusCode)
	body, _ := io.ReadAll(res.Body)
	if err := json.Unmarshal(body, &tasks); err != nil {
		fmt.Println(err)
		return
	}
	assert.Equal(t, "sample task1", tasks[0].Title)
}

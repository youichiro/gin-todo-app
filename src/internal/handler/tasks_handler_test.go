package handler

import (
	"example/web-service-gin/internal/client"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/sqlboiler/boil"
)

func TestTaskHandlerIndex(t *testing.T) {
	db := client.PostgresClientProvider{}
	db.Connect("test")
	defer db.Close()

	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer mockDB.Close()

	// Inject mock instance into boil.
	boil.SetDB(mockDB)
	db.Set(mockDB)

	rows := mock.NewRows([]string{"id", "title", "done"}).AddRow(1, "hoge", false)
	query := regexp.QuoteMeta(`SELECT "tasks".* FROM "tasks";`)
	mock.ExpectQuery(query).WillReturnRows(rows)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/tasks", nil)

	TaskHander{}.Index(c)

	assert.NoError(t, mock.ExpectationsWereMet())
	assert.Equal(t, 200, w.Code)
	expected := `[{"id": 1, "title": "hoge", "done": false, "created_at": "0001-01-01T00:00:00Z", "updated_at": "0001-01-01T00:00:00Z"}]`
	assert.JSONEq(t, expected, w.Body.String())
}

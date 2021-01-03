package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hiroyaonoe/todoapp-server/database"
)

type request struct {
	name     string // test名
	path     string
	cookie   map[string]string
	method   string // GET, POST, PUT, DELETE
	body     string // request body
	wantCode int
	wantData string
}

func TestUser_Test1(t *testing.T) {

	router := setUp(t)

	tests := []request{
		{
			name:   "userの作成",
			path:   "/user",
			method: "POST",
			body: `{
				"name":"username",
				"password":"password",
				"email":"example@example.com"
			}`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			url := "http://localhost" + router.Port + "/api/v1" + tt.path
			req := httptest.NewRequest(tt.method, url, bytes.NewBufferString(tt.body))
			router.Gin.ServeHTTP(w, req)

		})
	}
}

func setUp(t *testing.T) (r *Routing) {
	t.Helper()

	db := database.NewTestDB()
	db.Migrate()
	db.Connect().LogMode(true)

	// databaseを初期化する
	db.Connect().Exec("TRUNCATE TABLE *")

	user := new(database.UserRepository)
	task := new(database.TaskRepository)
	r = NewRouting(db, user, task)
}

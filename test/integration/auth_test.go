//go:build integration
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/server"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/worker"
	"github.com/ramdhanrizkij/bytecode-api/pkg/logger"
	"github.com/ramdhanrizkij/bytecode-api/test/integration/testhelper"
)

func TestAuth_Register(t *testing.T) {
	tdb := testhelper.SetupTestDB(t)
	defer tdb.Teardown(t)

	cfg := &config.Config{
		App: config.AppConfig{Name: "test-api"},
		JWT: config.JWTConfig{Secret: "test-secret", ExpiryHours: 1},
		Log: config.LogConfig{Level: "error"},
	}
	
	_ = logger.InitGlobal("error")
	wp := worker.NewWorkerPool(1, 10, logger.Log)
	sched := worker.NewScheduler(logger.Log)
	
	srv := server.NewServer(cfg, tdb.DB, logger.Log, wp, sched)
	srv.SetupRoutes()
	app := srv.AppForTest() // I need to add this helper to Server struct

	t.Run("Success Register", func(t *testing.T) {
		tdb.TruncateAll(t)
		// Seed "user" role which is required by Register
		tdb.DB.Exec("INSERT INTO roles (id, name) VALUES (gen_random_uuid(), 'user')")

		payload := map[string]string{
			"name":     "Test User",
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(payload)
		
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, _ := app.Test(req)
		
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		
		data := result["data"].(map[string]interface{})
		user := data["user"].(map[string]interface{})
		assert.Equal(t, "test@example.com", user["email"])
		assert.NotEmpty(t, data["token"])
	})

	t.Run("Duplicate Email", func(t *testing.T) {
		payload := map[string]string{
			"name":     "Test User 2",
			"email":    "test@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(payload)
		
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})
}

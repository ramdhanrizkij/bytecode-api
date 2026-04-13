//go:build integration
package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/config"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/server"
	"github.com/ramdhanrizkij/bytecode-api/internal/core/worker"
	"github.com/ramdhanrizkij/bytecode-api/pkg/logger"
	"github.com/ramdhanrizkij/bytecode-api/pkg/jwt"
	"github.com/ramdhanrizkij/bytecode-api/test/integration/testhelper"
)

func TestRole_CRUD(t *testing.T) {
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
	app := srv.AppForTest()

	// Generate a superadmin token for authorization
	token, _ := jwt.GenerateToken("admin-id", "admin@example.com", "superadmin", cfg.JWT.Secret, 1)

	t.Run("Create Role", func(t *testing.T) {
		tdb.TruncateAll(t)
		
		payload := map[string]string{
			"name":        "Manager",
			"description": "Manage things",
		}
		body, _ := json.Marshal(payload)
		
		req := httptest.NewRequest(http.MethodPost, "/api/v1/roles", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("Get All Roles", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/roles", nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

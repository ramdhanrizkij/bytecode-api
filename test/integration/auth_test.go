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
		JWT: config.JWTConfig{Secret: "test-secret", ExpiryHours: 1, RefreshExpiryHours: 24},
		Log: config.LogConfig{Level: "error"},
	}

	_ = logger.InitGlobal("error")
	wp := worker.NewWorkerPool(1, 10, logger.Log)
	sched := worker.NewScheduler(logger.Log)

	srv, err := server.NewServer(cfg, tdb.DB, logger.Log, wp, sched)
	assert.NoError(t, err)
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
		assert.NotEmpty(t, data["refresh_token"])
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

func TestAuth_RefreshAndLogout(t *testing.T) {
	tdb := testhelper.SetupTestDB(t)
	defer tdb.Teardown(t)

	cfg := &config.Config{
		App: config.AppConfig{Name: "test-api"},
		JWT: config.JWTConfig{Secret: "test-secret", ExpiryHours: 1, RefreshExpiryHours: 24},
		Log: config.LogConfig{Level: "error"},
	}

	_ = logger.InitGlobal("error")
	wp := worker.NewWorkerPool(1, 10, logger.Log)
	sched := worker.NewScheduler(logger.Log)

	srv, err := server.NewServer(cfg, tdb.DB, logger.Log, wp, sched)
	assert.NoError(t, err)
	srv.SetupRoutes()
	app := srv.AppForTest()

	tdb.TruncateAll(t)
	tdb.DB.Exec("INSERT INTO roles (id, name) VALUES (gen_random_uuid(), 'user')")

	registerPayload := map[string]string{
		"name":     "Refresh User",
		"email":    "refresh@example.com",
		"password": "password123",
	}
	registerBody, _ := json.Marshal(registerPayload)
	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")

	registerResp, _ := app.Test(registerReq)
	assert.Equal(t, http.StatusCreated, registerResp.StatusCode)

	var registerResult map[string]interface{}
	json.NewDecoder(registerResp.Body).Decode(&registerResult)
	refreshToken := registerResult["data"].(map[string]interface{})["refresh_token"].(string)
	assert.NotEmpty(t, refreshToken)

	refreshPayload := map[string]string{"refresh_token": refreshToken}
	refreshBody, _ := json.Marshal(refreshPayload)
	refreshReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(refreshBody))
	refreshReq.Header.Set("Content-Type", "application/json")

	refreshResp, _ := app.Test(refreshReq)
	assert.Equal(t, http.StatusOK, refreshResp.StatusCode)

	var refreshResult map[string]interface{}
	json.NewDecoder(refreshResp.Body).Decode(&refreshResult)
	refreshData := refreshResult["data"].(map[string]interface{})
	rotatedRefreshToken := refreshData["refresh_token"].(string)
	assert.NotEmpty(t, refreshData["token"])
	assert.NotEmpty(t, rotatedRefreshToken)
	assert.NotEqual(t, refreshToken, rotatedRefreshToken)

	reuseBody, _ := json.Marshal(refreshPayload)
	reuseReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(reuseBody))
	reuseReq.Header.Set("Content-Type", "application/json")

	reuseResp, _ := app.Test(reuseReq)
	assert.Equal(t, http.StatusUnauthorized, reuseResp.StatusCode)

	logoutPayload := map[string]string{"refresh_token": rotatedRefreshToken}
	logoutBody, _ := json.Marshal(logoutPayload)
	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", bytes.NewReader(logoutBody))
	logoutReq.Header.Set("Content-Type", "application/json")

	logoutResp, _ := app.Test(logoutReq)
	assert.Equal(t, http.StatusOK, logoutResp.StatusCode)

	afterLogoutBody, _ := json.Marshal(logoutPayload)
	afterLogoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewReader(afterLogoutBody))
	afterLogoutReq.Header.Set("Content-Type", "application/json")

	afterLogoutResp, _ := app.Test(afterLogoutReq)
	assert.Equal(t, http.StatusUnauthorized, afterLogoutResp.StatusCode)
}

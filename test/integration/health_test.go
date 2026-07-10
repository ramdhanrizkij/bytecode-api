//go:build integration

package integration

import (
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

func TestHealthCheck(t *testing.T) {
	tdb := testhelper.SetupTestDB(t)
	defer tdb.Teardown(t)

	cfg := &config.Config{
		App: config.AppConfig{Name: "test-api", Env: "test"},
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

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)

	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusOK), result["meta"].(map[string]interface{})["code"])
	assert.Equal(t, "service is healthy", result["meta"].(map[string]interface{})["message"])

	data := result["data"].(map[string]interface{})
	assert.Equal(t, "ok", data["status"])
	assert.Equal(t, "test-api", data["service"])
	assert.Equal(t, "test", data["environment"])
	assert.Equal(t, "up", data["database"])
	assert.Equal(t, "disabled", data["cache"])
	assert.Equal(t, "local", data["storage"])
}

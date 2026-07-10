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

func TestUser_ProfileLifecycle(t *testing.T) {
	tdb := testhelper.SetupTestDB(t)
	defer tdb.Teardown(t)

	cfg := &config.Config{
		App: config.AppConfig{Name: "test-api", Env: "test"},
		JWT: config.JWTConfig{Secret: "test-secret", ExpiryHours: 1, RefreshExpiryHours: 24},
		Log: config.LogConfig{Level: "error"},
		Storage: config.StorageConfig{
			Provider:      "local",
			DefaultBucket: "uploads",
			LocalPath:     t.TempDir(),
			BaseURL:       "/storage",
		},
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

	registerPayload := map[string]interface{}{
		"name":     "Profile User",
		"email":    "profile@example.com",
		"password": "password123",
	}
	registerBody, _ := json.Marshal(registerPayload)
	registerReq := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")

	registerResp, _ := app.Test(registerReq)
	assert.Equal(t, http.StatusCreated, registerResp.StatusCode)

	var registerResult map[string]interface{}
	json.NewDecoder(registerResp.Body).Decode(&registerResult)
	registerData := registerResult["data"].(map[string]interface{})
	token := registerData["token"].(string)

	getMeReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	getMeReq.Header.Set("Authorization", "Bearer "+token)

	getMeResp, _ := app.Test(getMeReq)
	assert.Equal(t, http.StatusOK, getMeResp.StatusCode)

	var getMeResult map[string]interface{}
	json.NewDecoder(getMeResp.Body).Decode(&getMeResult)
	getMeData := getMeResult["data"].(map[string]interface{})
	_, hasProfilePicture := getMeData["profile_picture"]
	assert.False(t, hasProfilePicture)

	updatePayload := map[string]interface{}{
		"name":            "Profile User Updated",
		"email":           "profile.updated@example.com",
		"profile_picture": "avatars/profile-user.png",
	}
	updateBody, _ := json.Marshal(updatePayload)
	updateReq := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/profile", bytes.NewReader(updateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("Authorization", "Bearer "+token)

	updateResp, _ := app.Test(updateReq)
	assert.Equal(t, http.StatusOK, updateResp.StatusCode)

	var updateResult map[string]interface{}
	json.NewDecoder(updateResp.Body).Decode(&updateResult)
	updateData := updateResult["data"].(map[string]interface{})
	profilePicture := updateData["profile_picture"].(map[string]interface{})
	assert.Equal(t, "uploads", profilePicture["bucket"])
	assert.Equal(t, "avatars/profile-user.png", profilePicture["key"])
	assert.Equal(t, "/storage/uploads/avatars/profile-user.png", profilePicture["url"])

	getMeAfterReq := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	getMeAfterReq.Header.Set("Authorization", "Bearer "+token)

	getMeAfterResp, _ := app.Test(getMeAfterReq)
	assert.Equal(t, http.StatusOK, getMeAfterResp.StatusCode)

	var getMeAfterResult map[string]interface{}
	json.NewDecoder(getMeAfterResp.Body).Decode(&getMeAfterResult)
	getMeAfterData := getMeAfterResult["data"].(map[string]interface{})
	assert.Equal(t, "profile.updated@example.com", getMeAfterData["email"])
	assert.Equal(t, "Profile User Updated", getMeAfterData["name"])
	assert.Equal(t, "avatars/profile-user.png", getMeAfterData["profile_picture"].(map[string]interface{})["key"])
}

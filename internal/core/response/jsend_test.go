package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSendSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := map[string]string{"foo": "bar"}
	SendSuccess(c, data)

	if w.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
	}

	var resp JSendResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Status != StatusSuccess {
		t.Errorf("got status %s, want success", resp.Status)
	}
}

func TestSendFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := "invalid input"
	SendFail(c, data, http.StatusUnprocessableEntity)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf(
			"got status %d, want %d",
			w.Code,
			http.StatusUnprocessableEntity,
		)
	}

	var resp JSendResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Status != StatusFail {
		t.Errorf("got status %s, want fail", resp.Status)
	}
}

func TestSendError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	SendError(c, "boom", http.StatusInternalServerError, nil)

	if w.Code != http.StatusInternalServerError {
		t.Errorf(
			"got status %d, want %d",
			w.Code,
			http.StatusInternalServerError,
		)
	}

	var resp JSendResponse
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Status != StatusError || resp.Message != "boom" {
		t.Errorf("got %+v, want error/boom", resp)
	}
}

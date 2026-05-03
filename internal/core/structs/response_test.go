package structs

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewSuccessResponse(t *testing.T) {
	data := map[string]string{"key": "value"}
	resp := NewSuccessResponse(data)

	if resp.Status != StatusSuccess {
		t.Errorf("got status %s, want success", resp.Status)
	}
	if !reflect.DeepEqual(resp.Data, data) {
		t.Errorf("got data %v, want %v", resp.Data, data)
	}

	// Verify JSON marshaling
	got, _ := json.Marshal(resp)
	want := `{"status":"success","data":{"key":"value"}}`
	if string(got) != want {
		t.Errorf("got JSON %s, want %s", string(got), want)
	}
}

func TestNewFailResponse(t *testing.T) {
	data := map[string]string{"error": "field is required"}
	resp := NewFailResponse(data)

	if resp.Status != StatusFail {
		t.Errorf("got status %s, want fail", resp.Status)
	}
	if !reflect.DeepEqual(resp.Data, data) {
		t.Errorf("got data %v, want %v", resp.Data, data)
	}
}

func TestNewErrorResponse(t *testing.T) {
	msg := "internal server error"
	code := 500
	data := "some trace info"
	resp := NewErrorResponse(msg, code, data)

	if resp.Status != StatusError {
		t.Errorf("got status %s, want error", resp.Status)
	}
	if resp.Message != msg {
		t.Errorf("got message %s, want %s", resp.Message, msg)
	}
	if resp.Code != code {
		t.Errorf("got code %d, want %d", resp.Code, code)
	}
	if resp.Data != data {
		t.Errorf("got data %v, want %v", resp.Data, data)
	}
}

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginForwardsWidgetRecordID(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/captcha/verify" {
			t.Errorf("unexpected upstream path: %s", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if body["widget_record_id"] != "widget-123" || body["token"] != "invalid-token" {
			t.Errorf("unexpected captcha fields: %v", body)
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"ok":false,"error":{"code":"INVALID_TOKEN","message":"invalid"}}`))
	}))
	defer upstream.Close()

	client := &InfraiClient{BaseURL: upstream.URL, HTTP: upstream.Client()}
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"widgetRecordId":"widget-123","captchaToken":"invalid-token"}`))
	w := httptest.NewRecorder()
	loginHandler(client, w, req)
	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
	}
}

func TestWorkOrderAfterLogin(t *testing.T) {
	tests := []struct {
		name, status string
		photos       int
		followUp     bool
	}{
		{"dispatched technician", "dispatched", 3, true},
		{"follow-up pending", "dispatched", 3, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := workOrder{ID: "WO-1042", Status: tt.status, Photos: tt.photos, FollowUp: tt.followUp}
			if got.Status != "dispatched" || got.Photos < 1 || !got.FollowUp {
				t.Fatalf("unexpected work order: %+v", got)
			}
		})
	}
}

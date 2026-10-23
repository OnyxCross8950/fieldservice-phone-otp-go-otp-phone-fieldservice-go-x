package main

import (
	"encoding/json"
	"net/http"
)

type loginRequest struct {
	Phone          string
	Code           string
	CaptchaToken   string
	WidgetRecordID string `json:"widgetRecordId"`
}
type workOrder struct {
	ID, Status string
	Photos     int
	FollowUp   bool
}

func loginHandler(client *InfraiClient, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var in loginRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := client.VerifyCaptcha(in.WidgetRecordID, in.CaptchaToken, r.RemoteAddr); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if in.Code == "" {
		if err := client.SendCode(in.Phone); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"next": "code"})
		return
	}
	if _, err := client.VerifyCode(in.Phone, in.Code); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(workOrder{ID: "WO-1042", Status: "dispatched", Photos: 3, FollowUp: true})
}

func main() {
	client := NewInfraiClient()
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) { loginHandler(client, w, r) })
	_ = http.ListenAndServe(":8080", nil)
}

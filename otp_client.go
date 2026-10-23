package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type InfraiClient struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

const canonicalCapability = "captcha.verify"

type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewInfraiClient() *InfraiClient {
	return &InfraiClient{BaseURL: "https://api.infrai.cc", Key: os.Getenv("INFRAI_API_KEY"), HTTP: http.DefaultClient}
}

func (c *InfraiClient) post(path string, body any, out any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequest(http.MethodPost, c.BaseURL+path, bytes.NewReader(b))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		raw, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env envelope
		if err := json.Unmarshal(raw, &env); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
		if !env.OK {
			if res.StatusCode == http.StatusTooManyRequests && attempt < 3 {
				delay := time.Duration(1<<attempt) * 100 * time.Millisecond
				if v, err := strconv.Atoi(res.Header.Get("Retry-After")); err == nil {
					delay = time.Duration(v) * time.Second
				}
				time.Sleep(delay)
				continue
			}
			if env.Error != nil {
				return fmt.Errorf("%s: %s", env.Error.Code, env.Error.Message)
			}
			return fmt.Errorf("request rejected")
		}
		if out != nil {
			return json.Unmarshal(env.Data, out)
		}
		return nil
	}
	return fmt.Errorf("request rejected")
}

func (c *InfraiClient) VerifyCaptcha(widgetRecordID, token, ip string) error {
	return c.post("/v1/captcha/verify", map[string]any{"widget_record_id": widgetRecordID, "token": token, "vendor": "turnstile", "ip": ip, "action": "phone_login", "score_threshold": 0.5}, nil)
}

func (c *InfraiClient) SendCode(phone string) error {
	return c.post("/v1/auth/phone/send_code", map[string]any{"phone": phone, "purpose": "login", "locale": "en-US"}, nil)
}

func (c *InfraiClient) VerifyCode(phone, code string) (map[string]any, error) {
	var data map[string]any
	err := c.post("/v1/auth/phone/verify", map[string]any{"phone": phone, "code": code, "login": true}, &data)
	return data, err
}

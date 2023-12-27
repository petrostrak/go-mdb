package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

func Test_readIDParams(t *testing.T) {
	var app application

	tests := []struct {
		ID      string
		expeted uuid.UUID
	}{
		{"121f03cd-ce8c-447d-8747-fb8cb7aa3a52", uuid.MustParse("121f03cd-ce8c-447d-8747-fb8cb7aa3a52")},
		{"2f0141f8-f325-4b15-9973-e7b34852e298", uuid.MustParse("2f0141f8-f325-4b15-9973-e7b34852e298")},
		{"3659fbd7-7ba2-4151-95a2-b977ebf79307", uuid.MustParse("3659fbd7-7ba2-4151-95a2-b977ebf79307")},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := req.Context()
		ctx = context.WithValue(ctx, httprouter.ParamsKey, httprouter.Params{
			{Key: "id", Value: tt.ID},
		})
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		handler.ServeHTTP(rr, req)

		result, _ := app.readIDParams(req)

		if result != tt.expeted {
			t.Errorf("Expected %v but got %v\n", tt.expeted, result)
		}
	}
}

func Test_writeJSON(t *testing.T) {
	var app application
	rr := httptest.NewRecorder()
	payload := make(map[string]any)
	payload["foo"] = false

	headers := make(http.Header)
	headers.Add("FOO", "BAR")
	err := app.writeJSON(rr, http.StatusOK, payload, headers)
	if err != nil {
		t.Errorf("failed to write JSON: %v\n", err)
	}
}

func Test_readJSON(t *testing.T) {
	var app application

	sampleJSON := map[string]any{
		"foo": "bar",
	}
	body, _ := json.Marshal(sampleJSON)

	var decodedJSON struct {
		Foo string `json:"foo"`
	}

	req, err := http.NewRequest("POST", "/", bytes.NewReader(body))
	if err != nil {
		t.Log("Error", err.Error())
	}

	rr := httptest.NewRecorder()
	defer req.Body.Close()

	err = app.readJSON(rr, req, &decodedJSON)
	if err != nil {
		t.Error("failed to decode json", err)
	}

	badJSON := `
		{
			"foo": "bar"
		}
		{
			"alpha": "beta"
		}`

	req, err = http.NewRequest("POST", "/", bytes.NewReader([]byte(badJSON)))
	if err != nil {
		t.Log("Error", err)
	}

	err = app.readJSON(rr, req, &decodedJSON)
	if err == nil {
		t.Error("did not get an error with bad json")
	}
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/petrostrak/gomdb/internal/validator"
)

var app application

func Test_readIDParams(t *testing.T) {

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

func Test_readCSV(t *testing.T) {
	tests := []struct {
		qs           url.Values
		key          string
		defaultValue []string
		expected     []string
	}{
		{url.Values{"genres": []string{"crime", "comedy"}}, "genres", []string{}, []string{"crime", "comedy"}},
		{url.Values{"genres": []string{"crime", "comedy", "adventure"}}, "genres", []string{}, []string{"crime", "comedy", "adventure"}},
		{url.Values{"": []string{}}, "genres", []string{}, []string{}},
	}

	for _, tt := range tests {
		csv := app.readCSV(tt.qs, tt.key, tt.defaultValue)

		for i, s := range csv {
			if s != tt.expected[i] {
				t.Errorf("expected %s but got %s\n", tt.expected[i], s)
			}
		}
	}
}

func Test_readString(t *testing.T) {
	tests := []struct {
		qs           url.Values
		key          string
		defaultValue string
		expected     string
	}{
		{url.Values{"genres": []string{"comedy"}}, "genres", "", "comedy"},
		{url.Values{"genres": []string{"crime", "comedy", "adventure"}}, "genres", "", "crime"},
		{url.Values{"": []string{}}, "genres", "", ""},
	}

	for _, tt := range tests {
		s := app.readString(tt.qs, tt.key, tt.defaultValue)

		if s != tt.expected {
			t.Errorf("expected %s but got %s\n", tt.expected, s)
		}
	}
}

func Test_readInt(t *testing.T) {
	tests := []struct {
		qs           url.Values
		key          string
		defaultValue int
		expected     int
	}{
		{url.Values{"page": []string{"4"}}, "page", 1, 4},
		{url.Values{"page_sige": []string{"abc"}}, "page_size", 20, 20},
	}

	var v validator.Validator

	for _, tt := range tests {
		i := app.readInt(tt.qs, tt.key, tt.defaultValue, &v)

		if i != tt.expected {
			t.Errorf("expected %d but got %d\n", tt.expected, i)
		}
	}
}

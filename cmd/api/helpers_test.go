package main

import (
	"context"
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
			{"id", tt.ID},
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

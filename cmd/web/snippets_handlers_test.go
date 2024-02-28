package main

import (
	"net/http"
	"testing"

	"github.com/svidlak/lets-go/internal/assert"
	"github.com/svidlak/lets-go/internal/models/mocks"
)

func TestHome(t *testing.T) {
	mockData := mocks.MockSnippet
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Home page",
			urlPath:  "/",
			wantCode: http.StatusOK,
			wantBody: mockData.Title,
		},
		{
			name:     "Not found page",
			urlPath:  "/non-existing",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			statusCode, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, statusCode, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestSnippetView(t *testing.T) {
	mockData := mocks.MockSnippet
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: mockData.Content,
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			statusCode, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, statusCode, tt.wantCode)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

package main

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andres085/task_manager/internal/assert"
	"github.com/andres085/task_manager/internal/models/mocks"
)

func TestCommonHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	commonHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com; img-src 'self' blob: data:;"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	expectedValue = "Go"
	assert.Equal(t, rs.Header.Get("Server"), expectedValue)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}

func TestLogRequest(t *testing.T) {
	spyLogger := &mocks.SpyLogger{}
	logger := slog.New(spyLogger)

	app := &application{
		logger: logger,
	}

	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	app.logRequest(next).ServeHTTP(rr, r)

	expectedValue := "received request"
	logEntry := spyLogger.Entries[0]

	assert.Equal(t, spyLogger.Called, true)
	assert.Equal(t, len(spyLogger.Entries), 1)
	assert.Equal(t, logEntry.Message, expectedValue)
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("test panic")
}

func TestRecoverPanic(t *testing.T) {
	app := newTestApplication(t)

	handler := app.recoverPanic(http.HandlerFunc(panicHandler))
	req, err := http.NewRequest(http.MethodGet, "/panic", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	expectedValue := "close"
	assert.Equal(t, rr.Header().Get("Connection"), expectedValue)
	assert.Equal(t, rr.Code, http.StatusInternalServerError)
}

func TestRequireAuthentication(t *testing.T) {
	t.Run("Unauthenthicated", func(t *testing.T) {
		app := newTestApplication(t)
		rr := httptest.NewRecorder()

		r, err := http.NewRequest(http.MethodGet, "/workspace/view", nil)
		if err != nil {
			t.Fatal(err)
		}

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		app.requireAuthentication(next).ServeHTTP(rr, r)

		expectedValue := "/user/login"

		assert.Equal(t, rr.Result().StatusCode, http.StatusSeeOther)
		assert.Equal(t, rr.Header().Get("Location"), expectedValue)
	})

	t.Run("Authenthicated", func(t *testing.T) {
		app := newTestApplication(t)

		authenticatedRequest, err := http.NewRequest(http.MethodGet, "/workspace/view", nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx := context.WithValue(authenticatedRequest.Context(), isAuthenticatedContextKey, true)
		authenticatedRequest = authenticatedRequest.WithContext(ctx)

		rr := httptest.NewRecorder()

		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Authenticated"))
		})

		app.requireAuthentication(next).ServeHTTP(rr, authenticatedRequest)

		assert.Equal(t, rr.Result().StatusCode, http.StatusOK)
		assert.Equal(t, rr.Body.String(), "Authenticated")
		assert.Equal(t, rr.Header().Get("Cache-Control"), "no-store")
	})
}

func TestNosurf(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodPost, "/workspace/create", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	noSurf(next).ServeHTTP(rr, r)

	assert.Equal(t, rr.Result().StatusCode, http.StatusBadRequest)
}

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckSessionHandler(t *testing.T) {
	t.Parallel()

	api := NewMyHandler()
	api.sessions["valid-session"] = &Session{
		UserID:    1,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	api.users = map[string]*User{
		"testuser": {ID: 1, Username: "testuser", Email: "test@example.com"},
	}

	tests := []TestCookieCase{
		{
			name: "Valid session",
			cookie: &http.Cookie{
				Name:  "session_id",
				Value: "valid-session",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"id\":1,\"username\":\"testuser\",\"email\":\"test@example.com\"}",
		},
		{
			name:           "No cookie provided",
			cookie:         nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid cookie\n",
		},
		{
			name: "Invalid session",
			cookie: &http.Cookie{
				Name:  "session_id",
				Value: "invalid-session",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Unauthorized\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/checkSession", nil)
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}

			recorder := httptest.NewRecorder()
			api.checkSession(recorder, req)

			resp := recorder.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			body := recorder.Body.String()
			if body != tc.expectedBody {
				t.Errorf("expected body %q, got %q", tc.expectedBody, body)
			}
		})
	}
}

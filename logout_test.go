package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestCookieCase struct {
	name           string
	cookie         *http.Cookie
	expectedStatus int
	expectedBody   string
}

func TestLogoutHandler(t *testing.T) {
	t.Parallel()

	api := NewMyHandler()
	api.sessions["valid-session"] = &Session{
		UserID:    1,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	testCases := []TestCookieCase{
		{
			name: "Successful logout",
			cookie: &http.Cookie{
				Name:  "session_id",
				Value: "valid-session",
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "\"Successfuly logged out\"",
		},
		{
			name:           "No cookie provided",
			cookie:         nil,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "No cookie\n",
		},
		{
			name: "Invalid session",
			cookie: &http.Cookie{
				Name:  "session_id",
				Value: "invalid-session",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "No cookie\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/logout", nil)
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}

			recorder := httptest.NewRecorder()
			api.logoutHandler(recorder, req)

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

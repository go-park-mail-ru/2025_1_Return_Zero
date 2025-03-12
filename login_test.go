package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	name           string
	requestBody    string
	expectedStatus int
}

func TestLogin(t *testing.T) {
	t.Parallel()

	api := NewMyHandler()
	password, err := HashPassword("vasya")
	require.NoError(t, err, "failed to hash password")
	api.users["Vasily"] = &User{
		ID:       0,
		Username: "Vasily",
		Password: password,
		Email:    "supervasya@gmail.com",
	}

	testCases := []TestCase{
		{
			name:           "Login with username",
			requestBody:    `{"username": "Vasily", "password": "vasya"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Login with email",
			requestBody:    `{"password": "vasya", "email": "supervasya@gmail.com"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Login with non-existing email",
			requestBody:    `{"password": "vasya", "email": "supervasya322@gmail.com"}`,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Login with wrong password",
			requestBody:    `{"password": "vasya322", "email": "supervasya@gmail.com"}`,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			body := bytes.NewReader([]byte(tc.requestBody))
			r := httptest.NewRequest("POST", "/login", body)
			w := httptest.NewRecorder()

			api.loginHandler(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

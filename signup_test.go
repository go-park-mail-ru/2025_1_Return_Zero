package main

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignup(t *testing.T) {
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
			name:           "Signup",
			requestBody:    `{"username": "Vladimir", "password": "vova", "email": "vova@mail.ru"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Username existed",
			requestBody:    `{"username": "Vasily", "password": "vasya11", "email": "vasily@yandex.ru"}`,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Email existed",
			requestBody:    `{"username": "Vladimir22", "password": "vova2", "email": "supervasya@gmail.com"}`,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Invalid username",
			requestBody:    `{"username": "b", "password": "Beb1234567", "email": "bebe123456@gmail.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid password",
			requestBody:    `{"username": "bebe123", "password": "B", "email": "bebe123456@gmail.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid email",
			requestBody:    `{"username": "bebe123", "password": "Beb1234567", "email": "bebe123456gmail.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			body := bytes.NewReader([]byte(tc.requestBody))
			r := httptest.NewRequest("POST", "/signup", body)
			w := httptest.NewRecorder()

			api.signupHandler(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}
		})
	}
}

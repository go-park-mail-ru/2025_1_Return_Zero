package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFilters(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		filters         Filters
		expectedFilters Filters
		expectedError   error
	}{
		{
			name:    "empty filters",
			filters: Filters{},
			expectedFilters: Filters{
				Offset: 0,
				Limit:  0,
			},
			expectedError: nil,
		},
		{
			name: "invalid offset",
			filters: Filters{
				Offset: 10001,
			},
			expectedFilters: Filters{
				Offset: 10000,
				Limit:  0,
			},
			expectedError: nil,
		},
		{
			name: "invalid limit",
			filters: Filters{
				Limit: 101,
			},
			expectedFilters: Filters{
				Offset: 0,
				Limit:  100,
			},
			expectedError: nil,
		},
		{
			name: "invalid offset and limit",
			filters: Filters{
				Offset: 10001,
				Limit:  101,
			},
			expectedFilters: Filters{
				Offset: 10000,
				Limit:  100,
			},
			expectedError: nil,
		},
		{
			name: "offset below zero",
			filters: Filters{
				Offset: -1,
			},
			expectedFilters: Filters{
				Offset: 0,
				Limit:  0,
			},
			expectedError: errors.New("invalid offset: should be greater than 0"),
		},
		{
			name: "limit below zero",
			filters: Filters{
				Limit: -1,
			},
			expectedFilters: Filters{
				Offset: 0,
				Limit:  0,
			},
			expectedError: errors.New("invalid limit: should be greater than 0"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testFilters := testCase.filters
			err := testFilters.Validate()
			assert.Equal(t, testCase.expectedError, err)
			assert.Equal(t, testCase.expectedFilters, testFilters)
		})
	}
}

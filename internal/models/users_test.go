package models

import (
	"testing"

	"github.com/hqr999/Snippetbox/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	// Skip the test if the "-short" flag is provided when running the test.
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}
	// Set up a suite of table-driven tests and expected results.
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			// Call the newTestDB() helper function to get a connection pool to
			// our test database. Calling this here -- inside t.Run() -- means
			// that fresh database tables and data will be set up and torn down
			// for each sub-test.
			db := newTestDB(t)

			// Create a new instance of the UserModel
			m := UserModel{db}

			// Call the UserModel.Exists() method and check that the return
			// value matches the expected values for the sub-test and there is
			// no error.
			exists, err := m.Exists(v.userID)
			assert.Equal(t, exists, v.want)
			assert.Nil(t, err)
		})
	}

}

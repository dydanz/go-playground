package postgres

import (
	"database/sql"
	"errors"
)

// MockTx for unit testing
type MockTx struct {
	ShouldFail bool
}

func (m *MockTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.ShouldFail {
		return nil, errors.New("mock exec failure")
	}
	return nil, nil
}

func (m *MockTx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.ShouldFail {
		return nil, errors.New("mock query failure")
	}
	return nil, nil
}

func (m *MockTx) QueryRow(query string, args ...interface{}) *sql.Row {
	return &sql.Row{}
}

func (m *MockTx) Commit() error {
	if m.ShouldFail {
		return errors.New("mock commit failure")
	}
	return nil
}

func (m *MockTx) Rollback() error {
	if m.ShouldFail {
		return errors.New("mock rollback failure")
	}
	return nil
}

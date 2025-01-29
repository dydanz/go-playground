package postgres

import (
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

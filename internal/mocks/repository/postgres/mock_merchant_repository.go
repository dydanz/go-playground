package postgres

import (
	"github.com/stretchr/testify/mock"

	"go-playground/internal/domain"
)

type MockMerchantRepository struct {
	mock.Mock
}

func (m *MockMerchantRepository) Create(merchant *domain.Merchant) error {
	args := m.Called(merchant)
	return args.Error(0)
}

func (m *MockMerchantRepository) GetByID(id string) (*domain.Merchant, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Merchant), args.Error(1)
}

func (m *MockMerchantRepository) GetAll() ([]domain.Merchant, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Merchant), args.Error(1)
}

func (m *MockMerchantRepository) Update(merchant *domain.Merchant) error {
	args := m.Called(merchant)
	return args.Error(0)
}

func (m *MockMerchantRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

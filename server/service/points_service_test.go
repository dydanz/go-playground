package service

import (
	"context"
	"testing"

	"go-playground/server/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock repositories
type mockPointsRepository struct {
	mock.Mock
}

type mockEventLogRepository struct {
	mock.Mock
}

// Implement PointsRepository interface
func (m *mockPointsRepository) GetCurrentBalance(ctx context.Context, customerID, programID uuid.UUID) (int, error) {
	args := m.Called(ctx, customerID, programID)
	return args.Int(0), args.Error(1)
}

func (m *mockPointsRepository) Create(ctx context.Context, ledger *domain.PointsLedger) (*domain.PointsLedger, error) {
	args := m.Called(ctx, ledger)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PointsLedger), args.Error(1)
}

func (m *mockPointsRepository) GetByCustomerAndProgram(ctx context.Context, customerID, programID uuid.UUID) ([]*domain.PointsLedger, error) {
	args := m.Called(ctx, customerID, programID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PointsLedger), args.Error(1)
}

func (m *mockPointsRepository) GetByTransactionID(ctx context.Context, programID uuid.UUID) (*domain.PointsLedger, error) {
	args := m.Called(ctx, programID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PointsLedger), args.Error(1)
}

func (m *mockPointsRepository) Delete(ctx context.Context, LedgerID uuid.UUID) error {
	args := m.Called(ctx, LedgerID)
	return args.Error(0)
}

// Implement EventLogRepository interface
func (m *mockEventLogRepository) Create(ctx context.Context, event *domain.EventLog) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockEventLogRepository) Update(ctx context.Context, event *domain.EventLog) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *mockEventLogRepository) GetByReferenceID(referenceID string) (*domain.EventLog, error) {
	args := m.Called(referenceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventLog), args.Error(1)
}

func (m *mockEventLogRepository) GetByID(cid string) (*domain.EventLog, error) {
	args := m.Called(cid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.EventLog), args.Error(1)
}

func (m *mockEventLogRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockEventLogRepository) GetByUserID(userID string) ([]domain.EventLog, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.EventLog), args.Error(1)
}

// PointsServiceTestSuite defines the test suite
type PointsServiceTestSuite struct {
	suite.Suite
	pointsRepo *mockPointsRepository
	eventRepo  *mockEventLogRepository
	service    *PointsService
}

// SetupTest is called before each test
func (s *PointsServiceTestSuite) SetupTest() {
	s.pointsRepo = new(mockPointsRepository)
	s.eventRepo = new(mockEventLogRepository)
	s.service = NewPointsService(s.pointsRepo, s.eventRepo)
}

// TestPointsServiceTestSuite runs the test suite
func TestPointsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PointsServiceTestSuite))
}

// Test cases for EarnPoints
func (s *PointsServiceTestSuite) TestEarnPoints_Success() {
	ctx := context.Background()
	customerID := uuid.New()
	programID := uuid.New()
	transactionID := uuid.New()

	req := &domain.PointsTransaction{
		TransactionID: transactionID.String(),
		CustomerID:    customerID.String(),
		ProgramID:     programID.String(),
		Points:        100,
		Type:          "earn",
	}

	s.pointsRepo.On("GetCurrentBalance", ctx, customerID, programID).Return(0, nil)
	s.pointsRepo.On("Create", ctx, mock.MatchedBy(func(l *domain.PointsLedger) bool {
		return l.MerchantCustomersID == customerID &&
			l.ProgramID == programID &&
			l.PointsEarned == 100 &&
			l.PointsBalance == 100
	})).Return(&domain.PointsLedger{
		LedgerID:            uuid.New(),
		MerchantCustomersID: customerID,
		ProgramID:           programID,
		PointsEarned:        100,
		PointsBalance:       100,
		TransactionID:       transactionID,
	}, nil)

	result, err := s.service.EarnPoints(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Points, result.Points)
	s.Equal("earn", result.Type)
}

func (s *PointsServiceTestSuite) TestEarnPoints_InvalidPoints() {
	ctx := context.Background()
	customerID := uuid.New()
	programID := uuid.New()
	transactionID := uuid.New()

	testCases := []struct {
		name   string
		points int
	}{
		{"Zero Points", 0},
		{"Negative Points", -100},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := &domain.PointsTransaction{
				TransactionID: transactionID.String(),
				CustomerID:    customerID.String(),
				ProgramID:     programID.String(),
				Points:        tc.points,
				Type:          "earn",
			}

			result, err := s.service.EarnPoints(ctx, req)

			s.Error(err)
			s.Nil(result)
			_, ok := err.(domain.ValidationError)
			s.True(ok)
		})
	}
}

// Test cases for RedeemPoints
func (s *PointsServiceTestSuite) TestRedeemPoints_Success() {
	ctx := context.Background()
	customerID := uuid.New()
	programID := uuid.New()
	transactionID := uuid.New()

	req := &domain.PointsTransaction{
		TransactionID: transactionID.String(),
		CustomerID:    customerID.String(),
		ProgramID:     programID.String(),
		Points:        50,
		Type:          "redeem",
	}

	s.pointsRepo.On("GetCurrentBalance", ctx, customerID, programID).Return(100, nil)
	s.pointsRepo.On("Create", ctx, mock.MatchedBy(func(l *domain.PointsLedger) bool {
		return l.MerchantCustomersID == customerID &&
			l.ProgramID == programID &&
			l.PointsRedeemed == 50 &&
			l.PointsBalance == 50
	})).Return(&domain.PointsLedger{
		LedgerID:            uuid.New(),
		MerchantCustomersID: customerID,
		ProgramID:           programID,
		PointsRedeemed:      50,
		PointsBalance:       50,
		TransactionID:       transactionID,
	}, nil)

	result, err := s.service.RedeemPoints(ctx, req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(req.Points, result.Points)
	s.Equal("redeem", result.Type)
}

func (s *PointsServiceTestSuite) TestRedeemPoints_InsufficientPoints() {
	ctx := context.Background()
	customerID := uuid.New()
	programID := uuid.New()
	transactionID := uuid.New()

	req := &domain.PointsTransaction{
		TransactionID: transactionID.String(),
		CustomerID:    customerID.String(),
		ProgramID:     programID.String(),
		Points:        100,
		Type:          "redeem",
	}

	s.pointsRepo.On("GetCurrentBalance", ctx, customerID, programID).Return(50, nil)

	result, err := s.service.RedeemPoints(ctx, req)

	s.Error(err)
	s.Nil(result)
	_, ok := err.(domain.BusinessLogicError)
	s.True(ok)
	s.Contains(err.Error(), "insufficient points balance")
}

// Test cases for GetBalance
func (s *PointsServiceTestSuite) TestGetBalance_Success() {
	ctx := context.Background()
	customerID := uuid.New()
	programID := uuid.New()
	expectedBalance := 100

	s.pointsRepo.On("GetCurrentBalance", ctx, customerID, programID).Return(expectedBalance, nil)

	result, err := s.service.GetBalance(ctx, customerID, programID)

	s.NoError(err)
	s.NotNil(result)
	s.Equal(customerID.String(), result.CustomerID)
	s.Equal(programID.String(), result.ProgramID)
	s.Equal(expectedBalance, result.Balance)
}

// Test cases for GetLedger
func (s *PointsServiceTestSuite) TestGetLedger_Success() {
	ctx := context.Background()
	customerID := uuid.New()
	programID := uuid.New()

	expectedLedgers := []*domain.PointsLedger{
		{
			LedgerID:            uuid.New(),
			MerchantCustomersID: customerID,
			ProgramID:           programID,
			PointsEarned:        100,
			PointsBalance:       100,
		},
		{
			LedgerID:            uuid.New(),
			MerchantCustomersID: customerID,
			ProgramID:           programID,
			PointsRedeemed:      50,
			PointsBalance:       50,
		},
	}

	s.pointsRepo.On("GetByCustomerAndProgram", ctx, customerID, programID).Return(expectedLedgers, nil)

	result, err := s.service.GetLedger(ctx, customerID, programID)

	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 2)
	s.Equal(expectedLedgers[0].PointsEarned, result[0].PointsEarned)
	s.Equal(expectedLedgers[1].PointsRedeemed, result[1].PointsRedeemed)
}

func (s *PointsServiceTestSuite) TestGetLedger_EmptyResult() {
	ctx := context.Background()
	customerID := uuid.New()
	programID := uuid.New()

	s.pointsRepo.On("GetByCustomerAndProgram", ctx, customerID, programID).Return([]*domain.PointsLedger{}, nil)

	result, err := s.service.GetLedger(ctx, customerID, programID)

	s.NoError(err)
	s.NotNil(result)
	s.Len(result, 0)
}

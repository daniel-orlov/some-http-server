// Code generated by MockGen. DO NOT EDIT.
// Source: quote_svc.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	types "some-http-server/internal/types"

	gomock "github.com/golang/mock/gomock"
)

// MockExternalSvcClient is a mock of ExternalSvcClient interface.
type MockExternalSvcClient struct {
	ctrl     *gomock.Controller
	recorder *MockExternalSvcClientMockRecorder
}

// MockExternalSvcClientMockRecorder is the mock recorder for MockExternalSvcClient.
type MockExternalSvcClientMockRecorder struct {
	mock *MockExternalSvcClient
}

// NewMockExternalSvcClient creates a new mock instance.
func NewMockExternalSvcClient(ctrl *gomock.Controller) *MockExternalSvcClient {
	mock := &MockExternalSvcClient{ctrl: ctrl}
	mock.recorder = &MockExternalSvcClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExternalSvcClient) EXPECT() *MockExternalSvcClientMockRecorder {
	return m.recorder
}

// CreateQuote mocks base method.
func (m *MockExternalSvcClient) CreateQuote(ctx context.Context, data *types.CreateQuoteRequestData) (*types.CreateQuoteResponseData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateQuote", ctx, data)
	ret0, _ := ret[0].(*types.CreateQuoteResponseData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateQuote indicates an expected call of CreateQuote.
func (mr *MockExternalSvcClientMockRecorder) CreateQuote(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateQuote", reflect.TypeOf((*MockExternalSvcClient)(nil).CreateQuote), ctx, data)
}

// MockQuoteRepo is a mock of QuoteRepo interface.
type MockQuoteRepo struct {
	ctrl     *gomock.Controller
	recorder *MockQuoteRepoMockRecorder
}

// MockQuoteRepoMockRecorder is the mock recorder for MockQuoteRepo.
type MockQuoteRepoMockRecorder struct {
	mock *MockQuoteRepo
}

// NewMockQuoteRepo creates a new mock instance.
func NewMockQuoteRepo(ctrl *gomock.Controller) *MockQuoteRepo {
	mock := &MockQuoteRepo{ctrl: ctrl}
	mock.recorder = &MockQuoteRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockQuoteRepo) EXPECT() *MockQuoteRepoMockRecorder {
	return m.recorder
}

// Read mocks base method.
func (m *MockQuoteRepo) Read(ctx context.Context, id, accountID string) (*types.FullQuoteData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", ctx, id, accountID)
	ret0, _ := ret[0].(*types.FullQuoteData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockQuoteRepoMockRecorder) Read(ctx, id, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockQuoteRepo)(nil).Read), ctx, id, accountID)
}

// Save mocks base method.
func (m *MockQuoteRepo) Save(ctx context.Context, quote *types.FullQuoteData) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, quote)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockQuoteRepoMockRecorder) Save(ctx, quote interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockQuoteRepo)(nil).Save), ctx, quote)
}
// Code generated by MockGen. DO NOT EDIT.
// Source: internal/core/ports/user_ports.go

// Package ports is a generated GoMock package.
package ports

import (
	reflect "reflect"

	v2 "github.com/gofiber/fiber/v2"
	gomock "github.com/golang/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// Checkout mocks base method.
func (m *MockUserRepository) Checkout(order Order) (int, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checkout", order)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Checkout indicates an expected call of Checkout.
func (mr *MockUserRepositoryMockRecorder) Checkout(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checkout", reflect.TypeOf((*MockUserRepository)(nil).Checkout), order)
}

// GetOffers mocks base method.
func (m *MockUserRepository) GetOffers() ([]Offer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOffers")
	ret0, _ := ret[0].([]Offer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOffers indicates an expected call of GetOffers.
func (mr *MockUserRepositoryMockRecorder) GetOffers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOffers", reflect.TypeOf((*MockUserRepository)(nil).GetOffers))
}

// GetStatus mocks base method.
func (m *MockUserRepository) GetStatus(id int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockUserRepositoryMockRecorder) GetStatus(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockUserRepository)(nil).GetStatus), id)
}

// Login mocks base method.
func (m *MockUserRepository) Login(email, password string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", email, password)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserRepositoryMockRecorder) Login(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserRepository)(nil).Login), email, password)
}

// Register mocks base method.
func (m *MockUserRepository) Register(username, email, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", username, email, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockUserRepositoryMockRecorder) Register(username, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserRepository)(nil).Register), username, email, password)
}

// UpdateOffers mocks base method.
func (m *MockUserRepository) UpdateOffers(food, medicine map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOffers", food, medicine)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOffers indicates an expected call of UpdateOffers.
func (mr *MockUserRepositoryMockRecorder) UpdateOffers(food, medicine interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOffers", reflect.TypeOf((*MockUserRepository)(nil).UpdateOffers), food, medicine)
}

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// Checkout mocks base method.
func (m *MockUserService) Checkout(order Order) (int, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checkout", order)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Checkout indicates an expected call of Checkout.
func (mr *MockUserServiceMockRecorder) Checkout(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checkout", reflect.TypeOf((*MockUserService)(nil).Checkout), order)
}

// GetOffers mocks base method.
func (m *MockUserService) GetOffers() ([]Offer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOffers")
	ret0, _ := ret[0].([]Offer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOffers indicates an expected call of GetOffers.
func (mr *MockUserServiceMockRecorder) GetOffers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOffers", reflect.TypeOf((*MockUserService)(nil).GetOffers))
}

// GetStatus mocks base method.
func (m *MockUserService) GetStatus(id int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockUserServiceMockRecorder) GetStatus(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockUserService)(nil).GetStatus), id)
}

// Login mocks base method.
func (m *MockUserService) Login(email, password string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", email, password)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserServiceMockRecorder) Login(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserService)(nil).Login), email, password)
}

// Register mocks base method.
func (m *MockUserService) Register(username, email, password, passConfirm string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", username, email, password, passConfirm)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockUserServiceMockRecorder) Register(username, email, password, passConfirm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserService)(nil).Register), username, email, password, passConfirm)
}

// UpdateOffers mocks base method.
func (m *MockUserService) UpdateOffers(food, medicine map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOffers", food, medicine)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOffers indicates an expected call of UpdateOffers.
func (mr *MockUserServiceMockRecorder) UpdateOffers(food, medicine interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOffers", reflect.TypeOf((*MockUserService)(nil).UpdateOffers), food, medicine)
}

// MockUserHandlers is a mock of UserHandlers interface.
type MockUserHandlers struct {
	ctrl     *gomock.Controller
	recorder *MockUserHandlersMockRecorder
}

// MockUserHandlersMockRecorder is the mock recorder for MockUserHandlers.
type MockUserHandlersMockRecorder struct {
	mock *MockUserHandlers
}

// NewMockUserHandlers creates a new mock instance.
func NewMockUserHandlers(ctrl *gomock.Controller) *MockUserHandlers {
	mock := &MockUserHandlers{ctrl: ctrl}
	mock.recorder = &MockUserHandlersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserHandlers) EXPECT() *MockUserHandlersMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *MockUserHandlers) Login(c *v2.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// Login indicates an expected call of Login.
func (mr *MockUserHandlersMockRecorder) Login(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserHandlers)(nil).Login), c)
}

// Register mocks base method.
func (m *MockUserHandlers) Register(c *v2.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockUserHandlersMockRecorder) Register(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserHandlers)(nil).Register), c)
}

// MockOfferHandlers is a mock of OfferHandlers interface.
type MockOfferHandlers struct {
	ctrl     *gomock.Controller
	recorder *MockOfferHandlersMockRecorder
}

// MockOfferHandlersMockRecorder is the mock recorder for MockOfferHandlers.
type MockOfferHandlersMockRecorder struct {
	mock *MockOfferHandlers
}

// NewMockOfferHandlers creates a new mock instance.
func NewMockOfferHandlers(ctrl *gomock.Controller) *MockOfferHandlers {
	mock := &MockOfferHandlers{ctrl: ctrl}
	mock.recorder = &MockOfferHandlersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOfferHandlers) EXPECT() *MockOfferHandlersMockRecorder {
	return m.recorder
}

// Checkout mocks base method.
func (m *MockOfferHandlers) Checkout(c *v2.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Checkout", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// Checkout indicates an expected call of Checkout.
func (mr *MockOfferHandlersMockRecorder) Checkout(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Checkout", reflect.TypeOf((*MockOfferHandlers)(nil).Checkout), c)
}

// GetOffers mocks base method.
func (m *MockOfferHandlers) GetOffers(c *v2.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOffers", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetOffers indicates an expected call of GetOffers.
func (mr *MockOfferHandlersMockRecorder) GetOffers(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOffers", reflect.TypeOf((*MockOfferHandlers)(nil).GetOffers), c)
}

// GetStatus mocks base method.
func (m *MockOfferHandlers) GetStatus(c *v2.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatus", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetStatus indicates an expected call of GetStatus.
func (mr *MockOfferHandlersMockRecorder) GetStatus(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatus", reflect.TypeOf((*MockOfferHandlers)(nil).GetStatus), c)
}

// UpdateOffers mocks base method.
func (m *MockOfferHandlers) UpdateOffers(c *v2.Ctx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOffers", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOffers indicates an expected call of UpdateOffers.
func (mr *MockOfferHandlersMockRecorder) UpdateOffers(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOffers", reflect.TypeOf((*MockOfferHandlers)(nil).UpdateOffers), c)
}

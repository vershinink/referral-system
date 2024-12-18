// Code generated by mockery v2.47.0. DO NOT EDIT.

package mocks

import (
	context "context"
	models "referral-rest-api/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// Database is an autogenerated mock type for the Database type
type Database struct {
	mock.Mock
}

// AccessCodeByID provides a mock function with given fields: ctx, userID
func (_m *Database) AccessCodeByID(ctx context.Context, userID int64) (string, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for AccessCodeByID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (string, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) string); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CodeByEmail provides a mock function with given fields: ctx, email
func (_m *Database) CodeByEmail(ctx context.Context, email string) (models.RefCode, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for CodeByEmail")
	}

	var r0 models.RefCode
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.RefCode, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.RefCode); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(models.RefCode)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CodeByID provides a mock function with given fields: ctx, userID
func (_m *Database) CodeByID(ctx context.Context, userID int64) (models.RefCode, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for CodeByID")
	}

	var r0 models.RefCode
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (models.RefCode, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) models.RefCode); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(models.RefCode)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateCode provides a mock function with given fields: ctx, code
func (_m *Database) CreateCode(ctx context.Context, code models.RefCode) (int64, error) {
	ret := _m.Called(ctx, code)

	if len(ret) == 0 {
		panic("no return value specified for CreateCode")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, models.RefCode) (int64, error)); ok {
		return rf(ctx, code)
	}
	if rf, ok := ret.Get(0).(func(context.Context, models.RefCode) int64); ok {
		r0 = rf(ctx, code)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, models.RefCode) error); ok {
		r1 = rf(ctx, code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateReferral provides a mock function with given fields: ctx, user
func (_m *Database) CreateReferral(ctx context.Context, user models.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateReferral")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *Database) CreateUser(ctx context.Context, user models.User) error {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.User) error); ok {
		r0 = rf(ctx, user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCode provides a mock function with given fields: ctx, userID
func (_m *Database) DeleteCode(ctx context.Context, userID int64) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCode")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UserByEmail provides a mock function with given fields: ctx, email
func (_m *Database) UserByEmail(ctx context.Context, email string) (models.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for UserByEmail")
	}

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (models.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) models.User); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UsersByReferrerID provides a mock function with given fields: ctx, refID
func (_m *Database) UsersByReferrerID(ctx context.Context, refID int64) ([]models.User, error) {
	ret := _m.Called(ctx, refID)

	if len(ret) == 0 {
		panic("no return value specified for UsersByReferrerID")
	}

	var r0 []models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]models.User, error)); ok {
		return rf(ctx, refID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []models.User); ok {
		r0 = rf(ctx, refID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, refID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewDatabase creates a new instance of Database. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *Database {
	mock := &Database{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

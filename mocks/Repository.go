// Code generated by mockery v2.31.1. DO NOT EDIT.

package mocks

import (
	entity "todolist_gin_gorm/internal/model/entity"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Create provides a mock function with given fields: title, description
func (_m *Repository) Create(title string, description string) (*entity.Todos, error) {
	ret := _m.Called(title, description)

	var r0 *entity.Todos
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*entity.Todos, error)); ok {
		return rf(title, description)
	}
	if rf, ok := ret.Get(0).(func(string, string) *entity.Todos); ok {
		r0 = rf(title, description)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Todos)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(title, description)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: user
func (_m *Repository) CreateUser(user *entity.Users) error {
	ret := _m.Called(user)

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.Users) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: todoID
func (_m *Repository) Delete(todoID int64) (int64, error) {
	ret := _m.Called(todoID)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (int64, error)); ok {
		return rf(todoID)
	}
	if rf, ok := ret.Get(0).(func(int64) int64); ok {
		r0 = rf(todoID)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(todoID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindUserByEmail provides a mock function with given fields: username
func (_m *Repository) FindUserByEmail(username string) (*entity.Users, error) {
	ret := _m.Called(username)

	var r0 *entity.Users
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*entity.Users, error)); ok {
		return rf(username)
	}
	if rf, ok := ret.Get(0).(func(string) *entity.Users); ok {
		r0 = rf(username)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Users)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *Repository) GetAll() ([]entity.Todos, error) {
	ret := _m.Called()

	var r0 []entity.Todos
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]entity.Todos, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []entity.Todos); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.Todos)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetID provides a mock function with given fields: todoID
func (_m *Repository) GetID(todoID int64) (*entity.Todos, error) {
	ret := _m.Called(todoID)

	var r0 *entity.Todos
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (*entity.Todos, error)); ok {
		return rf(todoID)
	}
	if rf, ok := ret.Get(0).(func(int64) *entity.Todos); ok {
		r0 = rf(todoID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Todos)
		}
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(todoID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: todoID, updates
func (_m *Repository) Update(todoID int64, updates map[string]interface{}) (*entity.Todos, error) {
	ret := _m.Called(todoID, updates)

	var r0 *entity.Todos
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, map[string]interface{}) (*entity.Todos, error)); ok {
		return rf(todoID, updates)
	}
	if rf, ok := ret.Get(0).(func(int64, map[string]interface{}) *entity.Todos); ok {
		r0 = rf(todoID, updates)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Todos)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, map[string]interface{}) error); ok {
		r1 = rf(todoID, updates)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

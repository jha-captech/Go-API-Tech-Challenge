// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "go-api-tech-challenge/internal/models"
)

// CourseGetter is an autogenerated mock type for the CourseGetter type
type CourseGetter struct {
	mock.Mock
}

// GetCourse provides a mock function with given fields: ctx, ID
func (_m *CourseGetter) GetCourse(ctx context.Context, ID int) (models.Course, error) {
	ret := _m.Called(ctx, ID)

	if len(ret) == 0 {
		panic("no return value specified for GetCourse")
	}

	var r0 models.Course
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (models.Course, error)); ok {
		return rf(ctx, ID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) models.Course); ok {
		r0 = rf(ctx, ID)
	} else {
		r0 = ret.Get(0).(models.Course)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCourseGetter creates a new instance of CourseGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCourseGetter(t interface {
	mock.TestingT
	Cleanup(func())
}) *CourseGetter {
	mock := &CourseGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

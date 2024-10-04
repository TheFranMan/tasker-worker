package service1

import "github.com/stretchr/testify/mock"

type Mock struct {
	mock.Mock
}

func (m *Mock) UserGet(id int) (*User, error) {
	args := m.Called(id)

	if nil == args.Get(0) {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

func (m *Mock) DeleteUser(id int) error {
	args := m.Called(id)

	return args.Error(0)
}

func (m *Mock) UpdateUser(id int, email string) error {
	args := m.Called(id)

	return args.Error(0)
}

package repository

import "context"

var storage = [5]User{
	{Id: 1, Login: "test1", Password: "pass1"},
	{Id: 2, Login: "test2", Password: "pass2"},
	{Id: 3, Login: "test3", Password: "pass3"},
	{Id: 4, Login: "test4", Password: "pass4"},
	{Id: 5, Login: "test5", Password: "pass5"},
}

type InMemoryUserRepository struct{}

func (s *InMemoryUserRepository) FindByUsernameAndPassword(ctx context.Context, username, password string) (User, error) {
	for _, user := range storage {
		if user.Login == username && user.Password == password {
			return user, nil
		}
	}
	return User{}, nil
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{}
}

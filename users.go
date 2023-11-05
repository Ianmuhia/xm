package main

import (
	"context"

	"xmserver/db"
)

type User struct {
	UserName string `json:"username"`
}

// UsersService represents a service for managing user operations.
type UsersService interface {
	// FindUser retrieves a user by their username.
	// It returns the user if found, or an error if not found or an error occurred.
	FindUser(ctx context.Context, username string) (*User, error)

	// CreateUser creates a new user with the specified username.
	// It returns the created user or an error if the user creation fails.
	CreateUser(ctx context.Context, username string) (*User, error)
}

// usersService represents an implementation of the UsersService interface.
type usersService struct {
	db *db.Queries
}

// NewUsersService creates a new UsersService instance using the provided bun.DB database connection.
// It returns the newly created UsersService.
func NewUsersService(q db.DBTX) UsersService {
	return &usersService{db: db.New(q)}
}

// CreateUser inserts a new user with the specified username into the database.
// It returns the created user or an error if the insertion fails.
func (us *usersService) CreateUser(ctx context.Context, username string) (*User, error) {
	user, err := us.db.InsertUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return &User{UserName: user}, nil
}

// FindUser retrieves a user from the database by their username.
// It returns the user if found, or an error if the retrieval fails.
func (us *usersService) FindUser(ctx context.Context, username string) (*User, error) {
	data, err := us.db.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	return &User{UserName: data}, nil
}

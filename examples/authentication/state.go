package main

import (
	"errors"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrResourceNotFound = errors.New("resource not found")
)

// DB represents a handler for some sort of state.
type DB struct {
	tokenUsers    map[string]string
	userResources map[string]map[int]string
}

// state contains some actual state.
var state = &DB{
	tokenUsers: map[string]string{
		"abcde": "alice",
		"12345": "bob",
	},

	userResources: map[string]map[int]string{
		"alice": map[int]string{
			1: "foo",
			2: "bar",
		},
		"bob": map[int]string{
			3: "baz",
		},
	},
}

// validateToken returns the user corresponding to the token.
// An error is returned if the token is not recognized.
func (db *DB) validateToken(token string) (string, error) {
	user, ok := db.tokenUsers[token]
	if !ok {
		return "", ErrInvalidToken
	}

	return user, nil
}

// resource returns the resource with id for user.
// An error is returned if the resource is not found.
func (db *DB) resource(user string, id int) (string, error) {
	resources, ok := db.userResources[user]
	if !ok {
		return "", ErrResourceNotFound
	}

	resource, ok := resources[id]
	if !ok {
		return "", ErrResourceNotFound
	}

	return resource, nil
}

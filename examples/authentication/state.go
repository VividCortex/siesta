package main

import (
	"errors"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrResourceNotFound = errors.New("resource not found")
)

type state struct {
	tokenUsers    map[string]string
	userResources map[string]map[int]string
}

var DB = &state{
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

// validateToken checks
func (db *state) validateToken(token string) (string, error) {
	user, ok := db.tokenUsers[token]
	if !ok {
		return "", ErrInvalidToken
	}

	return user, nil
}

func (db *state) resource(user string, id int) (string, error) {
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

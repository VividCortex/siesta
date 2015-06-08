package main

var tokenUsers = map[string]string{
	"abcde": "alice",
	"12345": "bob",
}

var userResources = map[string]map[int]string{
	"alice": map[int]string{
		1: "foo",
		2: "bar",
	},

	"bob": map[int]string{
		3: "baz",
	},
}

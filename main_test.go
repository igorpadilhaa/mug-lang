package main

import "testing"

func TestRun(t *testing.T) {
	err := Run(`
		user = "mingau"
		another = user
		print(hello, another)
	`)

	if err != nil {
		t.Error(err)
	}
}
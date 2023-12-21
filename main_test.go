package main

import (
	"testing"

	"github.com/igorpadilhaa/mug/engine"
)

func TestRun(t *testing.T) {
	engine.SetVar("hello", "Hello")

	err := Run(`
		user = "mingau"
		another = user
		print(hello, another)
		print(string(print()))
		print(string(1234))
	`)

	if err != nil {
		t.Error(err)
	}
}
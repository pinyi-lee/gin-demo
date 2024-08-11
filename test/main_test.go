package test

import (
	"gin-demo/app/main"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	main.Setup()
	defer main.Close()
	go main.RunServer()

	r := m.Run()
	os.Exit(r)
}

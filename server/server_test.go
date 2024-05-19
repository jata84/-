package server

import (
	"fmt"
	"testing"
)

func TestNewProjectRunner(t *testing.T) {
	s := NewServer()
	s.Run()
	err := s.runner.LoadAllProject("../tests")
	if err != nil {
		fmt.Println(err)
	} else {
		s.runner.Run("project_demo")
	}
}

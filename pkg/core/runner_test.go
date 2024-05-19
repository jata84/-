package core

import (
	"testing"
)

func TestNewProjectRunner(t *testing.T) {
	r := NewRunner()
	r.LoadAllProject("../../tests/")
	p := r.projects.Get("project_demo")
	p.Run()

}

func TestNewProjectRunnerSelenium(t *testing.T) {
	r := NewRunner()
	r.LoadAllProject("../../tests/")
	p := r.projects.Get("project_selenium")
	p.Run()
}

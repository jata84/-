package core

import (
	"testing"
)

func TestNewLoaderFromFile(t *testing.T) {
	loader := NewProjectLoader()
	loader.LoadFromPath("../../tests_2/project_demo/")

}

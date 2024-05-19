package core

import (
	"testing"
)

func TestNewLoaderFromFile(t *testing.T) {
	loader := NewProjectLoader()
	loader.LoadFromPath("../../tests/project_demo/")

}

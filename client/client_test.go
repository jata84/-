package client

import (
	"testing"
)

func TestNewNodeResponse(t *testing.T) {
	client := NewClient()
	defer client.Close()
	client.ShellCommandLoad("/home/jata/workspace/Backend/goTask/tests/project_demo")

}

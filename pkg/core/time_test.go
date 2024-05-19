package core

import (
	"testing"
)

func TestNewModuleWaitWringParameters(t *testing.T) {
	name := "TestWait"
	description := "Test Description"
	nextNode := &Node{}
	var parameter *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("wrong_seconds", "10", "FLOAT"),
	})

	_, node_wait := NewNodeModuleWait("time")

	node := NewNode(name, description, parameter, nextNode, node_wait)

	if node.name != name {
		t.Errorf("Expected name to be %s, but got %s", name, node.name)
	}

	if node.Run() != nil {
		t.Errorf("Error executing run command")
	}

}

func TestNewModuleWaitRun(t *testing.T) {
	name := "TestWait"
	description := "Test Description"
	nextNode := &Node{}
	var parameter *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("seconds", "1", "FLOAT"),
	})

	_, node_wait := NewNodeModuleWait("time")

	node := NewNode(name, description, parameter, nextNode, node_wait)

	if node.name != name {
		t.Errorf("Expected name to be %s, but got %s", name, node.name)
	}

	if node.Run() != nil {
		t.Errorf("Error executing run command")
	}

}

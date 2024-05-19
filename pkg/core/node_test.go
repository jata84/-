package core

import (
	"testing"
)

func TestNewNodeResponse(t *testing.T) {
	status := "success"
	responseDict := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	node := &Node{}
	nodeResponse := NewNodeResponse(status, node, responseDict)

	if nodeResponse.status != status {
		t.Errorf("Expected status to be %s, but got %s", status, nodeResponse.status)
	}

	if nodeResponse.node != node {
		t.Errorf("Expected node to be the same, but it's not")
	}

	if nodeResponse.response_dict["key1"] != "value1" {
		t.Errorf("Expected 'key1' in response_dict to be 'value1', but got %v", nodeResponse.response_dict["key1"])
	}
}

func TestNewNode(t *testing.T) {
	name := "TestNode"
	description := "Test Description"
	nextNode := &Node{}
	var parameter *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "value_test", "str"),
	})

	node := NewNode(name, description, parameter, nextNode, nil)

	if node.name != name {
		t.Errorf("Expected name to be %s, but got %s", name, node.name)
	}

	if node.description != description {
		t.Errorf("Expected description to be %s, but got %s", description, node.description)
	}
	/*
		if node.parameters[0]["name"] != "value1" {
			t.Errorf("Expected 'param1' in parameters to be 'value1', but got %v", node.parameters["param1"])
		}
	*/
	if node.next != nextNode {
		t.Errorf("Expected next node to be the same, but it's not")
	}

	if node.response != nil {
		t.Errorf("Expected response to be the same, but it's not")
	}

	if node.validate_parameters() != nil {
		t.Errorf("Error validating parameters")
	}
}

func TestNodeRequest(t *testing.T) {
	name := "TestNode"
	description := "Test Description"
	nextNode := &Node{}
	var parameter *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("url", "https://httpbin.org/uuid", "str"),
	})

	_, node_request := NewNodeModuleRequest("get")

	node := NewNode(name, description, parameter, nextNode, node_request)

	if node.name != name {
		t.Errorf("Expected name to be %s, but got %s", name, node.name)
	}

	if node.Run() != nil {
		t.Errorf("Error executing run command")
	}

}

func TestNodeName(t *testing.T) {
	name := "TestNode"
	node := &Node{name: name}

	result := node.Name()

	if result != name {
		t.Errorf("Expected Name() to return %s, but got %s", name, result)
	}
}

func TestNodeRunOK(t *testing.T) {
	name := "TestNode"
	description := "Test Description"
	nextNode := &Node{}
	var parameter *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "value_test", "str"),
	})

	node := NewNode(name, description, parameter, nextNode, nil)

	if node.Run() != nil {
		t.Errorf("Error executing run command")
	}
}

func TestNodeRunErrorParameters(t *testing.T) {
	name := "TestNode"
	description := "Test Description"
	nextNode := &Node{}

	var parameter *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("wrong_parameter", "value_test", "str"),
	})

	node := NewNode(name, description, parameter, nextNode, nil)

	node.Run()
	if node.status.status != NodeStatusError {
		t.Errorf("Run failed but status and Node Status its not Correct")
	}

}

func TestNodeRunErrorRun(t *testing.T) {
	name := "TestNode"
	description := "Test Description"
	nextNode := &Node{}
	var parameter *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "value_test", "str"),
	})
	node := NewNode(name, description, parameter, nextNode, nil)

	if node.Run() != nil {
		if node.status.status != NodeStatusError {
			t.Errorf("Run failed but status and Node Status its not Correct")
		}

	}
}

func TestNodeValidateParameterNotInMandatory(t *testing.T) {
	name := "TestNode"
	description := "Test Description"
	nextNode := &Node{}
	var parameter *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "value_test", "str"),
		NewNodeParameter("name_test_not_mandatory", "value_test", "str"),
	})
	node := NewNode(name, description, parameter, nextNode, nil)

	if node.Run() != nil {
		if node.status.status != NodeStatusError {
			t.Errorf("Run failed but status and Node Status its not Correct")
		}

	}
}

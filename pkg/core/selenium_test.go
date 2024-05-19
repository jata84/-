package core

import (
	"testing"
)

func TestSeleniumNewBrowser(t *testing.T) {
	name := "TestWait"
	description := "Test Description"
	datastore := NewDataStore(nil)
	project := NewProject("Test Project", datastore, nil)

	nextNode := &Node{}
	var parameter *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("browser", "1", "STRING"),
		NewNodeParameter("url", "https://efi.efiasistencia.com/#/login", "STRING"),
	})

	_, node_selenium := NewNodeModuleSeleniumService("selenium")

	node := NewNode(name, description, parameter, nextNode, node_selenium)

	node.project = project
	if node.name != name {
		t.Errorf("Expected name to be %s, but got %s", name, node.name)
	}

	if node.Run() != nil {
		t.Errorf("Error executing run command")
	}
}

func TestSeleniumNewBrowserFill(t *testing.T) {
	name := "TestWait"
	description := "Test Description"
	datastore := NewDataStore(nil)
	project := NewProject("Test Project", datastore, nil)

	nextNode := &Node{}
	var parameter *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("browser", "1", "STRING"),
		NewNodeParameter("url", "https://efi.efiasistencia.com/#/login", "STRING"),
	})

	_, node_selenium := NewNodeModuleSeleniumService("selenium")

	node := NewNode(name, description, parameter, nextNode, node_selenium)

	var parameter_second *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("xpath", "1", "STRING"),
	})

	_, node_selenium_fill := NewNodeModuleSeleniumFillInput("selenium")

	node_fill := NewNode(name, description, parameter_second, nextNode, node_selenium_fill)

	node.project = project
	node_fill.project = project
	if node.name != name {
		t.Errorf("Expected name to be %s, but got %s", name, node.name)
	}

	if node.Run() != nil {
		t.Errorf("Error executing run command")
	}

	if node_fill.Run() != nil {
		t.Errorf("Error executing run command")
	}
}

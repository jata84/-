package core

import (
	"testing"
)

func TestNewPipelineWithRequestNode(t *testing.T) {
	name := "TestNode"
	description := "Test Description"
	var parameter *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("url", "https://httpbin.org/uuid", "str"),
	})

	_, node_request := NewNodeModuleRequest("get")

	node_1 := NewNode(name, description, parameter, nil, node_request)
	node_2 := NewNode(name, description, parameter, nil, node_request)

	pipeline := NewPipeline("ss")
	pipeline.AddNode(node_1)
	pipeline.AddNode(node_2)

	pipeline.Run()

}

func TestNewPipelineWithSimpleNode(t *testing.T) {

	parameter1 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "node1", "str"),
	})
	parameter2 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "node2", "str"),
	})

	node_1 := NewNode("test_node1", "test_description", parameter1, nil, nil)
	node_2 := NewNode("test_node2", "test_description2", parameter2, nil, nil)

	pipeline := NewPipeline("Prueba Pipeline")
	pipeline.AddNode(node_1)
	pipeline.AddNode(node_2)

	pipeline.Run()

}

func TestNewPipelineWithSimpleNodeErrorParsingNode(t *testing.T) {

	parameter1 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("error_parameter", "node1", "str"),
	})
	parameter2 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "node2", "str"),
	})

	node_1 := NewNode("test_node1", "test_description", parameter1, nil, nil)
	node_2 := NewNode("test_node2", "test_description2", parameter2, nil, nil)

	pipeline := NewPipeline("Prueba Pipeline")
	pipeline.AddNode(node_1)
	pipeline.AddNode(node_2)

	err := pipeline.Run()
	if err == nil {
		t.Error("Pipeline error not dispatched")
	} else {

		if pipeline.status.status != PipeLineStatusFailed {
			t.Error("Error expected Pipeline Status not Failed")
		}

		if pipeline.status.nodeStatusList.Get(pipeline.node_running.name).status != NodeStatusError {
			t.Error("Error expected Pipeline Node Status not Failed")

		}
	}

}

func TestNewPipelineListTest(t *testing.T) {

	parameter1 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "node1", "str"),
	})
	parameter2 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "node2", "str"),
	})

	node_1 := NewNode("test_node1", "test_description", parameter1, nil, nil)
	node_2 := NewNode("test_node2", "test_description2", parameter2, nil, nil)

	pipeline := NewPipeline("Prueba Pipeline")
	pipeline.AddNode(node_1)
	pipeline.AddNode(node_2)

	project_pipeline := NewPipeLineList(nil)
	project_pipeline.Add(pipeline)
	err := project_pipeline.Run()
	if err != nil {
		t.Error(err)
	}
}

func TestNewPipelineListErrorTest(t *testing.T) {

	parameter1 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("wrong_parameters", "node1", "str"),
	})
	parameter2 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "node2", "str"),
	})

	node_1 := NewNode("test_node1", "test_description", parameter1, nil, nil)
	node_2 := NewNode("test_node2", "test_description2", parameter2, nil, nil)

	pipeline := NewPipeline("Prueba Pipeline")
	pipeline.AddNode(node_1)
	pipeline.AddNode(node_2)

	project_pipeline := NewPipeLineList(nil)
	project_pipeline.Add(pipeline)
	err := project_pipeline.Run()
	if err == nil {
		t.Error("Error expected")
	}
	if project_pipeline.status != PipeLineStatusFailed {
		t.Error("Error expected Pipeline Status not Failed")
	}

}

func TestNewPipelineWithSeleniumNode(t *testing.T) {

	parameter1 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "node1", "str"),
	})
	parameter2 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "node2", "str"),
	})

	node_1 := NewNode("test_node1", "test_description", parameter1, nil, nil)
	node_2 := NewNode("test_node2", "test_description2", parameter2, nil, nil)

	pipeline := NewPipeline("Prueba Pipeline")
	pipeline.AddNode(node_1)
	pipeline.AddNode(node_2)

	pipeline.Run()

}

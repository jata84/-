package core

import (
	"testing"
)

func TestProject(t *testing.T) {
	memory := NewDataStore(nil)
	serverConfig := make(map[string]interface{})
	serverConfig["port"] = 8888
	http := NewHttpServer(serverConfig)
	project := NewProject("test", memory, http)

	pipeline_list := NewPipeLineList(project)
	pipeline := NewPipeline("test_pipeline")
	parameter1 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "value_test", "str"),
	})
	parameter2 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "value_test", "str"),
	})
	parameter3 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "value_test", "str"),
	})

	node := NewNode("test1", "description1", parameter1, nil, nil)
	node1 := NewNode("test2", "description2", parameter2, nil, nil)
	node2 := NewNode("test3", "description3", parameter3, nil, nil)
	pipeline.AddNode(node)
	pipeline.AddNode(node1)
	pipeline.AddNode(node2)
	pipeline_list.Add(pipeline)
	project.AddPipeLineList(pipeline_list)

	project.Run()
}

func TestProjectErrorWrongParameter(t *testing.T) {
	memory := NewDataStore(nil)
	project := NewProject("test", memory, nil)

	pipeline_list := NewPipeLineList(project)
	pipeline := NewPipeline("test_pipeline")
	parameter1 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("wrong_parameter", "value_test", "str"),
	})
	parameter2 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "value_test", "str"),
	})
	parameter3 := NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "value_test", "str"),
	})

	node := NewNode("test1", "description1", parameter1, nil, nil)
	node1 := NewNode("test2", "description2", parameter2, nil, nil)
	node2 := NewNode("test3", "description3", parameter3, nil, nil)
	pipeline.AddNode(node)
	pipeline.AddNode(node1)
	pipeline.AddNode(node2)
	pipeline_list.Add(pipeline)
	project.AddPipeLineList(pipeline_list)

	err := project.Run()
	if err == nil {
		t.Errorf("Error expected")
	}

	if err != nil {
		if project.status != ProjectStatusFailed {
			t.Errorf("Project status should be error")
		}
	}
}

package core

import (
	"errors"
	"fmt"
	"net/http"
)

func NewNodeModuleHttpService(name string) (error, INodeModule) {

	switch name {
	case "endpoint":
		return NewNodeModuleHttpServiceEndpoint("NodeModuleHttpEndpoint")
	default:
		return errors.New(fmt.Sprintf("Not a valid module name %s", name)), nil
	}

}

type NodeModuleHttpServiceEndpoint struct {
	INodeModule
	name string
	node *Node
}

func NewNodeModuleHttpServiceEndpoint(name string) (error, INodeModule) {

	return nil, &NodeModuleHttpServiceEndpoint{
		name: "NodeModuleHttpEndpoint",
	}
}

func (nmr *NodeModuleHttpServiceEndpoint) SetNode(node *Node) {
	nmr.node = node
}

func (nmr *NodeModuleHttpServiceEndpoint) IsBackground() bool {
	return true
}

func (nmr *NodeModuleHttpServiceEndpoint) mandatory_parameters() *NodeParameterList {
	return NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("path", "", "STRING"),
		NewNodeParameter("method", "", "STRING"),
	})
}

func (nmr *NodeModuleHttpServiceEndpoint) get_name() string {
	return nmr.name
}

func (nmr *NodeModuleHttpServiceEndpoint) pre_run(np *NodeParameterList) error {
	return nil
}

func (nmr *NodeModuleHttpServiceEndpoint) run(np *NodeParameterList) (error, *NodeResponse) {

	/*
		path := nmr.node.parameters.get("path")
		method := nmr.node.parameters.get("method")
	*/
	hello := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	}
	http.HandleFunc("/hello", hello)
	select {}
	response := NewNodeResponse(NodeStatusOk, nil, nil)
	return nil, response

}

func (nmr *NodeModuleHttpServiceEndpoint) post_run(np *NodeParameterList) error {
	return nil
}

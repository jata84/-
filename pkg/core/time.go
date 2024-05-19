package core

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

const NodeModuleTimeName = "time"

func NewNodeModuleTime(name string) (error, INodeModule) {

	switch name {
	case "wait":
		return NewNodeModuleWait("wait")
	case "wait_until":
		return errors.New("Not implemented module"), nil
	default:
		return errors.New(fmt.Sprintf("Invalid module requests.%s", name)), nil
	}
}

type Module struct {
	name string
	node *Node
}

func (m *Module) SetNode(node *Node) {
	m.node = node
}

type NodeModuleWait struct {
	INodeModule
	Module
}

func (m *NodeModuleWait) SetNode(node *Node) {
	m.node = node
}

func NewNodeModuleWait(name string) (error, INodeModule) {
	nm := &NodeModuleWait{
		Module: Module{
			name: name,
			node: nil,
		},
	}

	return nil, nm
}

func (nmr *NodeModuleWait) mandatory_parameters() *NodeParameterList {
	return NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("seconds", nil, "FLOAT"),
	})
}

func (nmr *NodeModuleWait) get_name() string {
	return nmr.name
}

func (nmr *NodeModuleWait) pre_run(np *NodeParameterList) error {
	return nil
}

func (nmr *NodeModuleWait) run(np *NodeParameterList) (error, *NodeResponse) {

	var seconds string
	seconds = fmt.Sprintf("%v", nmr.node.parameters.get("seconds").parameter_value)

	secondsInt, err := strconv.Atoi(seconds)
	if err != nil {
		return err, nil
	}

	time.Sleep(time.Duration(secondsInt) * time.Second)

	var response_map map[string]interface{}
	response := NewNodeResponse(NodeStatusOk, nil, response_map)
	return nil, response

}

func (nmr *NodeModuleWait) post_run(np *NodeParameterList) error {
	return nil
}

func (nmr *NodeModuleWait) IsBackground() bool {
	return false
}

package core

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type NodeResponse struct {
	status        string
	node          *Node
	response_dict map[string]interface{}
}

func NewNodeResponse(status string, node *Node, response_dict map[string]interface{}) *NodeResponse {

	return &NodeResponse{
		status:        status,
		node:          node,
		response_dict: response_dict}
}
func (nr *NodeResponse) SetNode(node *Node) {
	nr.node = node
}

type NodeParameter struct {
	parameter_name  string
	parameter_value interface{}
	parameter_type  string
}

func NewNodeParameter(parameter_name string, parameter_value interface{}, parameter_type string) *NodeParameter {
	return &NodeParameter{
		parameter_name:  parameter_name,
		parameter_value: parameter_value,
		parameter_type:  parameter_type,
	}
}

type NodeParameterList struct {
	parameters map[string]*NodeParameter
}

func NewNodeParameterList(parameters []*NodeParameter) *NodeParameterList {
	if parameters != nil {
		node_parameters := make(map[string]*NodeParameter)
		for _, param := range parameters {
			node_parameters[param.parameter_name] = param
		}
		return &NodeParameterList{
			parameters: node_parameters,
		}
	} else {
		node_parameters := make(map[string]*NodeParameter)
		return &NodeParameterList{
			parameters: node_parameters,
		}
	}

}

func (np *NodeParameterList) Add(parameter *NodeParameter) error {
	if np.get(parameter.parameter_name) != nil {
		return errors.New("Parameter already exists")
	}
	np.parameters[parameter.parameter_name] = parameter
	return nil

}

func (np *NodeParameterList) get(name string) *NodeParameter {
	if val, ok := np.parameters[name]; ok {
		return val
	}
	return nil
}

func (np *NodeParameterList) list() []*NodeParameter {
	var parameters []*NodeParameter
	for _, val := range np.parameters {
		parameters = append(parameters, val)
	}
	return parameters
}

type INode interface {
	run()
}

type Node struct {
	background bool

	module               INodeModule
	name                 string
	description          string
	parameters           *NodeParameterList
	next                 *Node
	response             *NodeResponse
	mandatory_parameters *NodeParameterList
	status               *NodeStatus
	project              *Project
	pipeline             *Pipeline
}

func NewNode(name string, description string, parameters *NodeParameterList, next *Node, module INodeModule) *Node {

	var mandatory_parameters *NodeParameterList = NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("name_test", "value_test", "str"),
	})

	node := &Node{
		background:           false,
		module:               module,
		name:                 name,
		description:          description,
		parameters:           parameters,
		next:                 next,
		response:             nil,
		mandatory_parameters: mandatory_parameters,
		status:               nil,
	}

	node.status = NewNodeStatus(nil, parameters, node, NodeStatusPendant)
	if module != nil {
		module.SetNode(node)
	}
	node.background = module.IsBackground()

	return node
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) set_status(status string, response *NodeResponse) {
	n.status.setStatus(status)
	n.response = response
}

func (n *Node) run() error {

	n.set_status(NodeStatusRunning, nil)
	if n.background {

	}
	if n.module != nil {
		err, response := n.module.run(n.parameters)
		if err != nil {
			n.status.setStatus(NodeStatusError)
			return err
		} else {
			response.SetNode(n)
			n.status.setStatus(NodeStatusOk)
			n.status.setResponse(response)
		}

	} else {
		fmt.Println(n.parameters.get("name_test").parameter_value)
	}
	fmt.Println("RUN")
	return nil
}

func (n *Node) pre_run() error {

	if n.module != nil {
		return n.module.pre_run(n.parameters)
	}
	fmt.Println("PRE-RUN")
	return nil

}

func (n *Node) post_run() error {
	if n.module != nil {
		return n.module.post_run(n.parameters)
	}
	fmt.Println("POST-RUN")
	return nil
}

func (n *Node) validate_parameters() error {
	var parameter_list *NodeParameterList

	if n.module != nil {
		parameter_list = n.module.mandatory_parameters()
	} else {
		parameter_list = n.mandatory_parameters
	}

	mandatoryParamInfo := make(map[string]bool)
	var error_string strings.Builder
	error_string.WriteString("")

	// Almacena información sobre si un parámetro es obligatorio o no
	for _, param := range parameter_list.list() {
		mandatoryParamInfo[param.parameter_name] = true
	}

	for _, param := range n.parameters.list() {
		// Verifica solo los parámetros obligatorios
		if isMandatory, exists := mandatoryParamInfo[param.parameter_name]; exists && isMandatory {
			delete(mandatoryParamInfo, param.parameter_name) // Elimina para que no se repita en el siguiente bucle
		} else {
			// Puedes agregar lógica aquí para manejar los parámetros no obligatorios que no están en la lista
		}
	}

	// Verifica si quedaron parámetros obligatorios sin comprobar
	for paramName := range mandatoryParamInfo {
		fmt.Printf("[%s] Mandatory parameter '%s' is missing.\n", n.name, paramName)
		error_string.WriteString(fmt.Sprintf("[%s] Mandatory parameter '%s' is missing.\n", n.name, paramName))
	}

	if error_string.String() != "" {
		return errors.New(error_string.String())
	}

	return nil
}

func (n *Node) Run() (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("Panic recuperado en Run(%s)", r))
		}
	}()

	if err = n.validate_parameters(); err != nil {
		n.set_status(NodeStatusError, NewNodeResponse(NodeStatusError, n, nil))
		return errors.Wrap(err, fmt.Sprintf("node: %s - Validation:  ", n.name))
	}
	if err = n.pre_run(); err != nil {
		n.set_status(NodeStatusError, NewNodeResponse(NodeStatusError, n, nil))
		return errors.Wrap(err, fmt.Sprintf("node: %s - PreRun: ", n.name))
	}

	if err = n.run(); err != nil {
		n.set_status(NodeStatusError, NewNodeResponse(NodeStatusError, n, nil))

		return errors.Wrap(err, fmt.Sprintf("node: %s - Run: ", n.name))
	}

	if err = n.post_run(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("node: %s - PostRun: ", n.name))
	}
	node_response := NewNodeResponse(NodeStatusOk, n, nil)
	n.set_status(NodeStatusOk, node_response)
	return nil
}

type NodeStatus struct {
	response   *NodeResponse
	parameters *NodeParameterList
	node       *Node
	status     string
}

func NewNodeStatus(response *NodeResponse, parameters *NodeParameterList, node *Node, status string) *NodeStatus {
	return &NodeStatus{
		response:   response,
		parameters: node.parameters,
		node:       node,
		status:     status,
	}
}

func (ns *NodeStatus) setStatus(status string) {
	ns.status = status
}

func (ns *NodeStatus) setResponse(response *NodeResponse) {
	ns.response = response
}

type NodeStatusList struct {
	nodeStatusList map[string]*NodeStatus
}

func (nsl *NodeStatusList) Add(nodeStatus *NodeStatus) {
	nsl.nodeStatusList[nodeStatus.node.name] = nodeStatus
}

func (nsl *NodeStatusList) Get(name string) *NodeStatus {
	if val, ok := nsl.nodeStatusList[name]; ok {
		return val
	}
	return nil
}

func NewNodeStatusList() *NodeStatusList {
	return &NodeStatusList{
		nodeStatusList: make(map[string]*NodeStatus),
	}
}

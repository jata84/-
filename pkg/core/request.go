package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/imroc/req/v3"
)

const NodeModuleRequestName = "request"

type NodeModuleRequest struct {
	INodeModule
	name string
	node *Node
}

func NewNodeModuleRequest(name string) (error, INodeModule) {

	switch name {
	case "get":
		return nil, &NodeModuleRequest{
			name: NodeModuleRequestName,
		}
	case "post":
		return errors.New("Not implemented module"), nil
	default:
		return errors.New(fmt.Sprintf("Invalid module requests.%s", name)), nil
	}
}

func (nmr *NodeModuleRequest) SetNode(node *Node) {
	nmr.node = node
}
func (nmr *NodeModuleRequest) mandatory_parameters() *NodeParameterList {
	return NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("url", "https://httpbin.org/uuid", "STRING"),
	})
}

func (nmr *NodeModuleRequest) get_name() string {
	return nmr.name
}

func (nmr *NodeModuleRequest) pre_run(np *NodeParameterList) error {
	return nil
}

func (nmr *NodeModuleRequest) run(np *NodeParameterList) (error, *NodeResponse) {
	client := req.C()
	resp, err := client.R().
		Get(np.get("url").parameter_value.(string))
	if err != nil {
		return err, NewNodeResponse(NodeStatusError, nil, nil)
	}
	fmt.Println(resp)

	var response_map map[string]interface{}
	resp_body, err := io.ReadAll(resp.Body)
	json.Unmarshal(resp_body, &response_map)
	response := NewNodeResponse(NodeStatusOk, nil, response_map)
	return nil, response

}

func (nmr *NodeModuleRequest) post_run(np *NodeParameterList) error {
	return nil
}

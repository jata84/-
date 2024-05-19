package core

import (
	"errors"
	"fmt"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func NewNodeModuleSeleniumService(name string) (error, INodeModule) {

	switch name {
	case "start":
		return NewNodeModuleSeleniumServiceStart("start")
	case "fill_input":
		return NewNodeModuleSeleniumFillInput("fill_input")
	case "click":
		return NewNodeModuleSeleniumClick("click")
	case "wait_for_elements":
		return NewNodeModuleSeleniumWaitForElement("wait_for_elements")
	default:
		return errors.New(fmt.Sprintf("Not a valid module name %s", name)), nil
	}

}

type NodeModuleSeleniumServiceStart struct {
	INodeModule
	name string
	node *Node
}

func NewNodeModuleSeleniumServiceStart(name string) (error, INodeModule) {

	return nil, &NodeModuleSeleniumServiceStart{
		name: "NodeModuleSeleniumStart",
	}
}

func (nmr *NodeModuleSeleniumServiceStart) SetNode(node *Node) {
	nmr.node = node
}

func (nmr *NodeModuleSeleniumServiceStart) mandatory_parameters() *NodeParameterList {
	return NewNodeParameterList([]*NodeParameter{
		//NewNodeParameter("browser", "", "str"),
		NewNodeParameter("url", "", "str"),
	})
}

func (nmr *NodeModuleSeleniumServiceStart) get_name() string {
	return nmr.name
}

func (nmr *NodeModuleSeleniumServiceStart) pre_run(np *NodeParameterList) error {
	return nil
}

func (nmr *NodeModuleSeleniumServiceStart) parse_parameter_input(np *NodeParameterList) (*NodeModuleParameterFillInput, error) {

	inputs, ok := nmr.node.parameters.get("wait_until").parameter_value.(interface{})
	if !ok {
		return nil, errors.New("Error reading input parameter")
	}

	input_data, ok := inputs.(map[string]interface{})
	if !ok {
		return nil, errors.New("Error reading input parameter")
	}
	input_element := &NodeModuleParameterFillInput{
		Name:      input_data["name"].(string),
		Value:     input_data["value"].(string),
		InputType: input_data["type"].(string),
	}

	return input_element, nil

}

func (nmr *NodeModuleSeleniumServiceStart) IsBackground() bool {
	return false
}

func (nmr *NodeModuleSeleniumServiceStart) ElementPresentCondition(driver selenium.WebDriver, parameter *NodeModuleParameterFillInput) selenium.Condition {
	return func(wd selenium.WebDriver) (bool, error) {
		for {
			_, err := wd.FindElement(parameter.InputType, parameter.Name)
			if err == nil {
				return true, err
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func (nmr *NodeModuleSeleniumServiceStart) run(np *NodeParameterList) (error, *NodeResponse) {
	service, err := selenium.NewChromeDriverService("../../chromedriver", 4444)
	if err != nil {
		panic(err)
	}
	//defer service.Stop()
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"window-size=1920x1080",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"disable-gpu",
		// "--headless",  // comment out this line to see the browser
	}})

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}
	nmr.node.pipeline.pipeline_list.project.memory.Set("driver", driver)
	nmr.node.pipeline.pipeline_list.project.memory.Set("service", service)

	url := nmr.node.parameters.get("url").parameter_value
	if err := driver.Get(fmt.Sprintf("%v", url)); err != nil {
		panic(err)
	}
	wait_until := nmr.node.parameters.get("wait_until")
	if wait_until != nil {
		parameter, _ := nmr.parse_parameter_input(nmr.node.parameters)
		driver.Wait(nmr.ElementPresentCondition(driver, parameter))
	}

	response := NewNodeResponse(NodeStatusOk, nil, nil)
	return nil, response

}

func (nmr *NodeModuleSeleniumServiceStart) post_run(np *NodeParameterList) error {
	return nil
}

type NodeModuleParameterFillInput struct {
	Name      string
	Value     string
	InputType string
}

type NodeModuleSeleniumFillInput struct {
	INodeModule
	name string
	node *Node
}

func NewNodeModuleSeleniumFillInput(name string) (error, INodeModule) {

	return nil, &NodeModuleSeleniumFillInput{
		name: "NodeModuleSeleniumName",
	}
}

func (nmr *NodeModuleSeleniumFillInput) SetNode(node *Node) {
	nmr.node = node
}

func (nmr *NodeModuleSeleniumFillInput) mandatory_parameters() *NodeParameterList {
	return NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("inputs", "", "JSON"),
	})
}

func (nmr *NodeModuleSeleniumFillInput) get_name() string {
	return nmr.name
}

func (nmr *NodeModuleSeleniumFillInput) pre_run(np *NodeParameterList) error {
	return nil
}

func (nmr *NodeModuleSeleniumFillInput) parse_parameter_input(np *NodeParameterList) ([]*NodeModuleParameterFillInput, error) {
	inputs, ok := nmr.node.parameters.get("inputs").parameter_value.([]interface{})
	if !ok {
		return nil, errors.New("Error reading input parameter")
	}

	input_list := make([]*NodeModuleParameterFillInput, 0)

	for _, inputInterface := range inputs {
		input, ok := inputInterface.(map[string]interface{})
		if !ok {
			return nil, errors.New("Error reading input parameter")
		}
		input_list = append(input_list, &NodeModuleParameterFillInput{
			Name:      input["name"].(string),
			Value:     input["value"].(string),
			InputType: input["type"].(string),
		})

	}
	return input_list, nil

}

func (nmr *NodeModuleSeleniumFillInput) IsBackground() bool {
	return false
}

func (nmr *NodeModuleSeleniumFillInput) run(np *NodeParameterList) (error, *NodeResponse) {

	input_parameter, err := nmr.parse_parameter_input(np)
	if err != nil {
		panic(err)
	}

	driver := nmr.node.pipeline.pipeline_list.project.memory.Get("driver").(selenium.WebDriver)
	//parameters := nmr.node.parameters.get("inputs").parameter_value.(string)
	for _, p := range input_parameter {
		input, err := driver.FindElement(p.InputType, p.Name)
		if err != nil {
			panic(err)
		}
		if err := input.Clear(); err != nil {
			panic(err)
		}

		err = input.SendKeys(p.Value)
	}

	response := NewNodeResponse(NodeStatusOk, nil, nil)
	return nil, response

}

func (nmr *NodeModuleSeleniumFillInput) post_run(np *NodeParameterList) error {
	return nil
}

type NodeModuleSeleniumClick struct {
	INodeModule
	name string
	node *Node
}

func NewNodeModuleSeleniumClick(name string) (error, INodeModule) {

	return nil, &NodeModuleSeleniumClick{
		name: "NodeModuleSeleniumClick",
	}
}

func (nmr *NodeModuleSeleniumClick) SetNode(node *Node) {
	nmr.node = node
}

func (nmr *NodeModuleSeleniumClick) mandatory_parameters() *NodeParameterList {
	return NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("inputs", "", "JSON"),
	})
}

func (nmr *NodeModuleSeleniumClick) parse_parameter_input(np *NodeParameterList) (*NodeModuleParameterFillInput, error) {

	inputs, ok := nmr.node.parameters.get("inputs").parameter_value.(interface{})
	if !ok {
		return nil, errors.New("Error reading input parameter")
	}

	input_data, ok := inputs.(map[string]interface{})
	if !ok {
		return nil, errors.New("Error reading input parameter")
	}
	input_element := &NodeModuleParameterFillInput{
		Name:      input_data["name"].(string),
		Value:     input_data["value"].(string),
		InputType: input_data["type"].(string),
	}

	return input_element, nil

}

func (nmr *NodeModuleSeleniumClick) IsBackground() bool {
	return false
}

func (nmr *NodeModuleSeleniumClick) get_name() string {
	return nmr.name
}

func (nmr *NodeModuleSeleniumClick) pre_run(np *NodeParameterList) error {
	return nil
}

func (nmr *NodeModuleSeleniumClick) run(np *NodeParameterList) (error, *NodeResponse) {

	driver := nmr.node.pipeline.pipeline_list.project.memory.Get("driver").(selenium.WebDriver)
	input_parameter, err := nmr.parse_parameter_input(np)
	input, err := driver.FindElement(input_parameter.InputType, input_parameter.Name)
	if err != nil {
		panic(err)
	}

	err = input.Click()
	response := NewNodeResponse(NodeStatusOk, nil, nil)
	return nil, response

}

func (nmr *NodeModuleSeleniumClick) post_run(np *NodeParameterList) error {
	return nil
}

type NodeModuleSeleniumWaitForElement struct {
	INodeModule
	name string
	node *Node
}

func NewNodeModuleSeleniumWaitForElement(name string) (error, INodeModule) {

	return nil, &NodeModuleSeleniumWaitForElement{
		name: "NodeModuleSeleniumWaitForElement",
	}
}

func (nmr *NodeModuleSeleniumWaitForElement) SetNode(node *Node) {
	nmr.node = node
}

func (nmr *NodeModuleSeleniumWaitForElement) mandatory_parameters() *NodeParameterList {
	return NewNodeParameterList([]*NodeParameter{
		NewNodeParameter("inputs", "", "JSON"),
	})
}

func (nmr *NodeModuleSeleniumWaitForElement) IsBackground() bool {
	return false
}

func (nmr *NodeModuleSeleniumWaitForElement) parse_parameter_input(np *NodeParameterList) (*NodeModuleParameterFillInput, error) {

	inputs, ok := nmr.node.parameters.get("inputs").parameter_value.(interface{})
	if !ok {
		return nil, errors.New("Error reading input parameter")
	}

	input_data, ok := inputs.(map[string]interface{})
	if !ok {
		return nil, errors.New("Error reading input parameter")
	}
	input_element := &NodeModuleParameterFillInput{
		Name:      input_data["name"].(string),
		Value:     input_data["value"].(string),
		InputType: input_data["type"].(string),
	}

	return input_element, nil

}

func (nmr *NodeModuleSeleniumWaitForElement) get_name() string {
	return nmr.name
}

func (nmr *NodeModuleSeleniumWaitForElement) pre_run(np *NodeParameterList) error {
	return nil
}

func (nmr *NodeModuleSeleniumWaitForElement) run(np *NodeParameterList) (error, *NodeResponse) {

	driver := nmr.node.pipeline.pipeline_list.project.memory.Get("driver").(selenium.WebDriver)
	input_parameter, err := nmr.parse_parameter_input(np)
	input, err := driver.FindElement(input_parameter.InputType, input_parameter.Name)
	if err != nil {
		panic(err)
	}

	err = input.Click()
	response := NewNodeResponse(NodeStatusOk, nil, nil)
	return nil, response

}

func (nmr *NodeModuleSeleniumWaitForElement) post_run(np *NodeParameterList) error {
	return nil
}

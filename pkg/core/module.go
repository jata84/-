package core

import (
	"errors"
	"fmt"
	"strings"
)

type Modules struct {
	modules map[string]INodeModule
}

func NewModules(modules map[string]INodeModule) *Modules {
	return &Modules{modules}
}

func (m *Modules) Get(name string) INodeModule {
	return m.modules[name]
}

func (m *Modules) Set(name string, module INodeModule) {
	m.modules[name] = module
}

func CreateNodeModule(key string) (error, INodeModule) {

	keys := strings.Split(key, ".")

	if len(keys) != 2 {
		return errors.New("Not a valid module name"), nil
	}

	app_name := keys[0]
	module := keys[1]

	switch app_name {
	case "requests":
		return NewNodeModuleRequest(module)
	case "selenium":
		return NewNodeModuleSeleniumService(module)
	case "time":
		return NewNodeModuleTime(module)
	case "http":
		return NewNodeModuleHttpService(module)
	default:
		return errors.New(fmt.Sprintf("Not a valid module name %s", app_name)), nil
	}
}

type IModule interface {
	SetNode(node *Node)
	IsBackground() bool
}

type INodeModule interface {
	IModule
	get_name() string
	pre_run(np *NodeParameterList) error
	run(np *NodeParameterList) (error, *NodeResponse)
	post_run(np *NodeParameterList) error
	mandatory_parameters() *NodeParameterList
	//SetNode(node *Node)
}

type IModuleRegistration interface {
	CreateNodeModule(key string) (error, INodeModule)
}

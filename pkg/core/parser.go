package core

import (
	"fmt"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ProjectYaml struct {
	Name  string                 `yaml:"name"`
	Store map[string]interface{} `yaml:"store"`
	Http  map[string]interface{} `yaml:"http"`
	Run   map[string][]string    `yaml:"run"`
}

func NewProjectYamlfromYaml(data []byte) (*ProjectYaml, error) {
	var project_yaml ProjectYaml
	var err error
	project_yaml = ProjectYaml{}

	err = project_yaml.Parse(data)
	return &project_yaml, err
}

func (p *ProjectYaml) Parse(data []byte) error {
	err := yaml.Unmarshal(data, p)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProjectYaml) Validate() (error, *Project) {
	store := NewDataStore(p.Store)
	var httpServer *HttpServer = nil
	if p.Http != nil {
		httpServer = NewHttpServer(p.Http)

	}
	project := NewProject(p.Name, store, httpServer)
	project.AddRunList(p.Run)

	return nil, project
}

func NewProjectParser(data []byte) (*Project, error) {
	var project_yaml *ProjectYaml
	var err error
	project_yaml = &ProjectYaml{}

	err = project_yaml.Parse(data)
	err, project := project_yaml.Validate()

	return project, err

}

type NodeYaml struct {
	Parameters  map[string]interface{} `yaml:"parameters"`
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description"`
}

func NewNodeParser(data []byte) (error, *Node) {
	var node_yaml *NodeYaml
	var err error
	node_yaml = &NodeYaml{}

	err = node_yaml.Parse(data)
	err, node := node_yaml.Validate()

	return err, node
}

func (np *NodeYaml) Parse(data []byte) error {
	var err error

	err = yaml.Unmarshal(data, np)

	if err != nil {
		fmt.Printf("Error unmarshaling YAML: %v\n", err)
		return err
	}
	return nil
}

func (np *NodeYaml) Validate() (error, *Node) {
	return nil, nil
}

type ParameterListYaml struct {
	parameters map[string]interface{} `yaml:"parameters"`
}

func NewParameterListParser(data []byte, data_map map[string]interface{}) (error, *NodeParameterList) {
	if data != nil {
		var parameter_list_yaml *ParameterListYaml
		var err error
		parameter_list_yaml = &ParameterListYaml{}
		err = parameter_list_yaml.Parse(data)
		err, node := parameter_list_yaml.Validate()
		return err, node
	} else if data_map != nil {
		var parameter_list_yaml *ParameterListYaml
		var err error
		parameter_list_yaml = &ParameterListYaml{}
		parameter_list_yaml.parameters = data_map
		err, node := parameter_list_yaml.Validate()
		return err, node
	} else {
		return errors.New("Error Parsing data"), nil
	}

}

func (pl *ParameterListYaml) Parse(data []byte) error {
	var err error

	err = yaml.Unmarshal(data, pl)

	if err != nil {
		fmt.Printf("Error unmarshaling YAML: %v\n", err)
		return err
	}
	return nil
}

func (pl *ParameterListYaml) Validate() (error, *NodeParameterList) {
	parameter_list := NewNodeParameterList(nil)
	for k, v := range pl.parameters {
		fmt.Sprintf("%T", v)
		var node_parameter *NodeParameter
		switch v.(type) {
		case int:
			node_parameter = NewNodeParameter(k, fmt.Sprint(v), "INTEGER")
		case float64:
			node_parameter = NewNodeParameter(k, fmt.Sprint(v), "FLOAT")
		case string:
			node_parameter = NewNodeParameter(k, fmt.Sprint(v), "STRING")
		case map[string]interface{}:
			node_parameter = NewNodeParameter(k, v.(map[string]interface{}), "MAP")
		case []interface{}:
			node_parameter = NewNodeParameter(k, v.([]interface{}), "ARRAY")

		default:

			return errors.New("invalid type"), nil
		}

		parameter_list.Add(node_parameter)
	}
	return nil, parameter_list
}

type PipelineYaml struct {
	Data map[string]interface{}
}

func PipeLineParser(data []byte) (*PipeLineList, error) {
	var pipeline PipelineYaml
	var err error
	pipeline = PipelineYaml{}

	err = pipeline.Parse(data)
	pipeline_list, err := pipeline.Validate()
	return pipeline_list, err
}

func (p *PipelineYaml) Validate() (*PipeLineList, error) {
	var pipeline_list *PipeLineList = NewPipeLineList(nil)

	for pipeline_key, pipeline_value := range p.Data {
		str := fmt.Sprintf("%v", pipeline_key)
		new_pipeline := NewPipeline(str)
		for _, node_value := range pipeline_value.([]interface{}) {
			pipeline_yaml := PipelineYaml{}
			pipeline_yaml.Data = node_value.(map[string]interface{})
			err, _, name := pipeline_yaml.pop("name")
			if err == nil {
				err, key, module := pipeline_yaml.pop()
				if err != nil {
					return nil, err
				}
				err, module_run := CreateNodeModule(key)
				if err != nil {
					errors.Wrap(err, "[Parser]")
					return nil, err
				}
				parameter_data := module.(map[string]interface{})["parameters"]
				err, parameter_list := NewParameterListParser(nil, parameter_data.(map[string]interface{}))

				new_node := NewNode(fmt.Sprintf("%v", name), "", parameter_list, nil, module_run)
				if err != nil {
					return nil, err
				}
				new_pipeline.AddNode(new_node)
			}
			fmt.Printf("")

		}
		pipeline_list.Add(new_pipeline)
	}

	return pipeline_list, nil
}

func (p *PipelineYaml) Parse(data []byte) error {
	var err error
	p.Data = make(map[string]interface{})
	err = yaml.Unmarshal(data, p.Data)

	if err != nil {
		fmt.Printf("Error unmarshaling YAML: %v\n", err)
		return err
	}
	return nil
}
func (p *PipelineYaml) pop(key ...string) (error, string, interface{}) {
	if len(key) > 1 {
		return errors.New("too many keys for pop"), "", nil
	}
	if len(key) == 0 {
		for k, v := range p.Data {
			delete(p.Data, k)
			return nil, k, v
		}
	} else {
		v, ok := p.Data[key[0]]
		if ok {
			delete(p.Data, key[0])
			return nil, key[0], v
		}
		return errors.New(fmt.Sprintf("key %s not found", key)), key[0], v
	}
	return errors.New("key not found"), "", nil
}

package core

import (
	"os"
	"testing"
)

func TestNewProjectYamlWithYAMLFile(t *testing.T) {
	// Abrir el archivo YAML
	yamlString, err := os.ReadFile("../../tests/project_demo/project.yaml")
	//yamlFile, err := os.Open("../../tests/project_demo/project.yaml")
	if err != nil {
		t.Fatalf("Error al abrir el archivo YAML: %v", err)
	}
	project, err := NewProjectYamlfromYaml(yamlString)
	project.Name = "Project Demo"

	if project.Name != "Project Demo" {
		t.Errorf("Expected nil, but got '%v'", project.Name)
	}
}

func TestNewPipelineWithYAMLFile(t *testing.T) {
	yamlString, err := os.ReadFile("../../tests/project_demo/pipelines/request_pipelines.yaml")
	if err != nil {
		t.Fatalf("Error al abrir el archivo YAML: %v", err)
	}
	pipeline, err := PipeLineParser(yamlString)
	if pipeline == nil {
		t.Errorf("Expected pipeline but error occurred")
	}

}

/*Modules Parser*/

func TestParserNodeWait(t *testing.T) {

	yamlString, err := os.ReadFile("../../tests/project_demo/pipelines/wait_pipeline.yaml")
	if err != nil {
		t.Fatalf("Error al abrir el archivo YAML: %v", err)
	}
	pipeline, err := PipeLineParser(yamlString)
	if pipeline == nil {
		t.Errorf("Expected pipeline but error occurred")
	}
	if pipeline.Run() != nil {
		t.Errorf("Expected nil but error occurred")
	}

}

func TestParserNodeSelenium(t *testing.T) {

	yamlString, err := os.ReadFile("../../tests/project_demo/pipelines/selenium_pipelines.yaml")
	if err != nil {
		t.Fatalf("Error al abrir el archivo YAML: %v", err)
	}
	pipeline, err := PipeLineParser(yamlString)
	if pipeline == nil {
		t.Errorf("Expected pipeline but error occurred")
	}
	if pipeline.Run() != nil {
		t.Errorf("Expected nil but error occurred")
	}

}

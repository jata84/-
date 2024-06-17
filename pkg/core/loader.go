package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type ProjectLoader struct {
	project_validation  bool
	pipeline_validation bool
	run_validation      bool
}

func NewProjectLoader() *ProjectLoader {
	return &ProjectLoader{
		project_validation:  false,
		pipeline_validation: false,
		run_validation:      false,
	}
}

func (p *ProjectLoader) LoadFromPath(path string) (*Project, error) {
	var project *Project
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() && file.Name() == "project.yaml" {
			filePath := filepath.Join(path, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, errors.New("Error loading project.yaml")
			} else {
				project, _ = NewProjectParser(content)
				fmt.Printf(project.name)
				p.project_validation = true
			}
			fmt.Printf("Contenido de project.yaml:\n%s\n", string(content))
		}
	}
	pipelines_path := filepath.Join(path, "pipelines")
	pipeline_files, err := os.ReadDir(pipelines_path)
	if err != nil {
		log.Fatal(err)
	}
	for _, pipeline_file := range pipeline_files {
		pipeline_file := filepath.Join(pipelines_path, pipeline_file.Name())
		content, err := os.ReadFile(pipeline_file)
		if err != nil {
			return nil, errors.New("Error loading project.yaml")
		} else {
			parser, err := PipeLineParser(content)
			if err != nil {
				return nil, err
			}
			project.AddPipeLineList(parser)
		}
		fmt.Printf("Contenido de project.yaml:\n%s\n", string(content))

	}
	/*
		for _, pipeline_run := range project.run {
			project.run_pipeline.Add(project.pipelines.pipeline_list.Get(pipeline_run))
		}
	*/
	return project, nil

}

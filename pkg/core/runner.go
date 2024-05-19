package core

import (
	"errors"
	"fmt"
	"os"
)

type Runner struct {
	projects *ProjectList
}

func NewRunner() *Runner {
	return &Runner{
		projects: NewProjectList(),
	}
}

func (r *Runner) Run(name string) error {
	project := r.projects.Get(name)
	if project == nil {
		return errors.New("Project not found")
	}
	project.Run()
	return nil
}

func (r *Runner) Stop(name string) {
}

func (r *Runner) StopAll() {
}

func (r *Runner) Status(name string) {

}

func (r *Runner) List() map[string]*Project {
	return r.projects.GetAll()
}

func (r *Runner) ListStatus() map[string]string {
	status := make(map[string]string)
	for _, project := range r.projects.GetAll() {
		status[project.name] = project.status
	}
	return status
}

func (r *Runner) LoadProject(project_name string) {
	project_loader := NewProjectLoader()
	project, err := project_loader.LoadFromPath(project_name)
	if err != nil {
		fmt.Println("Error")
	}
	r.projects.Add(project)

}

func (r *Runner) LoadAllProject(parameter_path string) error {

	archivos, err := os.ReadDir(parameter_path)
	if err != nil {
		fmt.Println("Error reading projects folder", err)
		return err
	}

	for _, archivo := range archivos {
		if archivo.IsDir() {
			nombreCarpeta := archivo.Name()
			project_loader := NewProjectLoader()
			project, err := project_loader.LoadFromPath(fmt.Sprintf("%s/%s", parameter_path, nombreCarpeta))
			if err != nil {
				return err
			}
			r.projects.Add(project)

		}
	}
	return nil
}

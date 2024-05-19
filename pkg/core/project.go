package core

import (
	"fmt"
	"net/http"
)

type HttpServer struct {
	port int
}

func NewHttpServer(parameter map[string]interface{}) *HttpServer {
	port := parameter["port"].(int)
	return &HttpServer{
		port: port,
	}
}

type Project struct {
	name         string
	pipelines    *ProjectPipelines
	store        *DataStore
	memory       *DataStore
	http         *HttpServer
	run          *RunBook
	run_pipeline *PipeLineList
	status       string
}

func NewProject(name string, memory *DataStore, http *HttpServer) *Project {

	project := &Project{
		name:   name,
		store:  NewDataStore(nil),
		memory: memory,
		http:   http,
		run:    nil,

		status: ProjectStatusNotRunning,
	}

	if project.http != nil {
		project.HttpServer()
	}
	project.run = NewRunBook(project)
	run_pipeline := NewPipeLineList(project)
	pipelines := NewProjectPipelines(project)

	project.run_pipeline = run_pipeline
	project.pipelines = pipelines

	return project
}

func (p *Project) HttpServer() {

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%v", p.http.port), nil); err != nil {
			fmt.Println("Error starting the server:", err)
		}
	}()

}

/*
	func (p *Project) AddRun(run string) {
		p.run = append(p.run, run)
	}
*/
func (p *Project) AddRunBook(runbook *RunBook) {
	p.run = runbook
}

func (p *Project) AddPipeLine(pipeline *Pipeline) {
	p.pipelines.pipeline_list.Add(pipeline)
}

func (p *Project) AddPipeLineList(pipeline *PipeLineList) {
	for _, pipeline_list := range pipeline.pipeline_list {
		p.pipelines.pipeline_list.Add(pipeline_list)
	}

	pipeline.project = p

	/*TODO: Error control on this funcion*/
}
func (p *Project) Run() error {

	err := p.run_pipeline.Run()
	if err != nil {
		p.status = ProjectStatusFailed
	}
	return err
}

type ProjectPipelines struct {
	pipeline_list *PipeLineList
}

func NewProjectPipelines(project *Project) *ProjectPipelines {
	return &ProjectPipelines{
		pipeline_list: NewPipeLineList(project),
	}
}

func (p *ProjectPipelines) Run() {
	p.pipeline_list.Run()
}

type ProjectList struct {
	project_list map[string]*Project
}

func NewProjectList() *ProjectList {

	return &ProjectList{
		project_list: make(map[string]*Project),
	}
}

func (p *ProjectList) Add(project *Project) {
	p.project_list[project.name] = project
}

func (p *ProjectList) Get(name string) *Project {
	return p.project_list[name]
}

func (p *ProjectList) GetAll() map[string]*Project {
	return p.project_list
}

func (p *ProjectList) Run(name string) {
	p.project_list[name].Run()
}

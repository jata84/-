package core

import (
	"fmt"
	"sync"
)

type PipelineStatus struct {
	status         string
	nodeStatusList *NodeStatusList
}

func NewPipeLineStatus() *PipelineStatus {

	return &PipelineStatus{
		status:         PipeLineStatusPendant,
		nodeStatusList: NewNodeStatusList(),
	}

}

func (pls *PipelineStatus) setStatus(status string) {
	pls.status = status
}

func (pls *PipelineStatus) setNodeStatus(nodeStatus *NodeStatus) {
	pls.nodeStatusList.Add(nodeStatus)

}

type Pipeline struct {
	name          string
	node_counter  int
	node_init     *Node
	node_pointer  *Node
	node_running  *Node
	status        *PipelineStatus
	project       *Project
	pipeline_list *PipeLineList
}

func NewPipeline(name string) *Pipeline {

	return &Pipeline{
		name:          name,
		node_counter:  0,
		node_init:     nil,
		node_pointer:  nil,
		node_running:  nil,
		status:        NewPipeLineStatus(),
		project:       nil,
		pipeline_list: nil,
	}
}

func (p *Pipeline) AddNode(node *Node) {

	node.project = p.project
	node.pipeline = p
	if p.node_init == nil {
		p.node_init = node
		p.node_running = node
		p.node_pointer = node
	} else {
		p.node_pointer.next = node
		p.node_pointer = node
	}

	p.node_counter += 1

}

func (p *Pipeline) Reset() {
	p.node_running = p.node_init
}

func (p *Pipeline) Run() error {
	if p.node_running == nil {
		return nil
	}
	if p.node_running.background {
		p.Run_()
		return nil
	} else {
		err := p.node_running.Run()
		if err == nil && p.node_running.next != nil {
			p.node_running = p.node_running.next
			return p.Run()
		} else {
			p.status.setStatus(PipeLineStatusFailed)
			p.status.setNodeStatus(p.node_running.status)

		}
		return err
	}

}

func (p *Pipeline) Run_() error {
	//var wg sync.WaitGroup
	//wg.Add(1)
	go func() {

		if p.node_running == nil {
			return
		}

		err := p.node_running.Run()

		if err == nil && p.node_running.next != nil {
			p.node_running = p.node_running.next
			err = p.Run() // Use err variable to capture the error
		} else {
			p.status.setStatus(PipeLineStatusFailed)
			p.status.setNodeStatus(p.node_running.status)
		}
		if err != nil {
			fmt.Println("Error:", err)
		}
		//defer wg.Done()
	}()
	//wg.Wait()

	return nil // Return a default value here, you might want to change this based on your requirements
}

type PipeLineList struct {
	pipeline_list    map[string]*Pipeline
	project          *Project
	running_pipeline *Pipeline
	status           string
}

func NewPipeLineList(project *Project) *PipeLineList {
	return &PipeLineList{
		pipeline_list:    make(map[string]*Pipeline),
		project:          project,
		status:           PipeLineStatusPendant,
		running_pipeline: nil,
	}
}

func (pl *PipeLineList) Add(pipeline *Pipeline) {
	pipeline.pipeline_list = pl
	pipeline.project = pl.project
	pl.pipeline_list[pipeline.name] = pipeline

}

func (pl *PipeLineList) Get(name string) *Pipeline {
	return pl.pipeline_list[name]
}

func (pl *PipeLineList) _Run() error {
	var wg sync.WaitGroup
	for _, pipeline := range pl.pipeline_list {

		wg.Add(1)
		go func() {
			err := pipeline.Run()
			pl.running_pipeline = pipeline
			if err != nil {
				pl.status = PipeLineStatusFailed
				return
			}
			defer wg.Done()
		}()

	}
	wg.Wait()
	return nil
}

func (pl *PipeLineList) Run() error {
	for _, pipeline := range pl.pipeline_list {
		err := pipeline.Run()
		pl.running_pipeline = pipeline
		if err != nil {
			pl.status = PipeLineStatusFailed
			return err
		}

	}

	return nil
}

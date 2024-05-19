package server

import (
	"fmt"
	"goTask/pkg/core"
	"goTask/tcp"
	"log"
	"net"
)

var logger *log.Logger

type ProjectStatus struct {
	name   string
	status string
}

type Server struct {
	serve  *tcp.Server
	runner *core.Runner
}

func (s *Server) handleCommandLoad(addr *tcp.Addr, req *tcp.Message) error {

	uncompressZipIO(req.Data, "projects/")
	fileName, err := req.GetCommandArgs(1)
	if err != nil {
		return err
	}
	s.runner.LoadProject(fmt.Sprintf("projects/%s", fileName))

	/*
	 */

	err = s.serve.Send(*addr,
		tcp.NewResponseText("Project Loaded"),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) handleCommandList(addr *tcp.Addr, req *tcp.Message) error {

	projects := s.runner.ListStatus()

	var rows [][]string = make([][]string, 0)

	for k, v := range projects {
		var row []string = make([]string, 0)
		row = append(row, k)
		row = append(row, v)
		rows = append(rows, row)

	}

	data_table := map[string]interface{}{
		"status":       "OK",
		"description":  "Project list",
		"type":         tcp.RESPONSE_TABLE,
		"table_rows":   rows,
		"table_header": []string{"Project", "Status"},
	}
	err := s.serve.Send(*addr,
		tcp.NewResponse(data_table),
	)
	return err
}

func (s *Server) handleCommandStart(addr *tcp.Addr, req *tcp.Message) error {

	command, err := req.GetCommandArgs(1)

	if err != nil {
		return err
	}

	err = s.runner.Run(command)

	if err != nil {
		return err

	}

	err = s.serve.Send(*addr,
		tcp.NewResponseText("Project Started"),
	)
	return err
}

func (s *Server) handleCommands(addr *tcp.Addr, req *tcp.Message) error {

	switch req.GetCommand() {
	case "LOAD":
		return s.handleCommandLoad(addr, req)
	case "LIST":
		return s.handleCommandList(addr, req)
	case "RUN":
		return s.handleCommandList(addr, req)
	case "START":
		return s.handleCommandStart(addr, req)
	default:
		return nil
	}

}

func (s *Server) Run() {

	s.serve.OnConnect(func(conn *net.TCPConn, addr *tcp.Addr) {
		fmt.Println(fmt.Sprintf("one client connect, remote address=%s.", conn.RemoteAddr().String()))
	})

	s.serve.OnRecv(func(addr *tcp.Addr, req *tcp.Message) {

		if req.Type == tcp.COMMAND {
			s.handleCommands(addr, req)
		}

	})

	s.serve.OnDisconnect(func(addr *tcp.Addr) {

	})

	s.serve.Run(":8080")

}

func NewServer() *Server {
	s := &Server{
		runner: core.NewRunner(),
		serve:  tcp.NewServer(),
	}

	err := s.runner.LoadAllProject("projects")
	if err != nil {
		fmt.Println(err)
	}

	return s
}

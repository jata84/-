package client

import (
	"archive/zip"
	"bytes"
	"fmt"
	"goTask/tcp"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
)

const (
	TCP  = "tcp"
	HOST = "0.0.0.0"
	PORT = "9977"
)

type Client struct {
	Cli *tcp.Client
	Out chan int
}

func NewClient() *Client {

	client := tcp.NewClient("127.0.0.1:8080")
	out := make(chan int)

	// on receive data event
	client.OnRecv(func(recv *tcp.Message) {

		data, err := recv.GetResponse()
		if err != nil {
			fmt.Println("Error")
		} else {
			if data["type"] == tcp.RESPONSE_TABLE {
				// Append header
				t := table.NewWriter()
				header := data["table_header"].([]interface{})
				rows := data["table_rows"].([]interface{})
				t.AppendHeader(header)
				for _, r := range rows {
					t.AppendRow(r.([]interface{}))
				}

				fmt.Println(t.Render())
			}
			if data["type"] == tcp.RESPONSE_DESCRIPTION {
				fmt.Println(data["description"])
			}
		}
		out <- 1

	})

	c := &Client{
		Cli: client,
		Out: out,
	}

	return c
}

func (c *Client) compressProjectIO(srcFolder string) ([]byte, error) {
	buffer := new(bytes.Buffer)
	archive := zip.NewWriter(buffer)

	baseFolder := filepath.Base(srcFolder)

	err := filepath.Walk(srcFolder, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		relPath, err := filepath.Rel(srcFolder, filePath)
		if err != nil {
			return err
		}

		destPath := filepath.Join(baseFolder, relPath)

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		dest, err := archive.Create(destPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(dest, file)
		return err
	})

	if err != nil {
		return nil, err
	}

	err = archive.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *Client) compressProject(srcFolder string) (string, error) {
	projectFolder := filepath.Base(srcFolder)
	zipFile := fmt.Sprintf("%s.zip", projectFolder)
	zipf, err := os.Create(zipFile)
	if err != nil {
		return "", err
	}
	defer zipf.Close()

	archive := zip.NewWriter(zipf)
	defer archive.Close()

	baseFolder := filepath.Base(srcFolder)

	err = filepath.Walk(srcFolder, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		relPath, err := filepath.Rel(srcFolder, filePath)
		if err != nil {
			return err
		}

		destPath := filepath.Join(baseFolder, relPath)

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		dest, err := archive.Create(destPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(dest, file)
		return err
	})

	if err != nil {
		return "", err
	}

	return zipFile, nil
}

func (c *Client) readFileToBytes(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return fileBytes, nil
}

func (c *Client) SendProjectFile(filePath string) error {
	file_bytes, err := c.readFileToBytes(filePath)
	err = c.Cli.Send(&tcp.Message{
		Type: tcp.FILE,
		Data: file_bytes,
	})
	if err != nil {
		return err
	}
	return nil

	/*
		c.sendMessage(filePath)
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(c.conn, file)
		if err != nil {
			return err
		}

		fmt.Printf("Project %s Loaded", filePath)

		return nil
	*/
	return nil
}

func (c *Client) SendCommand(data []byte, command ...string) error {
	err := c.Cli.Send(&tcp.Message{
		Type:    tcp.COMMAND,
		Command: command,
		Data:    data,
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ShellCommandLoad(srcFolder string) error {

	compressed_file, err := c.compressProjectIO(srcFolder)
	if err != nil {
		return err
	}
	projectFolder := filepath.Base(srcFolder)
	err = c.SendCommand(compressed_file, "LOAD", projectFolder)
	return err
}

func (c *Client) ShellCommandList() error {
	err := c.SendCommand(nil, "LIST")
	return err
}

func (c *Client) ShellCommandStart(project string) error {
	err := c.SendCommand(nil, "START", project)
	return err
}

func (c *Client) Commands() error {

	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "shell",
				Aliases: []string{"a"},
				Usage:   "add a task to the list",
				Action: func(cCtx *cli.Context) error {
					//fmt.Println(": ", cCtx.Args().First())
					//ShellMain()
					return nil
				},
			},
			{
				Name:    "load",
				Aliases: []string{"l"},
				Usage:   "load a project into the service",
				Action: func(cCtx *cli.Context) error {
					path_parameter := cCtx.Args().Slice()[0]
					c.ShellCommandLoad(path_parameter)
					//fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "load and run a project",
				Action: func(cCtx *cli.Context) error {
					path_parameter := cCtx.Args().Slice()[0]
					c.ShellCommandLoad(path_parameter)
					//fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
			},
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "Start a project",
				Action: func(cCtx *cli.Context) error {
					path_parameter := cCtx.Args().Slice()[0]
					c.ShellCommandStart(path_parameter)
					//fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
			},

			{
				Name:    "list",
				Aliases: []string{"li"},
				Usage:   "list all projects in the service",
				Action: func(cCtx *cli.Context) error {
					c.ShellCommandList()
					//fmt.Println("completed task: ", cCtx.Args().First())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (c *Client) Close() {
	c.Cli.Close()
}

func (c *Client) Run() error {
	c.Commands()
	return nil
}

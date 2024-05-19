package server

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
)

func uncompressZip(zipFile, destFolder string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		filePath := filepath.Join(destFolder, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		zipFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zipFile.Close()

		targetFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, zipFile); err != nil {
			return err
		}
	}

	if err := os.Remove(zipFile); err != nil {
		return err
	}

	return nil
}

func uncompressZipIO(zipData []byte, destFolder string) error {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		filePath := filepath.Join(destFolder, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		zipFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zipFile.Close()

		targetFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, zipFile); err != nil {
			return err
		}
	}

	return nil
}

func listFolders(directoryPath string) []string {
	var folderNames []string

	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			folderNames = append(folderNames, file.Name())
		}
	}

	return folderNames
}

func receiveFile(conn net.Conn, fileName string) error {

	file, err := os.Create(fmt.Sprintf("%s.zip", fileName))
	if err != nil {
		return err
	}
	defer file.Close()

	written, err := io.Copy(file, conn)
	println(written)
	return err
}

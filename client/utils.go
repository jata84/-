package client

import (
	"encoding/base64"
	"io/ioutil"
)

func fileToBase64(filePath string) (string, error) {
	// Read the contents of the file.
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Encode the file data to Base64.
	base64Encoded := base64.StdEncoding.EncodeToString(fileData)

	return base64Encoded, nil
}

package provider

import (
	"fmt"
	"io"
	"os"
)

type LocalProvider struct {
	filePath string
}

func NewLocalProvider(path string) (Provider, error) {
	if path == "" {
		return nil, fmt.Errorf("local provider path is empty")
	}
	return &LocalProvider{
		filePath: path,
	}, nil
}

func (p *LocalProvider) Get() ([]byte, error) {
	if _, err := os.Stat(p.filePath); os.IsNotExist(err) {
		return nil, err
	}

	file, err := os.OpenFile(p.filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var content []byte
	content, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (p *LocalProvider) Disconnect() {}

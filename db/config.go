package db

import (
	"fmt"
	"os"
	"path"
)

const (
	MaxFileSize = 1000
)

type Config struct {
	path string
}

func (c *Config) GetValFilePath(no int) string {
	basename := fmt.Sprintf("val%04d", no)
	return path.Join(c.path, basename)
}

func (c *Config) GetLocFilePath(no int) string {
	basename := fmt.Sprintf("val%04d", no)
	return path.Join(c.path, basename)
}

func (c *Config) OpenValFile(no int, read bool) (*os.File, error) {
	path := c.GetValFilePath(no)
	if read {
		return os.Open(path)
	}
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func (c *Config) OpenLocFile(no int, read bool) (*os.File, error) {
	path := c.GetLocFilePath(no)
	if read {
		return os.Open(path)
	}
	return os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

// func (c *Config) GetFileNos(no int) ([]int, error) {
// }

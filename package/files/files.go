package files

import (
	"fmt"
	"github.com/gin-generator/ginctl/package/get"
	"os"
)

type LocalFileClient struct {
	Driver string
	Base   string
}

func NewClient() *LocalFileClient {
	driver := get.String("filesystem.driver")
	return &LocalFileClient{
		Driver: driver,
		Base:   get.String(fmt.Sprintf("filesystem.%s.base_path", driver)),
	}
}

// IsExist checks if the specified path (file or directory) exists and returns a boolean value.
func (l *LocalFileClient) IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

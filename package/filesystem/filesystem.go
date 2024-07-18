package filesystem

import (
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/gin-generator/ginctl/package/aliyun"
	"github.com/gin-generator/ginctl/package/get"
)

const (
	Local = "local"
	Oss   = "oss"
	Qi    = "qi"
	Minio = "minio"
)

type Driver interface {
	*oss.Client
}

type DriverType[T Driver] struct {
	driver T
}

func NewFileDriverType[T Driver](driver T) DriverType[T] {
	return DriverType[T]{driver: driver}
}

func (d DriverType[T]) Data() T {
	return d.driver
}

func NewClient[T Driver]() T {
	driver := get.Get("filesystem.driver")
	var client T
	switch driver {
	case Oss:
		client = NewFileDriverType(aliyun.NewOssClient()).driver
	case Local, Qi, Minio:
		panic("Not supported yet.")
	default:
		panic("The current file extension is not supported.")
	}
	return client
}

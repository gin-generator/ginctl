package cmd

import (
	"github.com/gin-generator/ginctl/package/helper"
	"os"
	"text/template"
)

type AppBase struct {
	App string
}

type Middleware struct {
	Name string
}

type Apply struct {
	App    string
	Module string
	Apply  string
}

type Deployer struct {
	App     string
	Date    string
	Version string
	Image   string
}

type Handler struct {
	Name    string
	Handler string
}

type Stub interface {
	AppBase | Apply | Deployer | Handler | Middleware
}

func CreateByStub[T Stub](filePath, stub string, stubStruct T) (err error) {
	if !helper.PathExists(filePath) {
		newFile, errs := os.Create(filePath)
		if errs != nil {
			return errs
		}
		defer newFile.Close()

		t, errs := StubData.ReadFile(stub)
		if errs != nil {
			return errs
		}

		temp, errs := template.New(filePath).Parse(string(t))
		if errs != nil {
			return errs
		}

		err = temp.Execute(newFile, stubStruct)
		return
	}
	return
}

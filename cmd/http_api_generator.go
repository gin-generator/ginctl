package cmd

import (
	"bytes"
	"fmt"
	"github.com/gin-generator/ginctl/cmd/base"
	"github.com/gin-generator/ginctl/package/helper"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type ApiLogic struct {
	Content string
}

type Body struct {
	LowerModel string
	Apply      string
	Mod        string
	Name       string
	Handler    string
}

type Operation struct {
	Opt         string
	Description string
}

type StubCode uint

const (
	FromStubBasic StubCode = iota
	FromStubImport
	FromStubLogicFunc
	FromStubTypes
	FromStubTypeFunc
	ToLogic

	FromWsLogic
	FromWsLogicRegister
	FromWsLogicImport
)

var StubMap = map[StubCode]string{
	// http
	FromStubBasic:     "stub/http/api/basic_logic.stub",
	FromStubImport:    "stub/http/api/logic_import.stub",
	FromStubLogicFunc: "stub/http/api/logic_func.stub",
	FromStubTypes:     "stub/http/api/types.stub",
	FromStubTypeFunc:  "stub/http/api/type_func.stub",
	ToLogic:           "stub/http/api/logic.stub",
	// websocket
	FromWsLogic:         "stub/websocket/logic/logic.stub",
	FromWsLogicImport:   "stub/websocket/logic/logic_import.stub",
	FromWsLogicRegister: "stub/websocket/logic/logic_register.stub",
}

// GenLogic generate apply logic.
func GenLogic(filePath string, from, to StubCode, body *Body) (err error) {
	dir := fmt.Sprintf("%s/%s", base.Pwd, strings.TrimLeft(filepath.Dir(filePath), "/"))
	err = helper.CreateDirIfNotExist(dir)
	if err != nil {
		return
	}

	filePath = fmt.Sprintf("%s/%s", base.Pwd, strings.TrimLeft(filePath, "/"))
	if helper.PathExists(filePath) {
		return
	}

	outFile, err := os.Create(filePath)
	if err != nil {
		return
	}

	var apiLogic ApiLogic
	c, err := StubData.ReadFile(StubMap[from])
	if err != nil {
		return
	}
	apiLogic.Content = string(c)

	tt, errs := StubData.ReadFile(StubMap[to])
	if errs != nil {
		return
	}
	tmp, err := template.New(filePath).Parse(string(tt))
	if err != nil {
		return
	}
	err = tmp.Execute(outFile, apiLogic)
	if err != nil {
		return
	}

	err = outFile.Close()
	if err != nil {
		return
	}

	if body != nil {
		tmp, err = template.ParseFiles(filePath)
		if err != nil {
			return
		}

		var output bytes.Buffer
		err = tmp.Execute(&output, body)
		if err != nil {
			return
		}

		err = os.WriteFile(filePath, output.Bytes(), os.ModePerm)
		return
	}

	return
}

func DoGenOperation(filePath, opt, desc string, code StubCode, errs chan error) {
	opt = helper.Capitalize(opt)

	tt, err := StubData.ReadFile(StubMap[code])
	if err != nil {
		errs <- err
		return
	}
	content := string(tt)

	// check operation is existed.
	address := fmt.Sprintf("%s/%s", base.Pwd, strings.TrimLeft(filePath, "/"))
	source, err := helper.GetFileContent(address)
	if err != nil {
		errs <- err
		return
	}
	if strings.Contains(source, opt) {
		return
	}

	err = helper.AppendToFile(address, content)
	if err != nil {
		errs <- err
		return
	}

	operate := &Operation{
		Opt:         opt,
		Description: fmt.Sprintf("%s %s", opt, desc),
	}
	err = ExecuteContent(address, operate)
	if err != nil {
		errs <- err
	}
}

func ExecuteContent(filePath string, opt *Operation) (err error) {
	tmp, err := template.ParseFiles(filePath)
	if err != nil {
		return
	}

	var output bytes.Buffer
	err = tmp.Execute(&output, opt)
	if err != nil {
		return
	}

	err = os.WriteFile(filePath, output.Bytes(), os.ModePerm)

	return
}

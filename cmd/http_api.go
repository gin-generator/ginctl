/*
Copyright © 2024 Joey <qcz19950516@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/gin-generator/ginctl/cmd/base"
	"github.com/gin-generator/ginctl/package/console"
	"github.com/spf13/cobra"
	"strings"
	"sync"
)

type Files struct {
	Code StubCode
	Name string
}

type Opts struct {
	Name string
	Desc string
}

var (
	logic     string
	operation string
	desc      string
	curd      bool
	OptMap    = []Opts{
		{
			Name: "Index",
			Desc: "Get page list",
		},
		{
			Name: "Show",
			Desc: "Get info",
		},
		{
			Name: "Create",
			Desc: "Save one source",
		},
		{
			Name: "Update",
			Desc: "Modifying a resource",
		},
		{
			Name: "Destroy",
			Desc: "Delete a resource",
		},
	}
	files = []Files{
		{
			Code: FromStubImport,
			Name: "logic",
		},
		{
			Code: FromStubTypes,
			Name: "types",
		},
	}
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "make http api",
	Long: `Example: ginctl http api -a demo -l user -c true, CURD operation to create a resource. 
Example: ginctl http api -a demo -l user -o ping -d test, to create a single operation method.`,
	RunE: GenApi,
}

func init() {
	httpCmd.AddCommand(apiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	apiCmd.Flags().StringVarP(&logic, "logic", "l", "", "Specify logic name")
	apiCmd.Flags().StringVarP(&operation, "operation", "o", "", "Specify operation name")
	apiCmd.Flags().StringVarP(&desc, "desc", "d", "", "Specify operation description")
	apiCmd.Flags().BoolVarP(&curd, "curd", "c", false, "Specifies whether you need to generate add, delete, update and get operations for the module")
}

func GenApi(_ *cobra.Command, _ []string) (err error) {

	// generate basic logic.
	err = MakeBasic()
	if err != nil {
		console.Error(err.Error())
		return
	}

	if operation != "" && curd {
		console.Error("Custom operations cannot be specified at the same time as CURD operations.")
		return
	}

	if logic == "" {
		console.Error("invalid logic name.")
		return
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 10)
	body := &Body{
		LowerModel: strings.ToLower(logic),
		Apply:      base.App,
		Mod:        base.Mod,
	}

	wg.Add(len(files))
	for _, info := range files {
		go func(wg *sync.WaitGroup, code StubCode, file string) {
			defer wg.Done()
			filePath := fmt.Sprintf("app/http/%s/logic/%s/%s.go", body.Apply, body.LowerModel, file)
			err = GenLogic(filePath, code, ToLogic, body)
			if err != nil {
				errChan <- err
			}
		}(&wg, info.Code, info.Name)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for errs := range errChan {
		if errs != nil {
			console.Error(errs.Error())
			return
		}
	}

	// make operation.
	errs := make(chan error, 20)
	files = []Files{
		{
			Code: FromStubLogicFunc,
			Name: "logic",
		},
		{
			Code: FromStubTypeFunc,
			Name: "types",
		},
	}
	if curd {
		// CURD
		for _, opt := range OptMap {
			for _, info := range files {
				filePath := fmt.Sprintf("app/http/%s/logic/%s/%s.go", body.Apply, body.LowerModel, info.Name)
				DoGenOperation(filePath, opt.Name, opt.Desc, info.Code, errs)
			}
		}
	} else {
		// Custom
		if operation == "" {
			console.Error("invalid operation name.")
			return
		}
		for _, info := range files {
			filePath := fmt.Sprintf("app/http/%s/logic/%s/%s.go", body.Apply, body.LowerModel, info.Name)
			if desc == "" {
				desc = operation
			}
			DoGenOperation(filePath, operation, desc, info.Code, errs)
		}
	}

	close(errs)
	for err = range errs {
		if err != nil {
			console.Error(err.Error())
			return
		}
	}

	console.Success("Logic Done.")

	return
}

func MakeBasic() (err error) {
	path := fmt.Sprintf("app/http/%s/logic/logic.go", base.App)
	body := &Body{
		Apply: base.App,
		Mod:   base.Mod,
	}
	err = GenLogic(path, FromStubBasic, ToLogic, body)
	return
}

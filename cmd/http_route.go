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
	"github.com/gin-generator/ginctl/package/helper"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"text/template"
)

// routeCmd represents the route command
var routeCmd = &cobra.Command{
	Use:   "route",
	Short: "make route",
	Long:  `Example: route -a web -m http`,
	RunE:  GenRoute,
}

func init() {
	httpCmd.AddCommand(routeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// routeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// routeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GenRoute(cmd *cobra.Command, _ []string) (err error) {

	filePath := fmt.Sprintf("%s/app/http/%s/%s/route.go", base.Pwd, base.App, cmd.Name())
	err = MakeRoute(filePath)
	if err != nil {
		console.Error(err.Error())
	} else {
		console.Success(fmt.Sprintf("Create route of %s done.", base.App))
	}
	return
}

type Template struct {
	Route string
}

func MakeRoute(filePath string) (err error) {

	dir := filepath.Dir(filePath)
	err = helper.CreateDirIfNotExist(dir)
	if err != nil {
		return
	}

	if !helper.PathExists(filePath) {
		var r Template
		r.Route = strcase.ToCamel(base.App)
		newFile, errs := os.Create(filePath)
		if errs != nil {
			return errs
		}
		defer newFile.Close()

		t, errs := StubData.ReadFile("stub/http/route/route.stub")
		if errs != nil {
			return errs
		}
		temp, errs := template.New(filePath).Parse(string(t))
		if errs != nil {
			return errs
		}

		err = temp.Execute(newFile, r)
		if err != nil {
			return
		}

		imports := fmt.Sprintf("%s \"%s/app/http/%s/route\"", base.App, base.Mod, base.App)
		router := fmt.Sprintf("%s/bootstrap/route.go", base.Pwd)
		err = helper.InsertImport(router, imports, "import ", "")
		if err != nil {
			return
		}

		tt, errs := StubData.ReadFile("stub/http/route/register_route.stub")
		if errs != nil {
			return errs
		}
		content := string(tt)
		content = fmt.Sprintf(content, r.Route, base.App, r.Route)
		err = helper.AppendToFile(router, content)
		if err != nil {
			return
		}
	}
	return
}

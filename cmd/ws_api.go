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
	"strings"
)

// wsApiCmd represents the wsApi command
var wsApiCmd = &cobra.Command{
	Use:   "api",
	Short: "make websocket api",
	Long:  `Example: ginctl ws api -a example -n example`,
	RunE:  GenWsApi,
}

func init() {
	wsCmd.AddCommand(wsApiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wsApiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wsApiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	wsApiCmd.Flags().StringP("name", "n", "", "api name")
}

func GenWsApi(cmd *cobra.Command, _ []string) (err error) {

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		console.Error("name is required")
		return
	}

	// register logic to route
	err = MakeWsApi(name)
	if err != nil {
		console.Error(err.Error())
		return
	}

	console.Success("Done.")
	return
}

func MakeWsApi(name string) (err error) {

	dir := fmt.Sprintf("%s/app/websocket/%s/logic", base.Pwd, base.App)
	err = helper.CreateDirIfNotExist(dir)
	if err != nil {
		return
	}

	// check route is existed
	filePath := fmt.Sprintf("%s/%s.go", dir, strings.ToLower(name))
	err = CreateByStub(filePath, StubMap[FromWsLogic], Handler{
		Name:    strings.ToLower(name),
		Handler: strcase.ToCamel(name),
	})
	if err != nil {
		return
	}
	route := fmt.Sprintf("%s/app/websocket/%s/route/route.go", base.Pwd, base.App)

	// add route
	imports, err := StubData.ReadFile(StubMap[FromWsLogicImport])
	if err != nil {
		return
	}
	importsContent := strings.Replace(string(imports), "{{.Module}}", base.Mod, -1)
	importsContent = strings.Replace(importsContent, "{{.App}}", base.App, -1)
	err = helper.InsertImport(route, importsContent, "import ", "")
	if err != nil {
		return
	}

	register, err := StubData.ReadFile(StubMap[FromWsLogicRegister])
	if err != nil {
		return
	}
	registerContent := strings.Replace(string(register), "{{.Name}}", strings.ToLower(name), -1)
	registerContent = strings.Replace(registerContent, "{{.Handler}}", strcase.ToCamel(name), -1)
	err = helper.InsertStringInFile(route, "router := websocket.NewRouter()", "\n\t"+registerContent)

	return
}

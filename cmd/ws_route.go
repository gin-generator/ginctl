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
)

// wsRouteCmd represents the wsRoute command
var wsRouteCmd = &cobra.Command{
	Use:   "route",
	Short: "make websocket route",
	Long:  `Example: ginctl ws route -a example -n example`,
	RunE:  GenWsRoute,
}

func init() {
	wsCmd.AddCommand(wsRouteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wsRouteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wsRouteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	wsRouteCmd.Flags().StringP("name", "n", "", "route name")
}

func GenWsRoute(_ *cobra.Command, _ []string) (err error) {
	err = MakeWsRoute()
	if err != nil {
		console.Error(err.Error())
		return
	}
	return
}

func MakeWsRoute() (err error) {
	dir := fmt.Sprintf("%s/app/websocket/%s/route", base.Pwd, base.App)
	err = helper.CreateDirIfNotExist(dir)
	if err != nil {
		return
	}
	filePath := fmt.Sprintf("%s/route.go", dir)
	stub := "stub/websocket/route/route.stub"
	app := Apply{
		Apply: strcase.ToCamel(base.App),
	}
	err = CreateByStub(filePath, stub, app)
	return
}

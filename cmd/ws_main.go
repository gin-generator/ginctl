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

// wsMainCmd represents the wsMain command
var wsMainCmd = &cobra.Command{
	Use:   "main",
	Short: "make websocket main",
	Long:  `Example: ginctl ws main -a example`,
	RunE:  GenWebsocket,
}

func init() {
	wsCmd.AddCommand(wsMainCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// wsMainCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// wsMainCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GenWebsocket(_ *cobra.Command, _ []string) (err error) {

	err = MakeWsRoute()
	if err != nil {
		console.Error(err.Error())
		return
	}

	err = MakeWsMain()
	if err != nil {
		console.Error(err.Error())
		return
	}

	console.Success("Done.")
	return
}

// MakeWsMain make an application of websocket
func MakeWsMain() (err error) {
	dir := fmt.Sprintf("%s/app/websocket/%s", base.Pwd, base.App)
	err = helper.CreateDirIfNotExist(dir)
	if err != nil {
		return
	}
	filePath := fmt.Sprintf("%s/%s.go", dir, base.App)
	stub := "stub/websocket/websocket.stub"
	app := Apply{
		App:    strings.ToLower(base.App),
		Module: base.Mod,
		Apply:  strcase.ToCamel(base.App),
	}
	err = CreateByStub(filePath, stub, app)
	return
}

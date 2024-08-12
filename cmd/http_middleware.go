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
	"strings"
	"text/template"
)

// middlewareCmd represents the middleware command
var middlewareCmd = &cobra.Command{
	Use:   "middleware",
	Short: "make http middleware",
	Long: `Generate middleware.
Example: http middleware --app demo --name auth.
Example: http middleware --app demo --name auth --global true.`,
	RunE: GenMiddleware,
}

func init() {
	httpCmd.AddCommand(middlewareCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// middlewareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	middlewareCmd.Flags().StringP("name", "n", "", "Middleware name")
	middlewareCmd.Flags().BoolP("global", "g", false, "Is it global middleware(default false)")
}

type Ware struct {
	Middleware string
	Pkg        string
}

func GenMiddleware(c *cobra.Command, _ []string) (err error) {

	middleware, err := c.Flags().GetString("name")
	if err != nil || middleware == "" {
		console.Error("invalid middleware name.")
		return
	}

	global, err := c.Flags().GetBool("global")
	if err != nil {
		console.Error("invalid global.")
		return
	}

	var ware Ware
	lower := strings.ToLower(middleware)
	dir, pkg := "", ""
	if global {
		pkg = "middlewares"
		dir = fmt.Sprintf("%s/middleware", base.Pwd)
	} else {
		pkg = "middleware"
		dir = fmt.Sprintf("%s/app/http/%s/middleware", base.Pwd, base.App)
	}

	err = helper.CreateDirIfNotExist(dir)
	if err != nil {
		console.Error(err.Error())
		return
	}
	// check middleware is existed.
	filePath := fmt.Sprintf("%s/%s.go", dir, lower)
	if !helper.PathExists(filePath) {
		ware.Middleware = strcase.ToCamel(middleware)
		ware.Pkg = pkg

		// create middleware file.
		newFile, ers := os.Create(filePath)
		if ers != nil {
			console.Error(ers.Error())
			return
		}
		defer newFile.Close()

		tp, errs := StubData.ReadFile("stub/http/middleware/middleware.stub")
		if errs != nil {
			return
		}
		temp, errs := template.New(filePath).Parse(string(tp))
		if errs != nil {
			console.Error(errs.Error())
			return
		}
		err = temp.Execute(newFile, ware)
		if err != nil {
			console.Error(err.Error())
			return
		}

		// insert offset.
		routePath, imports, t := "", "", ""
		if global {
			routePath = fmt.Sprintf("%s/bootstrap/route.go", base.Pwd)
			imports = fmt.Sprintf("\"%s/middleware\"", base.Mod)
			t = "\t"
		} else {
			routePath = fmt.Sprintf("%s/app/http/%s/route/route.go", base.Pwd, base.App)
			imports = fmt.Sprintf("\"%s/app/http/%s/middleware\"", base.Mod, base.App)
		}
		err = helper.InsertImport(routePath, imports, "import ", "")
		if err != nil {
			console.Error(err.Error())
			return
		}
		err = helper.InsertImport(routePath, fmt.Sprintf("\t%s.%s(),", pkg, ware.Middleware), "r.Use", t)
		if err != nil {
			console.Error(err.Error())
			return
		}
		console.Success(fmt.Sprintf("Create middleware of %s done.", lower))
	}

	return
}

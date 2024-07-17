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
	"github.com/gin-generator/ginctl/build/base"
	"github.com/gin-generator/ginctl/package/console"
	"github.com/gin-generator/ginctl/package/helper"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "make an application of http,grpc,websocket",
	RunE:  GenApply,
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GenApply(_ *cobra.Command, _ []string) (err error) {

	switch base.Module {
	case base.Http:
		err = MakeHttpApp()
	default:
		console.Error("not support")
		return
	}

	return
}

func MakeHttpApp() (err error) {

	// check
	appDir := fmt.Sprintf("%s/app/%s/%s", base.Pwd, base.Module, base.App)
	if helper.PathExists(appDir) {
		console.Info("Application of `" + base.App + "` has been created.")
		return
	}

	err = os.MkdirAll(appDir, os.ModePerm)
	if err != nil {
		console.Exit(err.Error())
	}

	var wg sync.WaitGroup
	errs := make(chan error, 5)
	wg.Add(5)

	// basic logic.go
	go func(wg *sync.WaitGroup, errs chan error) {
		defer wg.Done()
		err = MakeBasic()
		if err != nil {
			errs <- err
		}
	}(&wg, errs)

	// etc
	go func(wg *sync.WaitGroup, errs chan error) {
		defer wg.Done()
		err = MakeEtc()
		if err != nil {
			errs <- err
		}
	}(&wg, errs)

	// route
	go func(wg *sync.WaitGroup, errs chan error) {
		defer wg.Done()
		filePath := fmt.Sprintf("%s/app/%s/%s/route/route.go", base.Pwd, base.Module, base.App)
		err = MakeRoute(filePath)
		if err != nil {
			errs <- err
		}
	}(&wg, errs)

	// deploy
	go func(wg *sync.WaitGroup, errs chan error) {
		defer wg.Done()
		err = MakeDeployer()
		if err != nil {
			errs <- err
		}
	}(&wg, errs)

	// main.go
	go func(wg *sync.WaitGroup, errs chan error) {
		defer wg.Done()
		err = MakeHttpMain()
		if err != nil {
			errs <- err
		}
	}(&wg, errs)

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(errs)
	}(&wg)

	for err = range errs {
		if err != nil {
			console.Error(err.Error())
		}
	}

	console.Success("Done.")

	return
}

func MakeHttpMain() (err error) {
	filePath := fmt.Sprintf("%s/app/%s/%s/%s.go", base.Pwd, base.Module, base.App, base.App)
	stub := fmt.Sprintf("stub/apply/%s/%s.stub", base.Module, base.Module)
	app := Apply{
		Module: base.Mod,
		Apply:  strcase.ToCamel(base.App),
	}
	err = CreateByStub(filePath, stub, app)

	return
}

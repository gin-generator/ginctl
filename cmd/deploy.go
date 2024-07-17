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
	"github.com/spf13/cobra"
	"sync"
	"time"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "make deploy",
	Long:  `generate deploy for application.`,
	RunE:  GenDeploy,
}

func init() {
	rootCmd.AddCommand(deployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deployCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func GenDeploy(_ *cobra.Command, _ []string) (err error) {
	err = MakeDeployer()
	if err != nil {
		console.Error(err.Error())
	}
	console.Success("Create deploy done.")
	return
}

var deployer = []string{
	"DockerFile",
	"Makefile",
	"k8s",
	"gateway",
	"certificate",
}

func MakeDeployer() (err error) {

	dir := fmt.Sprintf("%s/app/%s/%s/deploy", base.Pwd, base.Module, base.App)
	err = helper.CreateDirIfNotExist(dir)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	errs := make(chan error, len(deployer))
	wg.Add(len(deployer))

	for _, deploy := range deployer {
		go func(deploy string, wg *sync.WaitGroup, errs chan error) {
			defer wg.Done()
			filePath := fmt.Sprintf("%s/%s", dir, deploy)
			if deploy == "k8s" || deploy == "gateway" || deploy == "certificate" {
				filePath += ".yaml"
			}
			stub := fmt.Sprintf("stub/deploy/%s.stub", deploy)
			ers := CreateByStub(filePath, stub, Deployer{
				App:     base.App,
				Date:    time.Now().Format("20060102"),
				Version: helper.GetModule(base.Pwd, "go"),
				Image:   "your-registry",
			})
			if ers != nil {
				errs <- ers
			}
		}(deploy, &wg, errs)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(errs)
	}(&wg)

	for err = range errs {
		if err != nil {
			return
		}
	}

	return
}

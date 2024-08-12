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
	b "github.com/gin-generator/ginctl/cmd/base"
	"github.com/gin-generator/ginctl/package/console"
	"github.com/gin-generator/ginctl/package/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"sync"
	"text/template"
)

// modelCmd represents the model command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "make model",
	Long: `Example: ginctl model --module http --table example.
Example: ginctl model --module http --table *.`,
	RunE: GenModelStruct,
}

func init() {

	RootCmd.AddCommand(modelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//modelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	modelCmd.Flags().StringP("module", "m", "", "Specify module name(http,grpc,websocket)")
	modelCmd.Flags().StringP("table", "t", "", "Specify table name")
	//modelCmd.Flags().StringP("path", "p", "stub", "Input model struct template file path (default $HOME/stub/model/model.stub)")
}

func GenModelStruct(c *cobra.Command, _ []string) (err error) {

	module, err := c.Flags().GetString("module")
	if err != nil {
		console.Exit("module name invalid.")
		return
	}

	// check config yaml is existed
	dir := fmt.Sprintf("%s/app/%s/%s/etc", b.Pwd, module, b.App)
	config := fmt.Sprintf("%s/env.yaml", dir)
	if !helper.PathExists(config) {
		console.Error("config not existed")
		return
	}

	viper.AddConfigPath(dir)
	viper.SetConfigType("yaml")
	viper.SetConfigName("env")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err = viper.ReadInConfig(); err != nil {
		console.Error(err.Error())
		return
	}

	// init db
	b.SetupDB()

	tableName, err := c.Flags().GetString("table")
	if err != nil {
		console.Exit(err.Error())
		return
	}
	if tableName == "" {
		console.Error("table name invalid.")
		return
	}

	// get sql database.
	database := viper.GetString(fmt.Sprintf("db.%s.database", b.DB.Config.Name()))
	// get dir.
	modelDir := fmt.Sprintf("%s/model/%s", b.Pwd, database)
	err = helper.CreateDirIfNotExist(modelDir)
	if err != nil {
		console.Error(err.Error())
		return
	}

	p := "stub/model/model.stub"
	t, err := StubData.ReadFile(p)
	if err != nil {
		console.Error(err.Error())
		return
	}
	temp, err := template.New(p).Parse(string(t))
	if err != nil {
		console.Error(err.Error())
		return
	}

	tables, err := GetTables(tableName)
	if err != nil {
		console.Error(err.Error())
		return
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 10)

	for _, table := range tables {
		filePath := fmt.Sprintf("%s/%s.go", modelDir, table.TableName)
		// check table struct is existed.
		if !helper.PathExists(filePath) {
			wg.Add(1)
			go MakeTable(temp, table, filePath, &wg, errChan)
		}
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err = range errChan {
		if err != nil {
			console.Error(err.Error())
			return
		}
	}
	console.Success("Done.")

	return
}

func MakeTable(temp *template.Template, table *Table, filePath string, wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	columns, ers := GetColumn(table.TableName)
	if ers != nil {
		errChan <- ers
		return
	}

	table.Struct = GenerateStruct(table.TableName, columns)
	table.Index = helper.GetFirstRuneLower(table.TableName)

	// Handling import packages.
	pkg := ""
	if strings.Contains(table.Struct, "json.RawMessage") {
		pkg += "\"encoding/json\"\n"
	}
	if strings.Contains(table.Struct, "time.Time") {
		pkg += "\t\"github.com/gin-generator/ginctl/package/time\""
	}
	if pkg != "" {
		table.Import = fmt.Sprintf("import (\n  %s\n)", pkg)
	}

	newFile, ers := os.Create(filePath)
	if ers != nil {
		errChan <- ers
		return
	}
	defer newFile.Close()

	err := temp.Execute(newFile, table)
	if err != nil {
		errChan <- err
		return
	}
}

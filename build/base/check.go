package base

import (
	"errors"
	"fmt"
	"github.com/gin-generator/ginctl/package/helper"
	"github.com/spf13/cobra"
	"strings"
)

const (
	Http      = "http"
	Grpc      = "grpc"
	Websocket = "websocket"
)

var (
	App     string
	Module  string
	Pwd     string
	Mod     string
	modules = []string{
		Http,
		Grpc,
		Websocket,
	}
	actionBlack = map[string]bool{
		"etc":        true,
		"route":      true,
		"api":        false,
		"model":      false,
		"middleware": true,
	}
	BlackCmd = []string{
		"ginctl",
		"version",
	}
)

func Check(cmd *cobra.Command) (err error) {
	App, err = cmd.Flags().GetString("app")
	if err != nil {
		return
	}
	if App == "" {
		return errors.New("invalid app name")
	}
	App = strings.ToLower(App)
	Module, err = cmd.Flags().GetString("module")
	if err != nil || Module == "" || !helper.Contains(modules, Module) {
		return errors.New("invalid module name")
	}
	Module = strings.ToLower(Module)
	// load or create app
	if is, ok := actionBlack[cmd.Name()]; ok && is {
		appPath := fmt.Sprintf("%s/app/%s/%s/%s", Pwd, Module, App, cmd.Name())
		err = helper.CreateDirIfNotExist(appPath)
	}
	return
}

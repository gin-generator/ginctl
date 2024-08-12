package base

import (
	"errors"
	"github.com/spf13/cobra"
	"strings"
)

const (
	Http      = "http"
	Grpc      = "grpc"
	Websocket = "websocket"
)

var (
	App      string
	Pwd      string
	Mod      string
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
	return
}

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/taodev/apiman/client/http"
	"github.com/taodev/apiman/logger"
	"github.com/taodev/apiman/storage"
)

var commandRun = &cobra.Command{
	Use:     "run",
	Short:   "run api",
	PreRun:  preRun,
	PostRun: postRun,
	Run: func(cmd *cobra.Command, args []string) {
		logger.DefaultNoPrint(false)

		var err error

		api := new(http.ApiHttp)
		if err = api.Load(workDir, configPath); err != nil {
			fmt.Println(err)
			os.Exit(-1)
			return
		}

		sessionDB := storage.NewFromMemory()
		var result http.ApiResult
		defer func() {
			logger.LogYaml(result)
		}()
		if result, err = api.Do(sessionDB); err != nil {
			fmt.Println("err: ", err)
			logger.LogYaml(result)
			os.Exit(-1)
			return
		}

		runPass = result.Pass()
	},
}

func init() {
	mainCommand.AddCommand(commandRun)
}

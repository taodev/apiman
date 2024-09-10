package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/taodev/apiman/client/http"
)

var commandRun = &cobra.Command{
	Use:   "run",
	Short: "run api",
	Run: func(cmd *cobra.Command, args []string) {
		workDir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}

		api := new(http.ApiHttp)
		if err = api.Load(workDir, configPath); err != nil {
			panic(err)
		}

		if err = api.Do(); err != nil {
			panic(err)
		}

		fmt.Println(api)
	},
}

func init() {
	mainCommand.AddCommand(commandRun)
}

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/taodev/apiman/client/http"
	"github.com/taodev/apiman/storage"
)

var commandRun = &cobra.Command{
	Use:     "run",
	Short:   "run api",
	PreRun:  preRun,
	PostRun: postRun,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(workDir) <= 0 {
			workDir, err = os.Getwd()
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			workDir, err = filepath.Abs(workDir)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		if err = os.Chdir(workDir); err != nil {
			fmt.Println(err)
			return
		}

		// 判断env配置文件是否存在
		envPath := filepath.Join(workDir, ".env.yaml")
		if err = storage.LoadEnv(envPath); err != nil {
			fmt.Println(err)
			return
		}

		api := new(http.ApiHttp)
		if err = api.Load(workDir, configPath); err != nil {
			panic(err)
		}

		sessionDB := storage.NewFromMemory()
		if err = api.Do(sessionDB); err != nil {
			panic(err)
		}

		fmt.Println(api)

		fmt.Println("safe exit")
	},
}

func init() {
	mainCommand.AddCommand(commandRun)
}

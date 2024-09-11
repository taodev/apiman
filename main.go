package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/taodev/apiman/logger"
	"github.com/taodev/apiman/storage"
)

var (
	configPath string
	configName string
	workDir    string
)

var mainCommand = &cobra.Command{
	Use: "apiman",
}

func init() {
	mainCommand.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default is $HOME/.apiman.yaml)")
	mainCommand.PersistentFlags().StringVarP(&workDir, "work", "w", "", "work dir (default is $PWD)")
}

func preRun(cmd *cobra.Command, args []string) {
	var err error

	// 配置工作目录
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

	// 初始化日志
	configName = filepath.Base(configPath)
	if i := strings.LastIndex(configName, ".yaml"); i > 0 {
		configName = configName[0:i]
	}
	logger.Default(configName)

	// 设置线程数
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	if err := mainCommand.Execute(); err != nil {
		panic(err)
	}
}

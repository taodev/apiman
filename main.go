package main

import "github.com/spf13/cobra"

var (
	configPath string
	workDir    string
)

var mainCommand = &cobra.Command{
	Use: "apiman",
}

func init() {
	mainCommand.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default is $HOME/.apiman.yaml)")
	mainCommand.PersistentFlags().StringVarP(&workDir, "work", "w", "", "work dir (default is $PWD)")
}

func main() {
	if err := mainCommand.Execute(); err != nil {
		panic(err)
	}
}

// Copyright © 2018 Zain Malik
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmalik/icte/pkgs/monitor"
	"github.com/zmalik/icte/pkgs/utils"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "icte",
	Short: "if change then execute",
	Long: `icte is a small utility to monitor changes in a directory or file, and if something changes it executes the command. For example:

icte directory1/ ls directory1/

icte is intended for stateless commands.
or to reload the state of a file or directory for apps where hot reload is not supported natively.
`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	args := readArgs()
	filesToMonitor, err := utils.GetFilesToMonitor(args)
	if err != nil {
		fmt.Println(err.Error())
		RootCmd.Usage()
		return
	}
	commandToExecute, err := utils.CommandToExecute(args)
	if err != nil {
		fmt.Println(err.Error())
		RootCmd.Usage()
		return
	}
	err = monitor.Monitor(filesToMonitor, commandToExecute)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func readArgs() []string {
	argsWithoutProg := os.Args[1:]
	err := utils.ValidateArgs(argsWithoutProg)
	if err != nil {
		fmt.Println(err.Error())
		RootCmd.Usage()
	}
	return argsWithoutProg
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.icte.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	RootCmd.SetUsageTemplate("Usage:\n    icte file files/ -c command\n")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".icte")           // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("HOME")) // adding home directory as first search path
	viper.AutomaticEnv()                   // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

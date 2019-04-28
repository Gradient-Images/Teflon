// Copyright © 2019 Máté Birkás <gadfly16@gmail.com>
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

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "teflon",
	Short: "Film and TV post-production automation framework",
	Long: `Teflon is an automation framework that makes it possible
to automate several aspects of the post-production process, like file
transformations, file sharing and job distribution.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// func init() {
// 	cobra.OnInitialize(initConfig)
//
// 	// Here you will define your flags and configuration settings.
// 	// Cobra supports persistent flags, which, if defined here,
// 	// will be global for your application.
// 	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.teflon.yaml)")
//
// 	// Cobra also supports local flags, which will only run
// 	// when this action is called directly.
// 	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }

// // initConfig reads in config file and ENV variables if set.
// func initConfig() {
// 	// teflonRoot := os.Getenv("TEFLON")
// 	// fmt.Println("Teflon root:", teflonRoot)
// 	// if cfgFile != "" {
// 	// 	// Use config file from the flag.
// 	// 	viper.SetConfigFile(cfgFile)
// 	// } else {
// 	// 	// Find home directory.
// 	// 	home, err := homedir.Dir()
// 	// 	if err != nil {
// 	// 		fmt.Println(err)
// 	// 		os.Exit(1)
// 	// 	}
// 	//
// 	// 	// Search config in home directory with name ".teflon" (without extension).
// 	// 	viper.AddConfigPath(home)
// 	// 	viper.SetConfigName(".teflon")
// 	// }
// 	//
// 	// viper.AutomaticEnv() // read in environment variables that match
// 	//
// 	// // If a config file is found, read it in.
// 	// if err := viper.ReadInConfig(); err == nil {
// 	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
// 	// }
// }
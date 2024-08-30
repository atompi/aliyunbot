/*
Copyright © 2024 Atom Pi <coder.atompi@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/atompi/aliyunbot/pkg/aliyunbot/handle"
	"github.com/atompi/aliyunbot/pkg/aliyunbot/options"
	"github.com/atompi/aliyunbot/pkg/utils"
	logkit "github.com/atompi/go-kits/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "aliyunbot",
	Short:   "一个通过阿里云 OpenAPI 操作阿里云资源配置的 CLI 工具",
	Long:    `通过阿里云 OpenAPI 操作阿里云资源配置。`,
	Version: options.Version,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		opts := options.NewOptions()

		level := opts.Core.Log.Level
		path := opts.Core.Log.Path
		maxSize := opts.Core.Log.MaxSize
		maxAge := opts.Core.Log.MaxAge
		compress := opts.Core.Log.Compress
		logger := logkit.InitLogger(level, path, maxSize, maxAge, compress)
		defer logger.Sync()
		undo := zap.ReplaceGlobals(logger)
		defer undo()

		go utils.GracefulExit()

		handle.Handle(opts)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./aliyunbot.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory with name "aliyunbot" (without extension).
		viper.AddConfigPath("./")
		viper.SetConfigType("yaml")
		viper.SetConfigName("aliyunbot")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		cobra.CheckErr(err)
	}
}

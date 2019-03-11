// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"log"
	"os"

	klinkregistry "github.com/k-box/k-link-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "klinkregistry",
	Short: "Manage K-Link applications and registrants",
	Long: `K-Link-Registry allows easy registration and management of
Applications and Registrants in a K-Link Network.`,
	Version: klinkregistry.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")

	rootCmd.PersistentFlags().String("db-host", "database", "Database Hostname")
	rootCmd.PersistentFlags().Int("db-port", 3306, "Database Port")
	rootCmd.PersistentFlags().String("db-user", "kregistry", "Database User")
	rootCmd.PersistentFlags().String("db-pass", "kregistry", "Database Password")
	rootCmd.PersistentFlags().String("db-name", "kregistry", "Database Name")

	rootCmd.PersistentFlags().String("smtp-host", "", "Mail sever hostname")
	rootCmd.PersistentFlags().Int("smtp-port", 25, "Mail server port")
	rootCmd.PersistentFlags().String("smtp-user", "registry", "Mail server user")
	rootCmd.PersistentFlags().String("smtp-pass", "registry", "Mail server password")
	rootCmd.PersistentFlags().String("smtp-from", "registry@example.com", "Mail sender address")

	rootCmd.PersistentFlags().String("assets", "", "Path to serve assets from. Default: Use embedded assets")
	rootCmd.PersistentFlags().String("migrations", "", "Path that contains the database migrations to run. Default: Use embedded migrations")

	rootCmd.PersistentFlags().Bool("enable-user-registration", false, "Enable user registration. Default false")

	viper.BindPFlag("db_host", rootCmd.PersistentFlags().Lookup("db-host"))
	viper.BindPFlag("db_port", rootCmd.PersistentFlags().Lookup("db-port"))
	viper.BindPFlag("db_user", rootCmd.PersistentFlags().Lookup("db-user"))
	viper.BindPFlag("db_pass", rootCmd.PersistentFlags().Lookup("db-pass"))
	viper.BindPFlag("db_name", rootCmd.PersistentFlags().Lookup("db-name"))

	viper.BindPFlag("smtp_host", rootCmd.PersistentFlags().Lookup("smtp-host"))
	viper.BindPFlag("smtp_port", rootCmd.PersistentFlags().Lookup("smtp-port"))
	viper.BindPFlag("smtp_user", rootCmd.PersistentFlags().Lookup("smtp-user"))
	viper.BindPFlag("smtp_pass", rootCmd.PersistentFlags().Lookup("smtp-pass"))
	viper.BindPFlag("smtp_from", rootCmd.PersistentFlags().Lookup("smtp-from"))

	viper.BindPFlag("assets_dir", rootCmd.PersistentFlags().Lookup("assets"))
	viper.BindPFlag("migrations_dir", rootCmd.PersistentFlags().Lookup("migrations"))

	viper.BindPFlag("enable_user_registration", rootCmd.PersistentFlags().Lookup("enable-user-registration"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Environment flags will be called like REGISTRY_ASSETS_DIR
	viper.SetEnvPrefix("registry")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".klinkregistry" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/klinkregistry")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}

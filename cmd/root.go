// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
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

	"github.com/blang/semver"

	"github.com/entropyx/gump/configuration"
	"github.com/entropyx/gump/file"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gump",
	Short: "Bump your project version into any file",
	Long: `Sometimes there is no easy way to update the references to your sofware version. Gump will do it for you!

		Given a simple configuration file (default is $PWD/gump.yml):

			version: 0.1.0
			files:
			- path: file/test_files/yet-another-configuration-file.yml
				keys:
				- info.extra.version
				prefix: 'docker-image:'

		Just do
			gump

		If you need to bump your major version
		  gump -M

`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Use viper Get functions insead of hardcoding
		var locations []string
		version := viper.Get("version").(string)
		newVersion, err := bumpVersion(cmd, version)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		files := viper.Get("files").([]interface{})
		for _, v := range files {
			var prefix string
			if infprefix := v.(map[interface{}]interface{})["prefix"]; infprefix != nil {
				prefix = infprefix.(string)
			}
			newVersionWithPrefix := prefix + newVersion
			ls := v.(map[interface{}]interface{})["keys"].([]interface{})
			for i := 0; i < len(ls); i++ {
				locations = append(locations, ls[i].(string))
			}
			f, err := file.Read(v.(map[interface{}]interface{})["path"].(string))
			if err != nil {
				fmt.Println("error:", err)
				return
			}

			err = f.Modify(newVersionWithPrefix, locations)
			if err != nil {
				fmt.Println("error:", err)
				return
			}
		}
		newConf := &configuration.Configuration{}
		viper.Set("version", newVersion)
		viper.Unmarshal(newConf)
		err = newConf.Write(viper.ConfigFileUsed())
		if err != nil {
			fmt.Println("error while updating the gump conf file")
		}
	},
}

func bumpVersion(cmd *cobra.Command, version string) (string, error) {
	vType := "patch"
	if minor, _ := cmd.Flags().GetBool("minor"); minor == true {
		vType = "minor"
	}
	if major, _ := cmd.Flags().GetBool("major"); major == true {
		vType = "major"
	}
	if set, _ := cmd.Flags().GetBool("set"); set == true {
		vType = "set"
	}
	if force, _ := cmd.Flags().GetString("force"); force != "" {
		vType = "force"
		version = force
	}
	return bump(vType, version)
}

func bump(typ, version string) (string, error) {
	v, err := semver.Make(version)
	if err != nil {
		return "", err
	}
	switch typ {
	case "major":
		v.Major++
	case "minor":
		v.Minor++
	case "set", "force":
		fmt.Println("Set version to", version)
		return version, nil
	case "patch":
		v.Patch++
	}
	newVersion := v.String()
	fmt.Printf("Bumped version from %s to %s \n", version, newVersion)
	return newVersion, nil
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $PWD/gump.yml)")
	RootCmd.Flags().BoolP("major", "M", false, "Bump major version")
	RootCmd.Flags().BoolP("minor", "m", false, "Bump minor version")
	RootCmd.Flags().BoolP("patch", "p", false, "Bump patch version")
	RootCmd.Flags().BoolP("set", "s", false, "Set the current version in config file")
	RootCmd.Flags().StringP("force", "f", "", "Manually select a version")
	// TODO: support prereleases
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("gump") // name of config file (without extension)
	viper.AddConfigPath(".")    // adding home directory as first search path
	viper.AutomaticEnv()        // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

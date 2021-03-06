/*
Copyright © 2020 Tim_Paik <timpaik@163.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"path"
	"path/filepath"
	"pdk/pkg"
	"strings"

	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:     "install <package(s)>",
	Short:   "Install package(s)",
	Long:    `Download and install the latest package from your local repository (default)`,
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		pkg.CheckRoot()
		var isLocal bool
		var isAutoYes bool
		var repoName string
		if isAutoYes, err = cmd.Flags().GetBool("local"); err != nil {
			fmt.Println(err)
			return
		}
		if isLocal, err = cmd.Flags().GetBool("local"); err != nil {
			fmt.Println(err)
			return
		} else if isLocal {
			if len(args) == 0 {
				fmt.Println("Error: no targets specified")
			}
			for _, PATH := range args {
				fmt.Println(pkg.Indent1 + "Installing from " + PATH)
				packageName := strings.FieldsFunc(strings.TrimSuffix(filepath.Base(PATH), path.Ext(PATH)), func(r rune) bool {
					if r == '-' {
						return true
					}
					return false
				})
				fmt.Println(pkg.Indent2 + "Unpacking...")
				if err := pkg.UnpackAndCallback(PATH, packageName[0]); err != nil {
					return
				}
				fmt.Println(pkg.Indent2 + "Done!")
			}
			return
		}
		if repoName, err = cmd.Flags().GetString("repoName"); err != nil {
			fmt.Println(err)
			return
		}
		if err := pkg.Install(args, repoName, isAutoYes); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")
	installCmd.PersistentFlags().String("repoName", pkg.DefaultRepo, "Install from specified repo (repo by default)")
	installCmd.Flags().Bool("local", false, "Unzip the installation from the local (default is false)")
	installCmd.Flags().Bool("y", false, "Automatic yes to prompts (default is false)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

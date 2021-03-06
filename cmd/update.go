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
	"pdk/pkg"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update <URL>",
	Short:   "Update the repo",
	Long:    `Download and update the repo from the URL.`,
	Aliases: []string{"u"},
	Run: func(cmd *cobra.Command, args []string) {
		pkg.CheckRoot()
		var repoName string
		if repoName, err = cmd.Flags().GetString("repoName"); err != nil {
			fmt.Println(err)
			return
		}
		if len(args) == 0 {
			if err := pkg.UpdateLocal(repoName); err != nil {
				fmt.Println(err)
				return
			}
		} else if len(args) == 1 {
			if err = pkg.Update(args[0], repoName); err != nil {
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println("Error: Too many targets specified")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")
	updateCmd.PersistentFlags().String("repoName", pkg.DefaultRepo, "Update the specified repo (repo by default)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

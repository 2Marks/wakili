/*
Copyright © 2024 MOSHOOD ADEPITAN <mkoabiola95@gmail.com>
*/
package cmd

import (
	"github.com/2marks/wakili/internal/proxy"
	"github.com/spf13/cobra"
)

var baseUrlIdentifier string

// proxyCmd represents the proxy command
var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "purge data from cache",
	Long:  `purge http responses data persisted to the cache`,
	Run: func(cmd *cobra.Command, args []string) {
		proxy.PurgeCachedDataHandler(baseUrlIdentifier)
	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)

	/** start configure flags */
	purgeCmd.Flags().BoolP("help", "", false, "help for this command")
	purgeCmd.Flags().StringVarP(&baseUrlIdentifier, "baseurl", "u", "", "base url to to purge cached data for")
	/** end configure flags */

}

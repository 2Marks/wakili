/*
Copyright © 2024 MOSHOOD ADEPITAN <mkoabiola95@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/2marks/wakili/internal/proxy"
	"github.com/spf13/cobra"
)

var (
	baseUrl string
	port    string
	cache   int
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "proxy requests to specified base url",
	Long:  `proxy requests to specified base url`,
	Run: func(cmd *cobra.Command, args []string) {
		isCachingEnabled := cache >= 0
		if isCachingEnabled {
			if err := proxy.InitCache(baseUrl); err != nil {
				log.Fatal(err)
			}
		}

		proxyServer := proxy.NewServer(baseUrl, port, cache)
		proxyServer.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	/** start configure flags */
	proxyCmd.Flags().BoolP("help", "", false, "help for this command")
	proxyCmd.Flags().StringVarP(&baseUrl, "baseurl", "u", "", "base url to proxy requests from")
	proxyCmd.Flags().StringVarP(&port, "port", "p", "5000", "port proxy server listens on")
	proxyCmd.Flags().IntVarP(&cache, "cache", "c", -1, "specify time duration (in seconds) to cache requests for")

	proxyCmd.MarkFlagRequired("baseurl")
	proxyCmd.MarkFlagRequired("port")
	proxyCmd.MarkFlagsRequiredTogether("baseurl", "port")
	/** end configure flags */

}

/*
Copyright © 2024 MOSHOOD ADEPITAN <mkoabiola95@gmail.com>
*/
package cmd

import (
	"github.com/2marks/wakili/internal/proxy"
	"github.com/spf13/cobra"
)

var (
	baseUrl string
	port    string
)

// proxyCmd represents the proxy command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "proxy requests to specified base url",
	Long:  `proxy requests to specified base url`,
	Run: func(cmd *cobra.Command, args []string) {
		proxyServer := proxy.NewServer(baseUrl, port)
		proxyServer.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	/** start configure flags */
	proxyCmd.Flags().BoolP("help", "", false, "help for this command")
	proxyCmd.Flags().StringVarP(&baseUrl, "baseurl", "u", "", "base url to proxy requests from")
	proxyCmd.Flags().StringVarP(&port, "port", "p", "5000", "port proxy server listens on")
	/** end configure flags */

}

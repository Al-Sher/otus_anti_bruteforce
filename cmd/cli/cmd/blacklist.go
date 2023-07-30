package cmd

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Al-Sher/otus_anti_bruteforce/internal/httpclient"
	"github.com/Al-Sher/otus_anti_bruteforce/internal/server"
	"github.com/spf13/cobra"
)

var blackListCommand = &cobra.Command{
	Use:   "blacklist",
	Short: "Add/remove from blacklist",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		method := args[0]
		if _, ok := allowedMethods[method]; !ok {
			return errNotFoundMethod
		}

		var err error
		var b []byte
		hc := httpclient.New(cfg.Host)
		vs := url.Values{}
		vs.Set(server.IPField, network)

		if method == "add" {
			b, err = hc.Post(context.Background(), "blackList", vs)
		} else {
			b, err = hc.Delete(context.Background(), "blackList", vs)
		}

		if err := checkResponse(b, err); err != nil {
			return err
		}

		fmt.Println("Addr successfully blacklisted")
		return nil
	},
}

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/Al-Sher/otus_anti_bruteforce/internal/config"
	"github.com/Al-Sher/otus_anti_bruteforce/internal/httpclient"
	"github.com/Al-Sher/otus_anti_bruteforce/internal/server"
)

var (
	net   string
	login string
	ip    string

	errNotSpecifiedLogin = errors.New("login not specified")
	errNotSpecifiedIP    = errors.New("ip not specified")
	errNotSpecifiedNet   = errors.New("addr with mask not specified")
)

func init() {
	flag.StringVar(&net, "n", "", "Addr with mask")
	flag.StringVar(&login, "l", "", "User login")
	flag.StringVar(&ip, "i", "", "User ip")
}

func main() {
	flag.Parse()

	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	hc := httpclient.New(cfg.Host)
	ctx := context.Background()

	cmd := flag.Arg(0)
	switch cmd {
	case "reset":
		reset(ctx, hc)
	case "blacklist":
		blacklist(ctx, hc)
	case "whitelist":
		whitelist(ctx, hc)
	default:
		fmt.Println("Command not found.")
	}
}

func reset(ctx context.Context, hc httpclient.HTTPClient) {
	if login == "" {
		fmt.Println("Error: ", errNotSpecifiedLogin)
		return
	}

	if ip == "" {
		fmt.Println("Error: ", errNotSpecifiedIP)
		return
	}

	vs := url.Values{}
	vs.Add(server.LoginField, login)
	vs.Add(server.IPField, ip)

	if err := hc.Get(ctx, "reset", vs); err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

func blacklist(ctx context.Context, hc httpclient.HTTPClient) {
	if net == "" {
		fmt.Println("Error: ", errNotSpecifiedNet)
		return
	}

	vs := url.Values{}
	vs.Add(server.IPField, ip)

	if err := hc.Get(ctx, "blackList", vs); err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

func whitelist(ctx context.Context, hc httpclient.HTTPClient) {
	if net == "" {
		fmt.Println("Error: ", errNotSpecifiedNet)
		return
	}

	vs := url.Values{}
	vs.Add(server.IPField, ip)

	if err := hc.Get(ctx, "blackList", vs); err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

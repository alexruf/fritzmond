package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/alexruf/fritzmond/collector"
	"github.com/alexruf/fritzmond/config"
	"github.com/alexruf/fritzmond/fritzbox"
	"github.com/alexruf/fritzmond/http"
)

var usage bool
var cfg config.Config

func main() {
	flag.BoolVar(&usage, "help", false, "Print this help message.")
	flag.StringVar(&cfg.Url, "url", "https://fritz.box", "URL of the FRITZ!Box.")
	flag.StringVar(&cfg.Username, "username", "", "Username to authenticate with the FRITZ!Box.")
	flag.StringVar(&cfg.Password, "password", "", "Password to authenticate with the FRITZ!Box.")
	flag.BoolVar(&cfg.TlsSkipVerify, "tlsSkipVerify", true, "Skip TLS certificate validation.")
	flag.UintVar(&cfg.Interval, "interval", 10, "Interval in seconds at which data should be fetched from the FRITZ!Box.")
	flag.Parse()

	if usage {
		flag.Usage()
		return
	}
	if err := cfg.Validate(); err != nil {
		log.Printf("Error in configuration: %s", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	httpClient := http.NewHttpClient(cfg.TlsSkipVerify)
	digestAuthClient := http.NewDigestAuthClient(httpClient, cfg.Username, cfg.Password)
	fb := fritzbox.New(ctx, digestAuthClient, cfg.Url)

	col := collector.New(ctx, cfg, fb)
	wg.Add(1)
	go col.Start(wg)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	cancel()
	wg.Wait()
}

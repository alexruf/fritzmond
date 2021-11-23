package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/alexruf/fritzmond/collector"
	"github.com/alexruf/fritzmond/config"
	"github.com/alexruf/fritzmond/fritzbox"
	"github.com/alexruf/fritzmond/http"
	"github.com/alexruf/fritzmond/ui"

	"github.com/nakabonne/tstorage"
)

var usage bool
var cfg config.Config

func main() {
	flag.BoolVar(&usage, "help", false, "Print this help message.")
	flag.StringVar(&cfg.Url, "url", "https://fritz.box", "URL of the FRITZ!Box.")
	flag.StringVar(&cfg.Username, "username", "", "Username to authenticate with the FRITZ!Box.")
	flag.StringVar(&cfg.Password, "password", "", "Password to authenticate with the FRITZ!Box.")
	flag.BoolVar(&cfg.SkipTlsVerify, "skipTlsVerify", true, "Skip TLS certificate validation.")
	flag.UintVar(&cfg.Interval, "interval", 10, "Interval in seconds at which data should be fetched from the FRITZ!Box.")
	flag.StringVar(&cfg.DbPath, "dbpath", "./data", "Path to the directory where database is stored.")
	flag.UintVar(&cfg.Port, "port", 8090, "Listen port for the web UI.")
	flag.BoolVar(&cfg.DisableWebUi, "disableWebUi", false, "Disable the web UI.")
	flag.Parse()

	if usage {
		flag.Usage()
		return
	}
	if err := cfg.Validate(); err != nil {
		log.Printf("Error in configuration: %s\n", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	storage, _ := tstorage.NewStorage(
		tstorage.WithDataPath(cfg.DbPath),
		tstorage.WithTimestampPrecision(tstorage.Seconds),
	)
	defer storage.Close()

	httpClient := http.NewHttpClient(cfg.SkipTlsVerify)
	digestAuthClient := http.NewDigestAuthClient(httpClient, cfg.Username, cfg.Password)
	fb := fritzbox.New(ctx, digestAuthClient, cfg.Url)

	col := collector.New(ctx, cfg, fb, storage)
	wg.Add(1)
	go col.Start(wg)

	var srv *nethttp.Server
	if !cfg.DisableWebUi {
		app := ui.New(storage)
		app.RegisterRoutes()

		wg.Add(1)
		srv = startHttpServer(wg)
	} else {
		log.Println("Web UI disabled")
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	if srv != nil {
		_ = srv.Shutdown(ctx)
	}
	cancel()

	wg.Wait()
}

func startHttpServer(wg *sync.WaitGroup) *nethttp.Server {
	srv := &nethttp.Server{
		Addr: fmt.Sprintf(":%d", cfg.Port),
	}

	go func() {
		defer wg.Done()
		log.Printf("HTTP Server listening on http://127.0.0.1:%d/\n", cfg.Port)
		if err := srv.ListenAndServe(); !errors.Is(err, nethttp.ErrServerClosed) {
			log.Fatalf("HTTP server error: %s", err)
		}
	}()

	return srv
}

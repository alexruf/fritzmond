package collector

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/alexruf/fritzmond/config"
	"github.com/alexruf/fritzmond/fritzbox"
)

type Collector struct {
	ctx      context.Context
	cfg      config.Config
	fritzbox fritzbox.Fritzbox
}

func New(ctx context.Context, cfg config.Config, fritzbox fritzbox.Fritzbox) Collector {
	return Collector{
		ctx:      ctx,
		cfg:      cfg,
		fritzbox: fritzbox,
	}
}

func (c Collector) Start(wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	ticker := time.NewTicker(time.Duration(c.cfg.Interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			c.run()
		case <-c.ctx.Done():
			log.Println("Collector: Stopped")
			return
		}
	}
}

func (c Collector) run() {
	log.Println("Collecting data from FRITZ!Box")
	commonLinkProperties, err := c.fritzbox.GetCommonLinkProperties()
	if err != nil {
		log.Println(err.Error())
	}
	if commonLinkProperties != nil {
		log.Printf("CommonLinkProperties: %+v\n", *commonLinkProperties)
	} else {
		log.Println("CommonLinkProperties: <empty>")
	}

	totalBytesSent, err := c.fritzbox.GetTotalBytesSent()
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("TotalBytesSent: %d\n", totalBytesSent)

	totalBytesReceived, err := c.fritzbox.GetTotalBytesReceived()
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("TotalBytesReceived: %d\n", totalBytesReceived)

	totalPacketsSent, err := c.fritzbox.GetTotalPacketsSent()
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("TotalPacketsSent: %d\n", totalPacketsSent)

	totalPacketsReceived, err := c.fritzbox.GetTotalPacketsReceived()
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("TotalPacketsReceived: %d\n", totalPacketsReceived)
}

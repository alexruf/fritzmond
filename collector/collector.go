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
		log.Printf("GetCommonLinkProperties: %+v\n", *commonLinkProperties)
	} else {
		log.Println("GetCommonLinkProperties: <empty>")
	}

	totalBytesSent, err := c.fritzbox.GetTotalBytesSent()
	if err != nil {
		log.Println(err.Error())
	}
	if totalBytesSent != nil {
		log.Printf("GetTotalBytesSent: %+v\n", *totalBytesSent)
	} else {
		log.Println("GetTotalBytesSent: <empty>")
	}

	totalBytesReceived, err := c.fritzbox.GetTotalBytesReceived()
	if err != nil {
		log.Println(err.Error())
	}
	if totalBytesReceived != nil {
		log.Printf("GetTotalBytesReceived: %+v\n", *totalBytesReceived)
	} else {
		log.Println("GetTotalBytesReceived: <empty>")
	}

	totalPacketsSent, err := c.fritzbox.GetTotalPacketsSent()
	if err != nil {
		log.Println(err.Error())
	}
	if totalPacketsSent != nil {
		log.Printf("GetTotalPacketsSent: %+v\n", *totalPacketsSent)
	} else {
		log.Println("GetTotalPacketsSent: <empty>")
	}

	totalPacketsReceived, err := c.fritzbox.GetTotalPacketsReceived()
	if err != nil {
		log.Println(err.Error())
	}
	if totalPacketsReceived != nil {
		log.Printf("GetTotalPacketsReceived: %+v\n", *totalPacketsReceived)
	} else {
		log.Println("GetTotalPacketsReceived: <empty>")
	}
}
